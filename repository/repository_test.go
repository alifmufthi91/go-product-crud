package repository_test

import (
	"context"
	"database/sql"
	"product-crud/models"
	"product-crud/repository"
	"regexp"
	"testing"
	"time"

	"github.com/go-test/deep"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestUserRepositorySuite(t *testing.T) {
	suite.Run(t, new(UserRepositorySuite))
}

type UserRepositorySuite struct {
	suite.Suite
	DB         *gorm.DB
	mock       sqlmock.Sqlmock
	repository repository.UserRepository
}

func (s *UserRepositorySuite) SetupSuite() {
	var (
		db  *sql.DB
		err error
	)

	db, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)

	s.DB, err = gorm.Open(mysql.New(mysql.Config{Conn: db, SkipInitializeWithVersion: true}), &gorm.Config{})
	require.NoError(s.T(), err)

	s.repository = repository.NewUserRepository(s.DB)
}

func (s *UserRepositorySuite) TestUserRepositoryGetByUserId() {

	product := models.Product{
		ID:                 1,
		ProductName:        "test",
		ProductDescription: "test",
		Photo:              "test.jpg",
		UploaderId:         1,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
		DeletedAt:          gorm.DeletedAt{},
	}

	user := models.User{
		ID:        1,
		FirstName: "Albert",
		LastName:  "Robb",
		Email:     "albert@robb@email.com",
		Products:  []models.Product{product},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: gorm.DeletedAt{},
	}

	userRows := sqlmock.NewRows([]string{
		"id",
		"first_name",
		"last_name",
		"email",
		"created_at",
		"updated_at",
		"deleted_at",
	}).AddRow(
		user.ID,
		user.FirstName,
		user.LastName,
		user.Email,
		user.CreatedAt,
		user.UpdatedAt,
		user.DeletedAt,
	)

	productRows := sqlmock.NewRows([]string{
		"id",
		"product_name",
		"product_description",
		"photo",
		"uploader_id",
		"created_at",
		"updated_at",
		"deleted_at",
	}).AddRow(
		product.ID,
		product.ProductName,
		product.ProductDescription,
		product.Photo,
		product.UploaderId,
		product.CreatedAt,
		product.UpdatedAt,
		product.DeletedAt,
	)

	const expectUser = "SELECT * FROM `users` WHERE users.id = ? AND `users`.`deleted_at` IS NULL"
	s.mock.ExpectQuery(regexp.QuoteMeta(expectUser)).WithArgs(user.ID).WillReturnRows(userRows)

	const expectProduct = "SELECT * FROM `products` WHERE `products`.`uploader_id` = ? AND `products`.`deleted_at` IS NULL"
	s.mock.ExpectQuery(regexp.QuoteMeta(expectProduct)).WithArgs(product.UploaderId).WillReturnRows(productRows)

	const expectUser2 = "SELECT * FROM `users` WHERE `users`.`id` = ? AND `users`.`deleted_at` IS NULL"
	s.mock.ExpectQuery(regexp.QuoteMeta(expectUser2)).WithArgs(product.UploaderId).WillReturnRows(userRows)

	res, err := s.repository.GetByUserId(context.Background(), 1)
	require.NoError(s.T(), err)
	require.Nil(s.T(), deep.Equal(user, *res))
}

func (s *UserRepositorySuite) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}
