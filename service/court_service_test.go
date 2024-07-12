package service

import (
	"errors"
	repomock "team2/shuttleslot/mock/repo_mock"
	"team2/shuttleslot/model"
	"team2/shuttleslot/model/dto"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

var mockCourt = model.Court{
	Id:    "1",
	Name:  "Court 1",
	Price: 100,
}
var updatedCourt = model.Court{
	Id:    "court_id",
	Name:  "Updated Court",
	Price: 60,
}

var payloadCourt = model.Court{
	Name:  "Updated Court",
	Price: 60,
}

type CourtServiceTestSuite struct {
	suite.Suite
	repoCourtMock *repomock.CourtRepositoryMock
	cS            CourtService
}

func (suite *CourtServiceTestSuite) SetupTest() {
	suite.repoCourtMock = new(repomock.CourtRepositoryMock)
	suite.cS = NewCourtService(suite.repoCourtMock)
}

func TestCourtServiceTestSuite(t *testing.T) {
	suite.Run(t, new(CourtServiceTestSuite))
}

func (suite *CourtServiceTestSuite) TestCreateCourt_Success() {
	suite.repoCourtMock.On("Create", mockCourt).Return(mockCourt, nil)

	court, err := suite.cS.CreateCourt(mockCourt)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), mockCourt, court)
}

func (suite *CourtServiceTestSuite) TestCreateCourt_Fail() {
	suite.repoCourtMock.On("Create", mockCourt).Return(model.Court{}, errors.New("error creating court"))

	court, err := suite.cS.CreateCourt(mockCourt)

	assert.Error(suite.T(), err)
	assert.EqualError(suite.T(), err, "error creating court")
	assert.Equal(suite.T(), model.Court{}, court)
}

func (suite *CourtServiceTestSuite) TestFindAllCourts_Success() {
	page := 1
	size := 10
	mockCourts := []model.Court{mockCourt}
	mockPaginate := dto.Paginate{
		Page:       page,
		Size:       size,
		TotalRows:  1,
		TotalPages: 1,
	}

	suite.repoCourtMock.On("FindAll", page, size).Return(mockCourts, mockPaginate, nil)

	courts, paginate, err := suite.cS.FindAllCourts(page, size)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), mockCourts, courts)
	assert.Equal(suite.T(), mockPaginate, paginate)
}

func (suite *CourtServiceTestSuite) TestFindAllCourts_Fail() {
	page := 1
	size := 10
	mockError := errors.New("error finding courts")

	suite.repoCourtMock.On("FindAll", page, size).Return([]model.Court{}, dto.Paginate{}, mockError)

	_, paginate, err := suite.cS.FindAllCourts(page, size)

	assert.Error(suite.T(), err)
	assert.EqualError(suite.T(), err, "error finding courts")
	assert.Nil(suite.T(), nil)
	assert.Equal(suite.T(), dto.Paginate{}, paginate)
}

func (suite *CourtServiceTestSuite) TestFindCourtById_Success() {
	suite.repoCourtMock.On("FindById", mockCourt.Id).Return(mockCourt, nil)

	court, err := suite.cS.FindCourtById(mockCourt.Id)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), mockCourt, court)
}

func (suite *CourtServiceTestSuite) TestFindCourtById_Fail() {
	suite.repoCourtMock.On("FindById", mockCourt.Id).Return(model.Court{}, errors.New("error"))

	court, err := suite.cS.FindCourtById(mockCourt.Id)

	assert.Error(suite.T(), err)
	assert.EqualError(suite.T(), err, "error")
	assert.Equal(suite.T(), model.Court{}, court)
}

func (suite *CourtServiceTestSuite) TestUpdateCourt_Success() {
	suite.repoCourtMock.On("FindById", "court_id").Return(mockCourt, nil)
	suite.repoCourtMock.On("Update", "court_id", payloadCourt).Return(updatedCourt, nil)

	_, err := suite.cS.UpdateCourt("court_id", payloadCourt)

	assert.NoError(suite.T(), err)
	suite.repoCourtMock.AssertExpectations(suite.T())
}

func (suite *CourtServiceTestSuite) TestUpdateCourt_FailFindById() {
	suite.repoCourtMock.On("FindById", "non_existing_id").Return(model.Court{}, errors.New("court not found"))

	_, err := suite.cS.UpdateCourt("non_existing_id", model.Court{Name: "Updated Court", Price: 60})

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "court not found", err.Error())
	suite.repoCourtMock.AssertExpectations(suite.T())
}

func (suite *CourtServiceTestSuite) TestUpdateCourt_FailUpdate() {
	suite.repoCourtMock.On("FindById", "court_id").Return(mockCourt, nil)
	suite.repoCourtMock.On("Update", "court_id", payloadCourt).Return(model.Court{}, errors.New("failed to update court"))

	_, err := suite.cS.UpdateCourt("court_id", payloadCourt)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "failed to update court", err.Error())
	suite.repoCourtMock.AssertExpectations(suite.T())
}

func (suite *CourtServiceTestSuite) TestUpdateCourt_EmptyPayloadFields() {
	mockCourt := model.Court{
		Id:    "court_id",
		Name:  "Old Court",
		Price: 50,
	}
	updatedCourt := model.Court{
		Id:    "court_id",
		Name:  "Old Court",
		Price: 60,
	}
	payload := model.Court{
		Price: 0,
	}

	suite.repoCourtMock.On("FindById", "court_id").Return(mockCourt, nil)

	suite.repoCourtMock.On("Update", "court_id", mock.MatchedBy(func(c model.Court) bool {
		return c.Name == "Old Court" && c.Price == 50
	})).Return(updatedCourt, nil)

	_, err := suite.cS.UpdateCourt("court_id", payload)

	assert.NoError(suite.T(), err)
	suite.repoCourtMock.AssertExpectations(suite.T())
}

func (suite *CourtServiceTestSuite) TestDeleteCourt_Success() {
	suite.repoCourtMock.On("FindById", mockCourt.Id).Return(mockCourt, nil)
	suite.repoCourtMock.On("Deleted", mockCourt.Id).Return(nil)

	err := suite.cS.DeleteCourt(mockCourt.Id)

	assert.NoError(suite.T(), err)
}

func (suite *CourtServiceTestSuite) TestDeleteCourt_Fail() {
	suite.repoCourtMock.On("FindById", mockCourt.Id).Return(model.Court{}, errors.New("court not found"))

	err := suite.cS.DeleteCourt(mockCourt.Id)

	assert.Error(suite.T(), err)
	assert.EqualError(suite.T(), err, "court not found")
}

func (suite *CourtServiceTestSuite) TestDeleteCourt_ErrorDeleting() {
	suite.repoCourtMock.On("FindById", mockCourt.Id).Return(model.Court{}, nil)
	suite.repoCourtMock.On("Deleted", mockCourt.Id).Return(errors.New("error deleting court"))

	err := suite.cS.DeleteCourt(mockCourt.Id)

	assert.Error(suite.T(), err)
	assert.EqualError(suite.T(), err, "error deleting court")
	suite.repoCourtMock.AssertExpectations(suite.T())
}
