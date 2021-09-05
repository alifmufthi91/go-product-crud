package repository

import (
	"database/sql"
	"fmt"
	"ibf-benevolence/database"
)

type baseRepository struct {
	db         *sql.DB
	table      string
	primaryKey string
	tableAlias string
}

type BaseRepository interface {
	findAll() (*sql.Rows, error)
}

func NewRepository(table string, primaryKey string, tableAlias string) BaseRepository {
	fmt.Println("Initializing base repository")
	dbconn := database.DBConnection()
	return baseRepository{
		db:         dbconn,
		table:      table,
		primaryKey: primaryKey,
		tableAlias: tableAlias,
	}
}

func (repo baseRepository) findAll() (*sql.Rows, error) {
	fmt.Println("Query: " + repo.findAllQuery())
	return repo.db.Query(repo.findAllQuery())
}

func (repo baseRepository) findQuery(key string, value string) string {
	return fmt.Sprintf("SELECT * FROM %s WHERE '%s' = '%s'", repo.table, key, value)
}

func (repo baseRepository) findAllQuery() string {
	return fmt.Sprintf("SELECT * FROM %s", repo.table)
}
