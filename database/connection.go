package database

import (
	"fmt"
	"ibf-benevolence/config"
	"ibf-benevolence/util/logger"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	dbConnOnce sync.Once
	conn       *gorm.DB
)

func DBConnection() (db *gorm.DB) {
	dbConnOnce.Do(func() {
		var env = config.GetEnv()
		dbHost := env.DBHost
		dbUser := env.DBUser
		dbPass := env.DBPassword
		dbPort := env.DBPort
		dbName := env.DBName
		url := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
		logger.Info("trying to connect DB : " + url)
		db, err := gorm.Open(mysql.Open(url))
		if err != nil {
			logger.Error(err.Error())
			panic(err.Error())
		}
		conn = db
	})

	return conn
}
