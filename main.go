package main

import (
	"flag"
	"fmt"
	"product-crud/cache"
	"product-crud/config"
	"product-crud/database"
	"product-crud/migrations"
	"product-crud/repository"
	"product-crud/server"
	"product-crud/service"
	"product-crud/util/logger"
)

var (
	migrate = flag.Bool("migrate", false, "auto migrate the data models")
)

func main() {
	flag.Parse()
	logger.Init()
	var env = config.GetEnv()

	db := database.DBConnection()
	redis := database.RedisConnection()

	userRepository := repository.NewUserRepository(db)
	productRepository := repository.NewProductRepository(db)
	userService := service.NewUserService(userRepository)
	productService := service.NewProductService(productRepository, userRepository)

	cache.InitCache(redis)

	if *migrate {
		migrations.Migrate()
	}

	logger.Info("Starting server..")

	r := server.NewRouter(db, userRepository, productRepository, userService, productService)
	logger.Info(fmt.Sprintf("Running Server on Port: %s", env.Port))
	err := r.Run(fmt.Sprintf("localhost:%s", env.Port))
	if err != nil {
		panic(err)
	}

}
