package migrations

import (
	"product-crud/database"
	"product-crud/models"
	"product-crud/util/logger"
)

func Migrate() {
	logger.Info("Starting auto migration..")
	database.DBConnection().AutoMigrate(models.User{}, models.Product{})
}
