package database

import (
	"fmt"
	"product-crud/config"
	"product-crud/util/logger"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	dbConnOnce sync.Once
	conn       *gorm.DB
)

func DBConnection() (db *gorm.DB) {
	dbConnOnce.Do(func() {
		var env = config.Env
		dbHost := env.DBHost
		dbUser := env.DBUser
		dbPass := env.DBPassword
		dbPort := env.DBPort
		dbName := env.DBName
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", dbHost, dbUser, dbPass, dbName, dbPort)
		logger.Info("trying to connect DB : " + dsn)
		db, err := gorm.Open(postgres.Open(dsn))
		if err != nil {
			logger.Error(err.Error())
			panic(err.Error())
		}
		conn = db
	})

	return conn
}
