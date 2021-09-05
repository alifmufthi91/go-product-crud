package database

import (
	"database/sql"
	"fmt"
	"ibf-benevolence/config"

	_ "github.com/go-sql-driver/mysql"
)

func DBConnection() (db *sql.DB) {
	var env = config.GetEnv()
	dbDriver := "mysql"
	dbHost := env.DBHost
	dbUser := env.DBUser
	dbPass := env.DBPassword
	dbPort := env.DBPort
	dbName := env.DBName
	url := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	fmt.Println("trying to connect DB : " + url)
	db, err := sql.Open(dbDriver, url)
	if err != nil {
		panic(err.Error())
	}
	return db
}
