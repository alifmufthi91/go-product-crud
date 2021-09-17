package server

import (
	"ibf-benevolence/controller"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {

	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(gin.Logger())
	router.Use(cors.Default())

	api := router.Group("api")
	{
		users := api.Group("users")
		userController := controller.NewUserController()
		{
			users.GET("/", userController.GetAllUser)
			users.POST("/", userController.AddUser)
		}
	}

	return router
}
