package service

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/itsahyarr/gofiber-boilerplate/internal/user/dto"
	"github.com/itsahyarr/gofiber-boilerplate/internal/user/repository"
	"github.com/itsahyarr/gofiber-boilerplate/pkg/database"
	"github.com/itsahyarr/gofiber-boilerplate/pkg/logger"
	"github.com/itsahyarr/gofiber-boilerplate/shared/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.uber.org/zap"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidOldPassword = errors.New("invalid old password")
)

// UserService defines the interface for user operations
type UserService interface {
	GetByID(ctx context.Context, id string) (*dto.UserResponse, error)
	GetAll(ctx context.Context, filter bson.M, page, pageSize int) ([]dto.UserResponse, int64, error)
	Update(ctx context.Context, id string, req *dto.UpdateUserRequest) (*dto.UserResponse, error)
	Delete(ctx context.Context, id string) error
	ChangePassword(ctx context.Context, id string, req *dto.ChangePasswordRequest) error
	RegisterWithStats(ctx context.Context, user *entity.User) error
}

type userServiceImpl struct {
	userRepo repository.UserRepository
	db       *database.MongoDB
}

// NewUserService creates a new user service
func NewUserService(userRepo repository.UserRepository, db *database.MongoDB) UserService {
	return &userServiceImpl{
		userRepo: userRepo,
		db:       db,
	}
}

func (s *userServiceImpl) GetByID(ctx context.Context, id string) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		logger.Error("failed to get user by id", zap.Error(err), zap.String("user_id", id))
		return nil, err
	}

	response := dto.ToUserResponse(user)
	return &response, nil
}

func (s *userServiceImpl) GetAll(ctx context.Context, filter bson.M, page, pageSize int) ([]dto.UserResponse, int64, error) {
	users, total, err := s.userRepo.FindAll(ctx, filter, page, pageSize)
	if err != nil {
		logger.Error("failed to get all users", zap.Error(err))
		return nil, 0, err
	}

	return dto.ToUserResponses(users), total, nil
}

func (s *userServiceImpl) Update(ctx context.Context, id string, req *dto.UpdateUserRequest) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		logger.Error("failed to find user for update", zap.Error(err), zap.String("user_id", id))
		return nil, err
	}

	// Update fields if provided
	if req.FirstName != nil {
		user.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		user.LastName = *req.LastName
	}
	if req.Role != nil {
		user.Role = *req.Role
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		logger.Error("failed to update user", zap.Error(err), zap.String("user_id", id))
		return nil, err
	}

	logger.Info("user updated successfully", zap.String("user_id", id))

	response := dto.ToUserResponse(user)
	return &response, nil
}

func (s *userServiceImpl) Delete(ctx context.Context, id string) error {
	if err := s.userRepo.Delete(ctx, id); err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return ErrUserNotFound
		}
		logger.Error("failed to delete user", zap.Error(err), zap.String("user_id", id))
		return err
	}

	logger.Info("user deleted successfully", zap.String("user_id", id))
	return nil
}

func (s *userServiceImpl) ChangePassword(ctx context.Context, id string, req *dto.ChangePasswordRequest) error {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return ErrUserNotFound
		}
		logger.Error("failed to find user for password change", zap.Error(err), zap.String("user_id", id))
		return err
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		return ErrInvalidOldPassword
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("failed to hash new password", zap.Error(err))
		return err
	}

	user.Password = string(hashedPassword)
	if err := s.userRepo.Update(ctx, user); err != nil {
		logger.Error("failed to update password", zap.Error(err), zap.String("user_id", id))
		return err
	}

	logger.Info("password changed successfully", zap.String("user_id", id))
	return nil
}

func (s *userServiceImpl) RegisterWithStats(ctx context.Context, user *entity.User) error {
	// 1. Start a Session
	session, err := s.db.Client.StartSession()
	if err != nil {
		logger.Error("failed to start mongodb session", zap.Error(err))
		return err
	}
	defer session.EndSession(ctx)

	// 2. Use WithTransaction to wrap your logic
	_, err = session.WithTransaction(ctx, func(sessCtx context.Context) (interface{}, error) {
		// --- OPERATION 1: Create User ---
		// We use the repository but we must ensure the repository methods can accept the session context
		// Currently the repo methods use the context passed to them.
		if err := s.userRepo.Create(sessCtx, user); err != nil {
			logger.Error("failed to create user in transaction", zap.Error(err))
			return nil, err
		}

		// --- OPERATION 2: Create associated data (Stats) ---
		statsCollection := s.db.Collection("user_stats")
		userStats := bson.M{
			"user_id":     user.ID,
			"login_count": 0,
			"points":      100, // Welcome bonus
		}

		if _, err := statsCollection.InsertOne(sessCtx, userStats); err != nil {
			logger.Error("failed to create user stats in transaction", zap.Error(err))
			return nil, err
		}

		return nil, nil
	})

	if err != nil {
		logger.Error("transaction failed and rolled back", zap.Error(err))
		return err
	}

	logger.Info("user and stats created successfully within transaction", zap.String("user_id", user.ID.Hex()))
	return nil
}
