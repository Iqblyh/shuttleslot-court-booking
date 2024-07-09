package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	mock "team2/shuttleslot/mock"
	servicemock "team2/shuttleslot/mock/service_mock"
	"team2/shuttleslot/model"
	"team2/shuttleslot/model/dto"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var payloadCourt = model.Court{
	Id:    "1",
	Name:  "Field 1",
	Price: 50000,
}

type CourtControllerTestSuite struct {
	suite.Suite
	courtServiceMock *servicemock.CourtServiceMock
	middlewareMock   *mock.AuthMiddlewareMock
	rg               *gin.RouterGroup
	courtController  *CourtController
}

func (suite *CourtControllerTestSuite) SetupTest() {
	suite.courtServiceMock = new(servicemock.CourtServiceMock)
	rg := gin.Default()
	suite.rg = rg.Group("/api/v1/courts")
	suite.middlewareMock = new(mock.AuthMiddlewareMock)
	suite.courtController = NewCourtController(suite.courtServiceMock, suite.middlewareMock, suite.rg)
	suite.courtController.Route()
}

func TestCourtControllerTestSuite(t *testing.T) {
	suite.Run(t, new(CourtControllerTestSuite))
}
func (suite *CourtControllerTestSuite) TestCreateCourt_Succes() {
	mockPayloadjson, err := json.Marshal(payloadCourt)
	assert.NoError(suite.T(), err)

	record := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, "/api/v1/courts", bytes.NewBuffer(mockPayloadjson))
	assert.NoError(suite.T(), err)

	req.Header.Set("Authorization", "Bearer "+token)
	ctx, _ := gin.CreateTestContext(record)
	ctx.Request = req

	suite.courtServiceMock.On("CreateCourt", payloadCourt).Return(payloadCourt, nil)
	suite.courtController.CreateCourtHandler(ctx)
	assert.Equal(suite.T(), http.StatusOK, record.Code)
}

func (suite *CourtControllerTestSuite) TestCreateCourt_FailedBinding() {
	record := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/", nil)

	req.Header.Set("Authorization", "Bearer "+token)
	ctx, _ := gin.CreateTestContext(record)
	ctx.Request = req
	suite.courtController.CreateCourtHandler(ctx)
	assert.Equal(suite.T(), http.StatusBadRequest, record.Code)
}
func (suite *CourtControllerTestSuite) TestCreateCourt_Failed() {
	mockPayloadjson, err := json.Marshal(payloadCourt)
	assert.NoError(suite.T(), err)

	record := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, "/api/v1/courts", bytes.NewBuffer(mockPayloadjson))
	assert.NoError(suite.T(), err)

	req.Header.Set("Authorization", "Bearer "+token)
	ctx, _ := gin.CreateTestContext(record)
	ctx.Request = req

	suite.courtServiceMock.On("CreateCourt", payloadCourt).Return(model.Court{}, errors.New("error"))
	suite.courtController.CreateCourtHandler(ctx)
	assert.Equal(suite.T(), http.StatusInternalServerError, record.Code)
}

func (suite *CourtControllerTestSuite) TestFindAllCourtsHandler_Success() {
	Paginate := dto.Paginate{
		Page:       1,
		Size:       10,
		TotalRows:  0,
		TotalPages: 0,
	}

	payloadCourt2 := []model.Court{
		{
			Id:    "1",
			Name:  "field 2",
			Price: 200000,
		},
	}
	mockPayloadjson, err := json.Marshal(payloadCourt)
	assert.NoError(suite.T(), err)
	record := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/api/v1/courts", bytes.NewBuffer(mockPayloadjson))
	assert.NoError(suite.T(), err)

	req.Header.Set("Authorization", "Bearer "+token)
	ctx, router := gin.CreateTestContext(record)
	router.GET("/api/v1/courts", suite.courtController.FindAllCourtsHandler)
	ctx.Request = req

	suite.courtServiceMock.On("FindAllCourts", Paginate.Page, Paginate.Size).Return(payloadCourt2, Paginate, nil)
	router.ServeHTTP(record, req)
	assert.Equal(suite.T(), http.StatusOK, record.Code)
	suite.courtServiceMock.AssertExpectations(suite.T())
}

