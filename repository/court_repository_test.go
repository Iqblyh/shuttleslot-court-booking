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

var mockCourt = model.Court{
	Id:        "1",
	Name:      "field 1",
	Price:     30000,
	CreatedAt: time.Time{},
	UpdatedAt: time.Time{},
}

type CourtRepositoryTestSuite struct {
	suite.Suite
	mockDb  *sql.DB
	mockSql sqlmock.Sqlmock
	repo    CourtRepository
}

func (suite *CourtRepositoryTestSuite) SetupTest() {
	db, mock, err := sqlmock.New()
	assert.NoError(suite.T(), err)
	suite.mockDb = db
	suite.mockSql = mock
	suite.repo = NewCourtRepository(suite.mockDb)
}

func TestCourtRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(CourtRepositoryTestSuite))
}

func (suite *CourtRepositoryTestSuite) TestCreateCourt_Success() {
	suite.mockSql.ExpectQuery("INSERT INTO courts").
		WithArgs(mockCourt.Name, mockCourt.Price).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "price", "created_at", "updated_at"}).
			AddRow(mockCourt.Id, mockCourt.Name, mockCourt.Price, mockCourt.CreatedAt, mockCourt.UpdatedAt))

	actual, err := suite.repo.Create(mockCourt)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), mockCourt, actual)
}

func (suite *CourtRepositoryTestSuite) TestCreateCourt_Failed() {
	suite.mockSql.ExpectQuery("INSERT INTO courts").
		WithArgs(mockCourt.Name, mockCourt.Price).
		WillReturnError(errors.New("insert failed"))

	_, err := suite.repo.Create(mockCourt)
	assert.Error(suite.T(), err)
}

func (suite *CourtRepositoryTestSuite) TestFindAll_Success() {
	page := 1
	size := 10
	offset := (page - 1) * size

	suite.mockSql.ExpectQuery("SELECT id, name, price, created_at, updated_at ").
		WithArgs(size, offset).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "price", "created_at", "updated_at"}).
			AddRow(mockCourt.Id, mockCourt.Name, mockCourt.Price, mockCourt.CreatedAt, mockCourt.UpdatedAt))

	actual, _, err := suite.repo.FindAll(page, size)

	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), actual, 1)
	assert.Equal(suite.T(), mockCourt, actual[0])

}

func (suite *CourtRepositoryTestSuite) TestFindAll_Failed() {
	page := 1
	size := 10
	offset := (page - 1) * size

	suite.mockSql.ExpectQuery("SELECT id, name, price, created_at, updated_at ").
		WithArgs(page, size, offset).
		WillReturnError(fmt.Errorf("database error"))

	actual, paginate, err := suite.repo.FindAll(page, size)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), []model.Court{}, actual)
	assert.Equal(suite.T(), dto.Paginate{}, paginate)
}

func (suite *CourtRepositoryTestSuite) TestFindAll_ScanError() {
	page := 1
	size := 10
	offset := (page - 1) * size

	suite.mockSql.ExpectQuery(regexp.QuoteMeta("SELECT id, name, price, created_at, updated_at FROM courts LIMIT $1 OFFSET $2")).
		WithArgs(size, offset).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "price", "created_at", "updated_at"}).
			AddRow(mockCourt.Id, mockCourt.Name, "invalid_price", mockCourt.CreatedAt, mockCourt.UpdatedAt))

	actual, paginate, err := suite.repo.FindAll(page, size)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), []model.Court{}, actual)
	assert.Equal(suite.T(), dto.Paginate{}, paginate)
}

func (suite *CourtRepositoryTestSuite) TestFindById_Success() {
	suite.mockSql.ExpectQuery("SELECT").
		WithArgs(mockCourt.Id).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "price", "created_at", "updated_at"}).
			AddRow(mockCourt.Id, mockCourt.Name, mockCourt.Price, mockCourt.CreatedAt, mockCourt.UpdatedAt))

	actual, err := suite.repo.FindById(mockCourt.Id)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), mockCourt, actual)
}

func (suite *CourtRepositoryTestSuite) TestFindById_Failed() {
	suite.mockSql.ExpectQuery("SELECT").
		WithArgs(mockCourt.Id).
		WillReturnError(errors.New("select failed"))

	_, err := suite.repo.FindById(mockCourt.Id)
	assert.Error(suite.T(), err)
}

func (suite *CourtRepositoryTestSuite) TestUpdateCourt_Success() {
	mockUpdatedAt := time.Now()
	suite.mockSql.ExpectQuery("UPDATE courts SET ").
		WithArgs(mockCourt.Name, mockCourt.Price, mockUpdatedAt, mockCourt.Id).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "price", "created_at", "updated_at"}).
			AddRow(mockCourt.Id, mockCourt.Name, mockCourt.Price, mockCourt.CreatedAt, mockUpdatedAt))

	mockCourt.UpdatedAt = mockUpdatedAt

	actual, err := suite.repo.Update(mockCourt.Id, mockCourt)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), mockCourt, actual)
}

func (suite *CourtRepositoryTestSuite) TestUpdateCourt_Failed() {
	suite.mockSql.ExpectQuery("UPDATE court").
		WithArgs(mockCourt.Name, mockCourt.Price, mockCourt.UpdatedAt, mockCourt.Id).
		WillReturnError(errors.New("update failed"))

	_, err := suite.repo.Update(mockCourt.Id, model.Court{})

	assert.Error(suite.T(), err)
}

func (suite *CourtRepositoryTestSuite) TestDeleteUser_Success() {
	suite.mockSql.ExpectExec("DELETE FROM court").
		WithArgs(mockCourt.Id).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := suite.repo.Deleted(mockCourt.Id)
	assert.NoError(suite.T(), err)

}

func (suite *CourtRepositoryTestSuite) TestDeleteUser_Failed() {
	suite.mockSql.ExpectExec("DELETE FROM court").
		WithArgs(mockCourt.Id).
		WillReturnError(errors.New("delete failed"))

	err := suite.repo.Deleted(mockCourt.Id)
	assert.Error(suite.T(), err)
	assert.EqualError(suite.T(), err, "delete failed")

}
