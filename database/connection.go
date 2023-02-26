package database

import (
	"fmt"
	"product-crud/config"
	"sync"
	"time"

	"product-crud/util/logger"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

var (
	dbConnOnce sync.Once
	conn       *gorm.DB
)

func DBConnection() *gorm.DB {
	newLogger := gormLogger.New(
		&log.Logger, // io writer
		gormLogger.Config{
			SlowThreshold:             time.Second,     // Slow SQL threshold
			LogLevel:                  gormLogger.Info, // Log level
			IgnoreRecordNotFoundError: true,            // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,           // Disable color
		},
	)
	dbConnOnce.Do(func() {
		var env = config.GetEnv()
		dbHost := env.DBHost
		dbUser := env.DBUser
		dbPass := env.DBPassword
		dbPort := env.DBPort
		dbName := env.DBName
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPass, dbHost, dbPort, dbName)
		logger.Info("trying to connect DB : " + dsn)
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: newLogger,
		})
		if err != nil {
			logger.Error(err.Error())
			panic(err)
		}
		conn = db
	})

	return conn
}