func (suite *CourtControllerTestSuite) TestFindAllCourtsHandler_Failed() {
	Paginate := dto.Paginate{
		Page: 1,
		Size: 10,
	}

	mockPayloadjson, err := json.Marshal(payloadCourt)
	assert.NoError(suite.T(), err)
	record := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/api/v1/courts", bytes.NewBuffer(mockPayloadjson))
	assert.NoError(suite.T(), err)

	req.Header.Set("Authorization", "Bearer "+token)
	ctx, router := gin.CreateTestContext(record)
	router.GET("/api/v1/courts", suite.courtController.FindAllCourtsHandler)
	ctx.Request = req

	suite.courtServiceMock.On("FindAllCourts", Paginate.Page, Paginate.Size).Return([]model.Court{}, Paginate, errors.New("not found"))

	router.ServeHTTP(record, req)
	assert.Equal(suite.T(), http.StatusNotFound, record.Code)
	suite.courtServiceMock.AssertExpectations(suite.T())
}
func (suite *CourtControllerTestSuite) TestFindAllCourtsHandler_InvalidPageSize() {
	record := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/api/v1/courts"+"?page=invalid&size=invalid", nil)
	assert.NoError(suite.T(), err)

	req.Header.Set("Authorization", "Bearer "+token)
	ctx, router := gin.CreateTestContext(record)
	router.GET("/api/v1/courts", suite.courtController.FindAllCourtsHandler)
	ctx.Request = req

	router.ServeHTTP(record, req)
	assert.Equal(suite.T(), http.StatusBadRequest, record.Code)
}

func (suite *CourtControllerTestSuite) TestFindCourtByIdHandler_Success() {
	record := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/api/v1/court/id/"+payloadCourt.Id, nil)
	assert.NoError(suite.T(), err)

	req.Header.Set("Authorization", "Bearer "+token)
	ctx, router := gin.CreateTestContext(record)
	router.GET("/api/v1/court/id/:id", suite.courtController.FindCourtByIdHandler)
	ctx.Request = req

	suite.courtServiceMock.On("FindCourtById", payloadCourt.Id).Return(payloadCourt, nil)
	router.ServeHTTP(record, req)
	assert.Equal(suite.T(), http.StatusOK, record.Code)
}
func (suite *CourtControllerTestSuite) TestFindCourtByIdHandler_Failed() {
	record := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/api/v1/court/id/"+payloadCourt.Id, nil)
	assert.NoError(suite.T(), err)

	req.Header.Set("Authorization", "Bearer "+token)
	ctx, router := gin.CreateTestContext(record)
	router.GET("/api/v1/court/id/:id", suite.courtController.FindCourtByIdHandler)
	ctx.Request = req

	suite.courtServiceMock.On("FindCourtById", payloadCourt.Id).Return(model.Court{}, errors.New("not found"))

	router.ServeHTTP(record, req)
	assert.Equal(suite.T(), http.StatusNotFound, record.Code)
	suite.courtServiceMock.AssertExpectations(suite.T())
}

func (suite *CourtControllerTestSuite) TestDeleteCourtHandler_Success() {
	suite.courtServiceMock.On("DeleteCourt", "1").Return(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/court/1", nil)
	req.Header.Set("Content-Type", "application/json")

	ctx, _ := gin.CreateTestContext(w)
	ctx.Params = gin.Params{{Key: "id", Value: "1"}}
	ctx.Request = req

	suite.courtController.DeleteCourtHandler(ctx)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func (suite *CourtControllerTestSuite) TestDeleteCourtHandler_Failed() {
	suite.courtServiceMock.On("DeleteCourt", "1").Return(errors.New("user not found"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/court/1", nil)
	req.Header.Set("Content-Type", "application/json")

	ctx, _ := gin.CreateTestContext(w)
	ctx.Params = gin.Params{{Key: "id", Value: "1"}}
	ctx.Request = req

	suite.courtController.DeleteCourtHandler(ctx)
	assert.Equal(suite.T(), http.StatusNotFound, w.Code)
}

func (suite *CourtControllerTestSuite) TestUpdateCourtHandler_Success() {

	suite.courtServiceMock.On("UpdateCourt", "1", payloadCourt).Return(payloadCourt, nil)

	w := httptest.NewRecorder()
	body, _ := json.Marshal(payloadCourt)
	req, _ := http.NewRequest(http.MethodPut, "/api/v1/court/1", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+token)

	ctx, _ := gin.CreateTestContext(w)
	ctx.Params = gin.Params{{Key: "id", Value: "1"}}
	ctx.Request = req

	suite.courtController.UpdateCourtHandler(ctx)
	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func (suite *CourtControllerTestSuite) TestUpdateCourtHandler_FailedBindJSON() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/api/v1/court/1", bytes.NewBuffer([]byte("{invalid json}")))
	req.Header.Set("Content-Type", "application/json")

	ctx, _ := gin.CreateTestContext(w)
	ctx.Params = gin.Params{{Key: "id", Value: "1"}}
	ctx.Request = req

	suite.courtController.UpdateCourtHandler(ctx)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}

func (suite *CourtControllerTestSuite) TestUpdateCourtHandler_FailedUpdateError() {
	suite.courtServiceMock.On("UpdateCourt", "1", model.Court{}).Return(model.Court{}, errors.New("update error"))

	w := httptest.NewRecorder()
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPut, "/api/v1/court/1", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+token)

	ctx, _ := gin.CreateTestContext(w)
	ctx.Params = gin.Params{{Key: "id", Value: "1"}}
	ctx.Request = req

	suite.courtController.UpdateCourtHandler(ctx)
	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
}
