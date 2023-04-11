package controller_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"product-crud/controller"
	"product-crud/dto/request"
	"product-crud/service"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type UserControllerSuite struct {
	suite.Suite
	userService    *service.MockUserService
	userController controller.UserController
	router         *gin.Engine
}

func TestUserControllerSuite(t *testing.T) {
	suite.Run(t, new(UserControllerSuite))
}

func (c *UserControllerSuite) SetupSuite() {
	c.userService = new(service.MockUserService)
	c.userController = controller.NewUserController(c.userService)

	c.router = gin.Default()
	c.router.POST("/login", c.userController.LoginUser)
}

func (s *UserControllerSuite) AfterTest(_, _ string) {
	s.userService.AssertExpectations(s.T())
}

func (c *UserControllerSuite) TestUserController_Login() {

	// define mock user service login behavior
	token := "test-token"
	c.userService.On("Login", mock.Anything).Return(token, nil).Once()

	// create new context
	req, _ := http.NewRequest(http.MethodPost, "/login", nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	// create the body payload for test
	body := &request.UserLoginRequest{
		Email:    "testlogin@mail.com",
		Password: "Password",
	}

	// assign payload to context request body
	bodyBytes, _ := json.Marshal(body)
	req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	ctx.Request = req

	// call the endpoint
	c.router.ServeHTTP(w, req)

	// check response status code
	require.Equal(c.T(), http.StatusOK, w.Code)

	var responseBody map[string]string
	json.Unmarshal(w.Body.Bytes(), &responseBody)

	// check if token is inserted in data attribute of response
	require.NotEmpty(c.T(), responseBody["data"])

	// check if token is the same as expected
	require.Equal(c.T(), token, responseBody["data"])
}
