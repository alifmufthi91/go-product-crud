package database

import (
	"database/sql"
	"fmt"
	"ibf-benevolence/config"
	"ibf-benevolence/util/logger"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

var (
	dbConnOnce sync.Once
	conn       *sql.DB
)

func DBConnection() (db *sql.DB) {
	dbConnOnce.Do(func() {
		var env = config.GetEnv()
		dbDriver := "mysql"
		dbHost := env.DBHost
		dbUser := env.DBUser
		dbPass := env.DBPassword
		dbPort := env.DBPort
		dbName := env.DBName
		url := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
		logger.Info("trying to connect DB : " + url)
		db, err := sql.Open(dbDriver, url)
		if err != nil {
			logger.Error(err.Error())
			panic(err.Error())
		}
		conn = db
	})

	return conn
}
