package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"team2/shuttleslot/model"
	"team2/shuttleslot/model/dto"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var mockUser = model.User{
	Id:          "1",
	Name:        "lala",
	PhoneNumber: "123",
	Email:       "@mail",
	Username:    "lalalele",
	Password:    "passwrod",
	Point:       0,
	Role:        "customer",
	CreatedAt:   time.Now(),
	UpdatedAt:   time.Now(),
}

type UserRepositoryTestSuite struct {
	suite.Suite
	mockDb  *sql.DB
	mockSql sqlmock.Sqlmock
	repo    UserRepository
}

func (suite *UserRepositoryTestSuite) SetupTest() {
	db, mock, err := sqlmock.New()
	assert.NoError(suite.T(), err)
	suite.mockDb = db
	suite.mockSql = mock
	suite.repo = NewUserRepository(suite.mockDb)
}

func TestUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}

func (suite *UserRepositoryTestSuite) TestCreateCustomer_Success() {
	suite.mockSql.ExpectQuery("INSERT INTO users").
		WithArgs(mockUser.Name, mockUser.PhoneNumber, mockUser.Email, mockUser.Username, mockUser.Password, mockUser.Role).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "phone_number", "email", "username", "password", "points", "role", "created_at", "updated_at"}).
			AddRow(mockUser.Id, mockUser.Name, mockUser.PhoneNumber, mockUser.Email, mockUser.Username, mockUser.Password, mockUser.Point, mockUser.Role, mockUser.CreatedAt, mockUser.UpdatedAt))

	actual, err := suite.repo.CreateCustomer(mockUser)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), mockUser, actual)
}

func (suite *UserRepositoryTestSuite) TestCreateCustomer_Failed() {
	suite.mockSql.ExpectQuery("INSERT INTO users").
		WithArgs(mockUser.Name, mockUser.Email, mockUser.Username, mockUser.Password, mockUser.PhoneNumber, mockUser.Role).
		WillReturnError(errors.New("insert failed"))

	_, err := suite.repo.CreateCustomer(mockUser)
	assert.Error(suite.T(), err)
}

func (suite *UserRepositoryTestSuite) TestCreateEmployee_Success() {
	suite.mockSql.ExpectQuery("INSERT INTO users").
		WithArgs(mockUser.Name, mockUser.PhoneNumber, mockUser.Email, mockUser.Username, mockUser.Password, mockUser.Role).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "phone_number", "email", "username", "password", "points", "role", "created_at", "updated_at"}).
			AddRow(mockUser.Id, mockUser.Name, mockUser.PhoneNumber, mockUser.Email, mockUser.Username, mockUser.Password, mockUser.Point, mockUser.Role, mockUser.CreatedAt, mockUser.UpdatedAt))

	actual, err := suite.repo.CreateEmployee(mockUser)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), mockUser, actual)
}

func (suite *UserRepositoryTestSuite) TestCreateEmployee_Failed() {
	suite.mockSql.ExpectQuery("INSERT INTO users").
		WithArgs(mockUser.Name, mockUser.Email, mockUser.Username, mockUser.Password, mockUser.PhoneNumber, mockUser.Role).
		WillReturnError(errors.New("insert failed"))

	_, err := suite.repo.CreateEmployee(mockUser)
	assert.Error(suite.T(), err)
}

func (suite *UserRepositoryTestSuite) TestCreateAdmin_Success() {
	suite.mockSql.ExpectQuery("INSERT INTO users").
		WithArgs(mockUser.Name, mockUser.PhoneNumber, mockUser.Email, mockUser.Username, mockUser.Password, mockUser.Role).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "phone_number", "email", "username", "password", "points", "role", "created_at", "updated_at"}).
			AddRow(mockUser.Id, mockUser.Name, mockUser.PhoneNumber, mockUser.Email, mockUser.Username, mockUser.Password, mockUser.Point, mockUser.Role, mockUser.CreatedAt, mockUser.UpdatedAt))

	actual, err := suite.repo.CreateAdmin(mockUser)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), mockUser, actual)
}

func (suite *UserRepositoryTestSuite) TestCreateAdmin_Failed() {
	suite.mockSql.ExpectQuery("INSERT INTO users").
		WithArgs(mockUser.Name, mockUser.Email, mockUser.Username, mockUser.Password, mockUser.PhoneNumber, mockUser.Role).
		WillReturnError(errors.New("insert failed"))

	_, err := suite.repo.CreateAdmin(mockUser)
	assert.Error(suite.T(), err)
}

func (suite *UserRepositoryTestSuite) TestFindUserByUsername_Success() {
	suite.mockSql.ExpectQuery("SELECT").
		WithArgs(mockUser.Username).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "phone_number", "email", "username", "password", "points", "role", "created_at", "updated_at"}).
			AddRow(mockUser.Id, mockUser.Name, mockUser.PhoneNumber, mockUser.Email, mockUser.Username, mockUser.Password, mockUser.Point, mockUser.Role, mockUser.CreatedAt, mockUser.UpdatedAt))

	actual, err := suite.repo.FindUserByUsername(mockUser.Username)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), mockUser, actual)
}

func (suite *UserRepositoryTestSuite) TestFindUserByUsername_Failed() {
	suite.mockSql.ExpectQuery("SELECT").
		WithArgs(mockUser.Username).
		WillReturnError(errors.New("select failed"))

	_, err := suite.repo.FindUserByUsername(mockUser.Username)
	assert.Error(suite.T(), err)
}

