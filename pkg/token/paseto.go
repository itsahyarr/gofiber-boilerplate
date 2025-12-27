package token

import (
	"errors"
	"time"

	"github.com/o1egl/paseto/v2"
)

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

// Payload contains the payload data of the token
type Payload struct {
	UserID    string    `json:"user_id"`
	Role      string    `json:"role"`
	TokenType string    `json:"token_type"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

// Valid checks if the token payload is valid
func (p *Payload) Valid() error {
	if time.Now().After(p.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}

// PasetoMaker is a PASETO token maker
type PasetoMaker struct {
	symmetricKey []byte
	paseto       *paseto.V2
}

// NewPasetoMaker creates a new PASETO token maker
func NewPasetoMaker(symmetricKey string) (*PasetoMaker, error) {
	if len(symmetricKey) != 32 {
		return nil, errors.New("symmetric key must be exactly 32 characters")
	}

	return &PasetoMaker{
		symmetricKey: []byte(symmetricKey),
		paseto:       paseto.NewV2(),
	}, nil
}

// CreateAccessToken creates a new access token for a specific user
func (m *PasetoMaker) CreateAccessToken(userID, role string, duration time.Duration) (string, *Payload, error) {
	payload := &Payload{
		UserID:    userID,
		Role:      role,
		TokenType: "access",
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	token, err := m.paseto.Encrypt(m.symmetricKey, payload, nil)
	if err != nil {
		return "", nil, err
	}

	return token, payload, nil
}

// CreateRefreshToken creates a new refresh token for a specific user
func (m *PasetoMaker) CreateRefreshToken(userID, role string, duration time.Duration) (string, *Payload, error) {
	payload := &Payload{
		UserID:    userID,
		Role:      role,
		TokenType: "refresh",
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	token, err := m.paseto.Encrypt(m.symmetricKey, payload, nil)
	if err != nil {
		return "", nil, err
	}

	return token, payload, nil
}

// VerifyToken verifies the token and returns the payload
func (m *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}

	err := m.paseto.Decrypt(token, m.symmetricKey, payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	return payload, nil
}
