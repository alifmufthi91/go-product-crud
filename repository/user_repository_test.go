package repository_test

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"product-crud/dto/app"
	"product-crud/models"
	"product-crud/repository"
	"regexp"
	"testing"

	"github.com/go-test/deep"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type UserRepositorySuite struct {
	suite.Suite
	DB         *gorm.DB
	mock       sqlmock.Sqlmock
	repository repository.UserRepository
}

func TestUserRepositorySuite(t *testing.T) {
	suite.Run(t, new(UserRepositorySuite))
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

func (s *UserRepositorySuite) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func (s *UserRepositorySuite) TestUserRepository_GetByUserId() {

	user := models.User{
		FirstName: "Albert",
		LastName:  "Robb",
		Email:     "albert@robb@email.com",
		Products:  nil,
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

	const expectUser = "SELECT * FROM `users` WHERE users.id = ? AND `users`.`deleted_at` IS NULL"
	s.mock.ExpectQuery(regexp.QuoteMeta(expectUser)).WithArgs(1).WillReturnRows(userRows)

	res, err := s.repository.GetByUserId(context.Background(), 1)
	require.NoError(s.T(), err)
	require.Nil(s.T(), deep.Equal(user, *res))
}

func (s *UserRepositorySuite) TestUserRepository_GetByEmail() {

	user := models.User{
		FirstName: "Albert",
		LastName:  "Robb",
		Email:     "albert@robb@email.com",
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

	const expectUser = "SELECT * FROM `users` WHERE email = ? AND `users`.`deleted_at` IS NULL"
	s.mock.ExpectQuery(regexp.QuoteMeta(expectUser)).WithArgs(user.Email).WillReturnRows(userRows)

	res, err := s.repository.GetByEmail(context.Background(), "albert@robb@email.com")
	require.NoError(s.T(), err)
	require.Nil(s.T(), deep.Equal(user, *res))
}

func (s *UserRepositorySuite) TestUserRepository_GetAll() {

	user := models.User{
		FirstName: "Albert",
		LastName:  "Robb",
		Email:     "albert@robb@email.com",
		Products:  nil,
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

	const expectUser = "SELECT * FROM `users` WHERE `users`.`deleted_at` IS NULL ORDER BY created_at asc LIMIT 5"
	s.mock.ExpectQuery(regexp.QuoteMeta(expectUser)).WillReturnRows(userRows)

	const expectCount = "SELECT count(*) FROM `users` WHERE `users`.`deleted_at` IS NULL"
	s.mock.ExpectQuery(regexp.QuoteMeta(expectCount)).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	pagination := app.Pagination{
		Limit: 5,
		Page:  1,
		Sort:  "created_at asc",
	}
	var count int64
	res, err := s.repository.GetAllUser(context.Background(), pagination, &count)
	require.NoError(s.T(), err)
	require.Nil(s.T(), deep.Equal([]*models.User{&user}, res))
}

func (s *UserRepositorySuite) TestUserRepository_IsExistingEmail() {

	existRows := sqlmock.NewRows([]string{
		"count(*) > 0",
	}).AddRow(
		true,
	)

	const expectUser = "SELECT count(*) > 0 FROM `users` WHERE email = ? AND `users`.`deleted_at` IS NULL"
	s.mock.ExpectQuery(regexp.QuoteMeta(expectUser)).WithArgs("albert@robb@email.com").WillReturnRows(existRows)

	res, err := s.repository.IsExistingEmail(context.Background(), "albert@robb@email.com")
	require.NoError(s.T(), err)
	require.Nil(s.T(), deep.Equal(true, *res))
}

func (s *UserRepositorySuite) TestUserRepository_AddUser() {

	bv := []byte("password")
	hasher := sha256.New()
	hasher.Write(bv)

	user := models.User{
		FirstName: "Albert",
		LastName:  "Robb",
		Email:     "albert@robb@email.com",
		Password:  hasher.Sum(nil),
	}

	s.mock.ExpectBegin()
	s.mock.ExpectExec("INSERT INTO `users`").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), user.FirstName, user.LastName, user.Email, user.Password).WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()

	_, err := s.repository.AddUser(context.Background(), user)
	require.NoError(s.T(), err)
}