func (suite *UserRepositoryTestSuite) TestFindUserById_Success() {
	suite.mockSql.ExpectQuery("SELECT").
		WithArgs(mockUser.Id).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "phone_number", "email", "username", "password", "points", "role", "created_at", "updated_at"}).
			AddRow(mockUser.Id, mockUser.Name, mockUser.PhoneNumber, mockUser.Email, mockUser.Username, mockUser.Password, mockUser.Point, mockUser.Role, mockUser.CreatedAt, mockUser.UpdatedAt))

	actual, err := suite.repo.FindUserById(mockUser.Id)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), mockUser, actual)
}

func (suite *UserRepositoryTestSuite) TestFindUserById_Failed() {
	suite.mockSql.ExpectQuery("SELECT id, name, phone_number, email, username, password, points, role").
		WithArgs(mockUser.Id).
		WillReturnError(errors.New("select failed"))

	_, err := suite.repo.FindUserById(mockUser.Id)
	assert.Error(suite.T(), err)
}

func (suite *UserRepositoryTestSuite) TestFindUserByRole_Success() {
	page := 1
	size := 10
	offset := (page - 1) * size

	suite.mockSql.ExpectQuery("SELECT id, name, phone_number, email, username, password, points, role, created_at, updated_at").
		WithArgs(mockUser.Role, size, offset).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "phone_number", "email", "username", "password", "points", "role", "created_at", "updated_at"}).
			AddRow(mockUser.Id, mockUser.Name, mockUser.PhoneNumber, mockUser.Email, mockUser.Username, mockUser.Password, mockUser.Point, mockUser.Role, mockUser.CreatedAt, mockUser.UpdatedAt))

	actual, paginate, err := suite.repo.FindUserByRole(mockUser.Role, page, size)

	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), actual, 1)
	assert.Equal(suite.T(), mockUser, actual[0])

	assert.Equal(suite.T(), page, paginate.Page)
	assert.Equal(suite.T(), size, paginate.Size)
	assert.Equal(suite.T(), 1, paginate.TotalRows)
	assert.Equal(suite.T(), 1, paginate.TotalPages)
}

func (suite *UserRepositoryTestSuite) TestFindUserByRole_Failure() {
	page := 1
	size := 10
	offset := (page - 1) * size

	suite.mockSql.ExpectQuery("SELECT id, name, phone_number, email, username, password, points, role").
		WithArgs(mockUser.Role, size, offset).
		WillReturnError(fmt.Errorf("database error"))

	actual, paginate, err := suite.repo.FindUserByRole(mockUser.Role, page, size)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), []model.User{}, actual)
	assert.Equal(suite.T(), dto.Paginate{}, paginate)
}

func (suite *UserRepositoryTestSuite) TestFindUserByRole_ScanError() {
	page := 1
	size := 10
	offset := (page - 1) * size

	suite.mockSql.ExpectQuery(regexp.QuoteMeta("SELECT id, name, phone_number, email, username, password, points, role, created_at, updated_at FROM users WHERE role = $1 LIMIT $2 OFFSET $3")).
		WithArgs(mockUser.Role, size, offset).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "phone_number", "email", "username", "password", "points", "role", "created_at", "updated_at"}).
			AddRow(mockUser.Id, mockUser.Name, mockUser.PhoneNumber, mockUser.Email, mockUser.Username, mockUser.Password, "invalid_point", "invalid_role", mockUser.CreatedAt, mockUser.UpdatedAt))

	actual, paginate, err := suite.repo.FindUserByRole(mockUser.Role, page, size)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), []model.User{}, actual)
	assert.Equal(suite.T(), dto.Paginate{}, paginate)
}

func (suite *UserRepositoryTestSuite) TestUpdateUser_Success() {
	mockUpdatedAt := time.Now()

	suite.mockSql.ExpectQuery("UPDATE users SET ").
		WithArgs(mockUser.Name, mockUser.PhoneNumber, mockUser.Email, mockUser.Username, mockUser.Password, mockUpdatedAt, mockUser.Id).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "phone_number", "email", "username", "password", "points", "role", "created_at", "updated_at"}).
			AddRow(mockUser.Id, mockUser.Name, mockUser.PhoneNumber, mockUser.Email, mockUser.Username, mockUser.Password, mockUser.Point, mockUser.Role, mockUser.CreatedAt, mockUpdatedAt))

	mockUser.UpdatedAt = mockUpdatedAt

	actual, err := suite.repo.UpdateUser(mockUser.Id, mockUser)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), mockUser, actual)
}

func (suite *UserRepositoryTestSuite) TestUpdateUser_Failed() {
	mockUpdatedAt := time.Now()
	suite.mockSql.ExpectQuery("UPDATE users SET ").
		WithArgs(mockUser.Name, mockUser.PhoneNumber, mockUser.Email, mockUser.Username, mockUser.Password, mockUpdatedAt, mockUser.Id).
		WillReturnError(errors.New("update failed"))

	_, err := suite.repo.UpdateUser(mockUser.Id, mockUser)

	assert.Error(suite.T(), err)
}

func (suite *UserRepositoryTestSuite) TestDeleteUser_Success() {
	suite.mockSql.ExpectExec("DELETE FROM users WHERE id = ?").
		WithArgs(mockUser.Id).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := suite.repo.DeleteUser(mockUser.Id)
	assert.NoError(suite.T(), err)
}

func (suite *UserRepositoryTestSuite) TestDeleteUser_Failed() {
	suite.mockSql.ExpectExec("DELETE FROM users WHERE id = ?").
		WithArgs(mockUser.Id).
		WillReturnError(errors.New("delete failed"))

	err := suite.repo.DeleteUser(mockUser.Id)
	assert.Error(suite.T(), err)
	assert.EqualError(suite.T(), err, "delete failed")
}
