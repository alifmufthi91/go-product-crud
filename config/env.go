package config

import (
	"log"
	"os"
	"product-crud/util/logger"
	"sync"

	"github.com/joho/godotenv"
)

var env Env
var once sync.Once

type Env struct {
	Port string `envconfig:"PORT"`

	DBHost     string `envconfig:"DB_HOST"`
	DBUser     string `envconfig:"DB_USERNAME"`
	DBPassword string `envconfig:"DB_PASSWORD"`
	DBName     string `envconfig:"DB_NAME"`
	DBPort     string `envconfig:"DB_PORT"`
}

func GetEnv() *Env {
	once.Do(func() {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
		env.Port = os.Getenv("PORT")
		env.DBHost = os.Getenv("DB_HOST")
		env.DBUser = os.Getenv("DB_USERNAME")
		env.DBPassword = os.Getenv("DB_PASSWORD")
		env.DBName = os.Getenv("DB_NAME")
		env.DBPort = os.Getenv("DB_PORT")
		logger.Info("Environment config set")
	})
	return &env
}
