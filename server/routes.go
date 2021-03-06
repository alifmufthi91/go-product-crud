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
			users.GET("/", middlewares.Auth, userController.GetAllUser)
			users.GET("/:id", middlewares.Auth, userController.GetUserById)
			users.POST("/", userController.RegisterUser)
			users.POST("/login", userController.LoginUser)
		}
		products := api.Group("products")
		productController := controller.NewProductController()
		{
			products.GET("/", middlewares.Auth, productController.GetAllProduct)
			products.GET("/:id", middlewares.Auth, productController.GetProductById)
			products.POST("/", middlewares.Auth, productController.AddProduct)
			products.PATCH("/:id", middlewares.Auth, productController.UpdateProduct)
			products.DELETE("/:id", middlewares.Auth, productController.DeleteProduct)
		}
		files := api.Group("files")
		filesController := controller.NewFileController()
		{
			files.POST("/upload", middlewares.Auth, middlewares.BodySizeMiddleware, filesController.Upload)
			files.GET("/:name", middlewares.Auth, filesController.Download)
		}
	}

	return router
}
