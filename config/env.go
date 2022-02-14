package config

import (
	"log"
	"os"
	"product-crud/util/logger"
	"sync"

	"github.com/joho/godotenv"
)

var Env Environment
var once sync.Once

type Environment struct {
	Port string `envconfig:"PORT"`

	DBHost     string `envconfig:"DB_HOST"`
	DBUser     string `envconfig:"DB_USERNAME"`
	DBPassword string `envconfig:"DB_PASSWORD"`
	DBName     string `envconfig:"DB_NAME"`
	DBPort     string `envconfig:"DB_PORT"`

	JWTSECRET string `envconfig:"JWT_SECRET"`
}

func GetEnv() *Environment {
	once.Do(func() {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
		Env.Port = os.Getenv("PORT")
		Env.DBHost = os.Getenv("DB_HOST")
		Env.DBUser = os.Getenv("DB_USERNAME")
		Env.DBPassword = os.Getenv("DB_PASSWORD")
		Env.DBName = os.Getenv("DB_NAME")
		Env.DBPort = os.Getenv("DB_PORT")
		Env.JWTSECRET = os.Getenv("JWT_SECRET")
		logger.Info("Environment config set")
	})
	return &Env
}
