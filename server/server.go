package server

import (
	"fmt"
	"ibf-benevolence/config"
)

func Init() {
	config.GetEnv()
	r := NewRouter()
	fmt.Println("Starting server..")
	r.Run()
}
