package repository

import (
	"fmt"
	"ibf-benevolence/database"

	"github.com/jmoiron/sqlx"
)

type baseRepository struct {
	db         *sqlx.DB
	table      string
	primaryKey string
	tableAlias string
}

type BaseRepository interface {
	findAll(dest interface{}) error
}

func NewRepository(table string, primaryKey string, tableAlias string) BaseRepository {
	dbconn := database.DBConnection()
	return baseRepository{
		db:         dbconn,
		table:      table,
		primaryKey: primaryKey,
		tableAlias: tableAlias,
	}
}

func (repo baseRepository) findAll(dest interface{}) error {
	return repo.db.Select(dest, repo.findAllQuery())
}

func (repo baseRepository) findQuery(key string, value string) string {
	return fmt.Sprintf("SELECT * FROM %s WHERE '%s' = '%s'", repo.table, key, value)
}

func (repo baseRepository) findAllQuery() string {
	return fmt.Sprintf("SELECT * FROM %s", repo.table)
}
