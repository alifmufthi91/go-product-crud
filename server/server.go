package server

import (
	"fmt"
	"product-crud/config"
	"product-crud/migrations"
	"product-crud/util/logger"
)

func Init() {
	logger.Info("Starting server..")
	var env = config.GetEnv()
	r := NewRouter()
	migrations.Migrate()
	logger.Info(fmt.Sprintf("Running Server on Port: %s", env.Port))
	r.Run(fmt.Sprintf("localhost:%s", env.Port))
}
