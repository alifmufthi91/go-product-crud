package main

import (
	"fmt"
	"product-crud/cache"
	"product-crud/config"
	"product-crud/database"
	"product-crud/repository"
	"product-crud/server"
	"product-crud/service"
	"product-crud/util/logger"
)

func main() {
	logger.Init()
	config.InitEnv()

	var env = config.GetEnv()

	db := database.DBConnection()
	redis := database.RedisConnection()

	userRepository := repository.NewUserRepository(db)
	productRepository := repository.NewProductRepository(db)
	userService := service.NewUserService(userRepository)
	productService := service.NewProductService(productRepository, userRepository)

	cache.InitCache(redis)

	logger.Info("Starting server..")

	r := server.NewRouter(db, userRepository, productRepository, userService, productService)
	// migrations.Migrate()
	logger.Info(fmt.Sprintf("Running Server on Port: %s", env.Port))
	err := r.Run(fmt.Sprintf("localhost:%s", env.Port))
	if err != nil {
		panic(err)
	}

}
