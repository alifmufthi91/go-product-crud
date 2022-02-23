package database

import (
	"fmt"
	"log"
	"os"
	"product-crud/config"
	"sync"
	"time"

	"product-crud/util/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

var (
	dbConnOnce sync.Once
	conn       *gorm.DB
)

func DBConnection() *gorm.DB {
	newLogger := gormLogger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		gormLogger.Config{
			SlowThreshold:             time.Second,     // Slow SQL threshold
			LogLevel:                  gormLogger.Info, // Log level
			IgnoreRecordNotFoundError: true,            // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,           // Disable color
		},
	)
	dbConnOnce.Do(func() {
		var env = config.Env
		dbHost := env.DBHost
		dbUser := env.DBUser
		dbPass := env.DBPassword
		dbPort := env.DBPort
		dbName := env.DBName
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", dbHost, dbUser, dbPass, dbName, dbPort)
		logger.Info("trying to connect DB : " + dsn)
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: newLogger,
		})
		if err != nil {
			logger.Error(err.Error())
			panic(err.Error())
		}
		conn = db
	})

	return conn
}
