package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/itsahyarr/gofiber-boilerplate/internal/user/repository"
	"github.com/itsahyarr/gofiber-boilerplate/internal/user/repository/mock"
	"github.com/itsahyarr/gofiber-boilerplate/shared/entity"
)

func TestGetByID(t *testing.T) {
	// 1. Setup Mock
	mockRepo := &mock.MockUserRepository{
		FindAllFunc: func(ctx context.Context, filter bson.M, page, pageSize int) ([]*entity.User, int64, error) {
			return nil, 0, nil
		},
		FindByIDFunc: func(ctx context.Context, id string) (*entity.User, error) {
			objID, _ := bson.ObjectIDFromHex("658bd7c1f1e29e0001bcdefg")
			return &entity.User{
				ID:        objID,
				Email:     "test@example.com",
				FirstName: "John",
				LastName:  "Doe",
			}, nil
		},
	}

	// 2. Initialize Service with Mock
	// Note: We pass nil for MongoDB since GetByID doesn't use it for transactions
	service := NewUserService(mockRepo, nil)

	// 3. Call Method
	res, err := service.GetByID(context.Background(), "658bd7c1f1e29e0001bcdefg")

	// 4. Assertions
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "test@example.com", res.Email)
	assert.Equal(t, "John", res.FirstName)
}

func TestGetByID_NotFound(t *testing.T) {
	// 1. Setup Mock to return NotFound
	mockRepo := &mock.MockUserRepository{
		FindAllFunc: func(ctx context.Context, filter bson.M, page, pageSize int) ([]*entity.User, int64, error) {
			return nil, 0, nil
		},
		FindByIDFunc: func(ctx context.Context, id string) (*entity.User, error) {
			return nil, repository.ErrUserNotFound
		},
	}

	service := NewUserService(mockRepo, nil)

	// 2. Call Method
	res, err := service.GetByID(context.Background(), "invalid-id")

	// 3. Assertions
	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, ErrUserNotFound, err)
}
