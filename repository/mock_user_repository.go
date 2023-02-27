package repository

import (
	"context"
	"product-crud/dto/app"
	"product-crud/models"

	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetAllUser(ctx context.Context, pagination app.Pagination, count *int64) ([]*models.User, error) {
	args := m.Called(ctx, pagination, count)
	return args.Get(0).([]*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByUserId(ctx context.Context, userId uint) (*models.User, error) {
	args := m.Called(ctx, userId)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) AddUser(ctx context.Context, user models.User) (*models.User, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) IsExistingEmail(ctx context.Context, email string) (*bool, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(*bool), args.Error(1)
}
