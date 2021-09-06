package repository

import (
	"fmt"
	"ibf-benevolence/database"
	"ibf-benevolence/util/logger"

	"github.com/jmoiron/sqlx"
)

type baseRepository struct {
	db         *sqlx.DB
	table      string
	primaryKey string
	tableAlias string
}

type BaseRepository interface {
	selectAll(dest interface{}) error
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

func (repo baseRepository) selectAll(dest interface{}) error {
	logger.Info("Querying: " + repo.selectAllQuery())
	return repo.db.Select(dest, repo.selectAllQuery())
}

func (repo baseRepository) selectQuery(key string, value string) string {
	return fmt.Sprintf("SELECT * FROM %s WHERE '%s' = '%s'", repo.table, key, value)
}

func (repo baseRepository) selectAllQuery() string {
	return fmt.Sprintf("SELECT * FROM %s", repo.table)
}
