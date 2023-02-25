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

	FilePath string `envconfig:"FILE_PATH"`

	RedisHost     string `envconfig:"REDIS_HOST"`
	RedisPort     string `envconfig:"REDIS_PORT"`
	RedisPassword string `envconfig:"REDIS_PASSWORD"`
}

func InitEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	Env = Environment{
		Port:          os.Getenv("PORT"),
		DBHost:        os.Getenv("DB_HOST"),
		DBUser:        os.Getenv("DB_USERNAME"),
		DBPassword:    os.Getenv("DB_PASSWORD"),
		DBName:        os.Getenv("DB_NAME"),
		DBPort:        os.Getenv("DB_PORT"),
		JWTSECRET:     os.Getenv("JWT_SECRET"),
		FilePath:      os.Getenv("FILE_PATH"),
		RedisHost:     os.Getenv("REDIS_HOST"),
		RedisPort:     os.Getenv("REDIS_PORT"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
	}
	logger.Info("Environment config set")
}

func GetEnv() *Environment {
	once.Do(InitEnv)
	return &Env
}
