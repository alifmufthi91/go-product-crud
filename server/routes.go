package server

import (
	"product-crud/controller"
	"product-crud/middlewares"

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
			users.GET("/:id", middlewares.Auth, userController.GetUserById)
			users.POST("/", userController.RegisterUser)
			users.POST("/login", userController.LoginUser)
		}
		products := api.Group("products")
		productController := controller.NewProductController()
		{
			// products.GET("/", productController)
			// users.GET("/:id", middlewares.Auth, userController.GetUserById)
			products.POST("/", middlewares.Auth, productController.AddProduct)
		}
	}

	return router
}
