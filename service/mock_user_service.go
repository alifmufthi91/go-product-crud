package service

import (
	"product-crud/dto/app"
	"product-crud/dto/request"
	"product-crud/dto/response"

	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) GetAll(pagination app.Pagination) (app.PaginatedResult[response.GetUserResponse], error) {
	args := m.Called(pagination)
	return args.Get(0).(app.PaginatedResult[response.GetUserResponse]), args.Error(1)
}

func (m *MockUserService) GetById(userId uint) (response.GetUserResponse, error) {
	args := m.Called(userId)
	return args.Get(0).(response.GetUserResponse), args.Error(1)
}

func (m *MockUserService) Register(userInput request.UserRegisterRequest) (response.GetUserResponse, error) {
	args := m.Called(userInput)
	return args.Get(0).(response.GetUserResponse), args.Error(1)
}

func (m *MockUserService) Login(userInput request.UserLoginRequest) (string, error) {
	args := m.Called(userInput)
	return args.Get(0).(string), args.Error(1)
}
