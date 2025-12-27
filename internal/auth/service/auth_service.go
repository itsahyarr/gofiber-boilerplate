package service

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/itsahyarr/gofiber-boilerplate/internal/auth/dto"
	"github.com/itsahyarr/gofiber-boilerplate/internal/auth/repository"
	"github.com/itsahyarr/gofiber-boilerplate/internal/config"
	userRepo "github.com/itsahyarr/gofiber-boilerplate/internal/user/repository"
	"github.com/itsahyarr/gofiber-boilerplate/pkg/logger"
	"github.com/itsahyarr/gofiber-boilerplate/pkg/token"
	"github.com/itsahyarr/gofiber-boilerplate/shared/entity"
	"go.uber.org/zap"
)

var (
	ErrEmailAlreadyExists  = errors.New("email already exists")
	ErrInvalidCredentials  = errors.New("invalid email or password")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrUserNotActive       = errors.New("user account is not active")
)

// AuthService defines the interface for authentication operations
type AuthService interface {
	Register(ctx context.Context, req *dto.RegisterRequest) (*dto.AuthResponse, error)
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.AuthResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*dto.TokenResponse, error)
	Logout(ctx context.Context, userID string) error
}

type authServiceImpl struct {
	userRepo   userRepo.UserRepository
	tokenRepo  repository.TokenRepository
	tokenMaker *token.PasetoMaker
	config     *config.Config
}

// NewAuthService creates a new authentication service
func NewAuthService(
	userRepository userRepo.UserRepository,
	tokenRepository repository.TokenRepository,
	tokenMaker *token.PasetoMaker,
	cfg *config.Config,
) AuthService {
	return &authServiceImpl{
		userRepo:   userRepository,
		tokenRepo:  tokenRepository,
		tokenMaker: tokenMaker,
		config:     cfg,
	}
}

func (s *authServiceImpl) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.AuthResponse, error) {
	// Check if email already exists
	exists, err := s.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		logger.Error("failed to check email existence", zap.Error(err))
		return nil, err
	}
	if exists {
		return nil, ErrEmailAlreadyExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("failed to hash password", zap.Error(err))
		return nil, err
	}

	// Create user
	user := &entity.User{
		Email:     req.Email,
		Password:  string(hashedPassword),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      entity.RoleUser, // Default role
		IsActive:  true,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		logger.Error("failed to create user", zap.Error(err))
		return nil, err
	}

	// Generate tokens
	accessToken, _, err := s.tokenMaker.CreateAccessToken(
		user.ID.Hex(),
		string(user.Role),
		s.config.Token.AccessTokenDuration,
	)
	if err != nil {
		logger.Error("failed to create access token", zap.Error(err))
		return nil, err
	}

	refreshToken, _, err := s.tokenMaker.CreateRefreshToken(
		user.ID.Hex(),
		string(user.Role),
		s.config.Token.RefreshTokenDuration,
	)
	if err != nil {
		logger.Error("failed to create refresh token", zap.Error(err))
		return nil, err
	}

	// Store refresh token in Redis
	if err := s.tokenRepo.Store(ctx, user.ID.Hex(), refreshToken, s.config.Token.RefreshTokenDuration); err != nil {
		logger.Error("failed to store refresh token", zap.Error(err))
		return nil, err
	}

	logger.Info("user registered successfully", zap.String("user_id", user.ID.Hex()))

	return &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         toUserResponse(user),
	}, nil
}

func (s *authServiceImpl) Login(ctx context.Context, req *dto.LoginRequest) (*dto.AuthResponse, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, userRepo.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		logger.Error("failed to find user", zap.Error(err))
		return nil, err
	}

	// Check if user is active
	if !user.IsActive {
		return nil, ErrUserNotActive
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Generate tokens
	accessToken, _, err := s.tokenMaker.CreateAccessToken(
		user.ID.Hex(),
		string(user.Role),
		s.config.Token.AccessTokenDuration,
	)
	if err != nil {
		logger.Error("failed to create access token", zap.Error(err))
		return nil, err
	}

	refreshToken, _, err := s.tokenMaker.CreateRefreshToken(
		user.ID.Hex(),
		string(user.Role),
		s.config.Token.RefreshTokenDuration,
	)
	if err != nil {
		logger.Error("failed to create refresh token", zap.Error(err))
		return nil, err
	}

	// Store refresh token in Redis
	if err := s.tokenRepo.Store(ctx, user.ID.Hex(), refreshToken, s.config.Token.RefreshTokenDuration); err != nil {
		logger.Error("failed to store refresh token", zap.Error(err))
		return nil, err
	}

	logger.Info("user logged in successfully", zap.String("user_id", user.ID.Hex()))

	return &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         toUserResponse(user),
	}, nil
}

func (s *authServiceImpl) RefreshToken(ctx context.Context, refreshTokenStr string) (*dto.TokenResponse, error) {
	// Verify refresh token
	payload, err := s.tokenMaker.VerifyToken(refreshTokenStr)
	if err != nil {
		return nil, ErrInvalidRefreshToken
	}

	// Check if refresh token exists in Redis
	storedToken, err := s.tokenRepo.Get(ctx, payload.UserID)
	if err != nil {
		return nil, ErrInvalidRefreshToken
	}

	// Verify stored token matches
	if storedToken != refreshTokenStr {
		return nil, ErrInvalidRefreshToken
	}

	// Generate new tokens
	accessToken, _, err := s.tokenMaker.CreateAccessToken(
		payload.UserID,
		payload.Role,
		s.config.Token.AccessTokenDuration,
	)
	if err != nil {
		logger.Error("failed to create access token", zap.Error(err))
		return nil, err
	}

	newRefreshToken, _, err := s.tokenMaker.CreateRefreshToken(
		payload.UserID,
		payload.Role,
		s.config.Token.RefreshTokenDuration,
	)
	if err != nil {
		logger.Error("failed to create refresh token", zap.Error(err))
		return nil, err
	}

	// Update refresh token in Redis
	if err := s.tokenRepo.Store(ctx, payload.UserID, newRefreshToken, s.config.Token.RefreshTokenDuration); err != nil {
		logger.Error("failed to store refresh token", zap.Error(err))
		return nil, err
	}

	logger.Info("tokens refreshed successfully", zap.String("user_id", payload.UserID))

	return &dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *authServiceImpl) Logout(ctx context.Context, userID string) error {
	// Delete refresh token from Redis
	if err := s.tokenRepo.Delete(ctx, userID); err != nil {
		logger.Error("failed to delete refresh token", zap.Error(err), zap.String("user_id", userID))
		return err
	}

	logger.Info("user logged out successfully", zap.String("user_id", userID))
	return nil
}

func toUserResponse(user *entity.User) dto.UserResponse {
	return dto.UserResponse{
		ID:        user.ID.Hex(),
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      string(user.Role),
	}
}
