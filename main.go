package main

import (
	"product-crud/cache"
	"product-crud/config"
	"product-crud/server"
	"product-crud/util/logger"
)

func main() {
	logger.Init()
	config.EnvInit()
	cache.Init()
	server.Init()
}
