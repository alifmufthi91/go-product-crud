package service_test

import (
	"crypto/sha256"
	"errors"
	"product-crud/dto/request"
	"product-crud/dto/response"
	"product-crud/models"
	"product-crud/repository"
	"product-crud/service"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type UserServiceSuite struct {
	suite.Suite
	userRepository *repository.MockUserRepository
	userService    service.IUserService
}

func TestUserServiceSuite(t *testing.T) {
	suite.Run(t, new(UserServiceSuite))
}

func (s *UserServiceSuite) SetupSuite() {
	s.userRepository = new(repository.MockUserRepository)
	s.userService = service.NewUserService(s.userRepository)
}

func (s *UserServiceSuite) AfterTest(_, _ string) {
	s.userRepository.AssertExpectations(s.T())
}

func (s *UserServiceSuite) TestUserService_RegisterUser() {
	userRequest := request.UserRegisterRequest{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		Password:  "Password",
	}

	bv := []byte(userRequest.Password)
	hasher := sha256.New()
	hasher.Write(bv)

	user := models.User{
		FirstName: userRequest.FirstName,
		LastName:  userRequest.LastName,
		Email:     userRequest.Email,
		Password:  hasher.Sum(nil),
	}

	response := response.GetUserResponse{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Products:  []response.GetProductResponse{},
	}

	// expect success
	existing := false
	s.userRepository.On("IsExistingEmail", mock.Anything, userRequest.Email).Return(&existing, nil).Once()
	s.userRepository.On("AddUser", mock.Anything, user).Return(&user, nil).Once()
	newUser, err := s.userService.Register(userRequest)
	require.NoError(s.T(), err)
	require.Equal(s.T(), response, *newUser)

	// expect error because email is existed
	existing = true
	expectedErr := errors.New("email is already exist")
	s.userRepository.On("IsExistingEmail", mock.Anything, user.Email).Return(&existing, expectedErr).Once()
	newUser, err = s.userService.Register(userRequest)
	require.Error(s.T(), err, expectedErr.Error())
	require.Nil(s.T(), newUser)

	// expect error because process AddUser is having problem
	existing = false
	var emptyUser *models.User
	expectedErr = errors.New("error happen during add user")
	s.userRepository.On("IsExistingEmail", mock.Anything, user.Email).Return(&existing, nil).Once()
	s.userRepository.On("AddUser", mock.Anything, user).Return(emptyUser, expectedErr).Once()
	newUser, err = s.userService.Register(userRequest)
	require.Error(s.T(), err, expectedErr.Error())
	require.Nil(s.T(), newUser)
}
