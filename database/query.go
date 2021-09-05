package database

// import (
// 	"database/sql"
// 	"log"
// )

// type Query func(string, *sql.DB) *sql.Rows

// func (db *Db) query(sql string, connection *sql.DB) *sql.Rows {
// 	log.Printf("query: %s initiated", sql)
// 	rows, err := connection.Query(sql)
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	log.Printf("query: %s success", sql)
// 	return rows
// }
