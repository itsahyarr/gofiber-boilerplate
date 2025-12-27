package mock

import (
	"context"

	"github.com/itsahyarr/gofiber-boilerplate/shared/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// MockUserRepository is a manual mock for the UserRepository interface
type MockUserRepository struct {
	CreateFunc        func(ctx context.Context, user *entity.User) error
	FindByIDFunc      func(ctx context.Context, id string) (*entity.User, error)
	FindByEmailFunc   func(ctx context.Context, email string) (*entity.User, error)
	FindAllFunc       func(ctx context.Context, filter bson.M, page, pageSize int) ([]*entity.User, int64, error)
	UpdateFunc        func(ctx context.Context, user *entity.User) error
	DeleteFunc        func(ctx context.Context, id string) error
	ExistsByEmailFunc func(ctx context.Context, email string) (bool, error)
}

func (m *MockUserRepository) Create(ctx context.Context, user *entity.User) error {
	return m.CreateFunc(ctx, user)
}

func (m *MockUserRepository) FindByID(ctx context.Context, id string) (*entity.User, error) {
	return m.FindByIDFunc(ctx, id)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	return m.FindByEmailFunc(ctx, email)
}

func (m *MockUserRepository) FindAll(ctx context.Context, filter bson.M, page, pageSize int) ([]*entity.User, int64, error) {
	return m.FindAllFunc(ctx, filter, page, pageSize)
}

func (m *MockUserRepository) Update(ctx context.Context, user *entity.User) error {
	return m.UpdateFunc(ctx, user)
}

func (m *MockUserRepository) Delete(ctx context.Context, id string) error {
	return m.DeleteFunc(ctx, id)
}

func (m *MockUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	return m.ExistsByEmailFunc(ctx, email)
}
