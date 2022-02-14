package main

import (
	"product-crud/server"
	"product-crud/util/logger"
)

func main() {
	logger.Init()
	server.Init()
}
