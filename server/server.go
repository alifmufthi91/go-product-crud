package server

import (
	"fmt"
	"product-crud/config"
	"product-crud/util/logger"
)

func Init() {
	logger.Info("Starting server..")
	var env = config.GetEnv()
	r := NewRouter()
	// migrations.Migrate()
	logger.Info(fmt.Sprintf("Running Server on Port: %s", env.Port))
	err := r.Run(fmt.Sprintf("localhost:%s", env.Port))
	if err != nil {
		panic(err)
	}
}
