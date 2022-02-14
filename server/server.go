package server

import (
	"fmt"
	"ibf-benevolence/config"
	"ibf-benevolence/util/logger"
)

func Init() {
	logger.Info("Starting server..")
	var env = config.GetEnv()
	r := NewRouter()
	logger.Info(fmt.Sprintf("Running Server on Port: %s", env.Port))
	r.Run(fmt.Sprintf("localhost:%s", env.Port))
}
