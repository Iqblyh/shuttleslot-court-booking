package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	mock "team2/shuttleslot/mock"
	servicemock "team2/shuttleslot/mock/service_mock"
	"team2/shuttleslot/model"
	"team2/shuttleslot/model/dto"
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var payload = dto.LoginRequest{
	Username: "lalalele",
	Password: "password",
}

var payloadUser = model.User{
	Id:          "1",
	Name:        "lala",
	PhoneNumber: "08989",
	Email:       "lala@mail",
	Username:    "lalalele",
	Password:    "password",
	Point:       0,
	Role:        "admin",
}
var token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJFbmlnbWFDYW1wIiwiZXhwIjoxNzE5ODAzNjgwLCJpYXQiOjE3MTk4MDAwODAsInVzZXJJZCI6ImJhNDgxNGUwLTZkODMtNDM0Mi05ZWExLTYwODVjOGNjNWJmMSIsInJvbGUiOiJhZG1pbiJ9.oEzYcGrNbBfOW5zd11doq-mtZixdMCLaA9HkkTO-PEk"

type UserControllerTestSuite struct {
	suite.Suite
	userServiceMock *servicemock.UserServiceMock
	middlewareMock  *mock.AuthMiddlewareMock
	rg              *gin.RouterGroup
	userController  *UserController
}

func (suite *UserControllerTestSuite) SetupTest() {
	suite.userServiceMock = new(servicemock.UserServiceMock)
	rg := gin.Default()
	suite.rg = rg.Group("/api/v1/users")
	suite.middlewareMock = new(mock.AuthMiddlewareMock)
	suite.userController = NewUserController(suite.userServiceMock, suite.middlewareMock, suite.rg)
	suite.userController.Route()
}

func TestUserControllerTestSuite(t *testing.T) {
	suite.Run(t, new(UserControllerTestSuite))
}

func (suite *UserControllerTestSuite) TestLogin_Success() {
	response := dto.LoginResponse{Token: "token123"}
	mockPayloadjson, err := json.Marshal(payload)
	assert.NoError(suite.T(), err)

	record := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, "/api/v1/users/login", bytes.NewBuffer(mockPayloadjson))
	assert.NoError(suite.T(), err)
	req.Header.Set("Content-Type", "application/json")

	ctx, _ := gin.CreateTestContext(record)
	ctx.Request = req

	suite.userServiceMock.On("Login", payload).Return(response, nil)
	suite.userController.LoginHandler(ctx)
	assert.Equal(suite.T(), http.StatusOK, record.Code)
}

func (suite *UserControllerTestSuite) TestLogin_FailedBinding() {
	record := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/", nil)

	req.Header.Set("Authorization", "Bearer "+token)
	ctx, _ := gin.CreateTestContext(record)
	ctx.Request = req
	suite.userController.LoginHandler(ctx)
	assert.Equal(suite.T(), http.StatusBadRequest, record.Code)
}
func (suite *UserControllerTestSuite) TestLogin_Failed() {
	response := dto.LoginResponse{}
	mockPayloadjson, err := json.Marshal(payload)
	assert.NoError(suite.T(), err)

	record := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, "/api/v1/users/login", bytes.NewBuffer(mockPayloadjson))
	assert.NoError(suite.T(), err)

	req.Header.Set("Content-Type", "application/json")
	ctx, _ := gin.CreateTestContext(record)
	ctx.Request = req

	suite.userServiceMock.On("Login", payload).Return(response, errors.New("error"))
	suite.userController.LoginHandler(ctx)
	assert.Equal(suite.T(), http.StatusInternalServerError, record.Code)
}

func (suite *UserControllerTestSuite) TestCreateAdminHandler_Success() {
	mockPayloadjson, err := json.Marshal(payloadUser)
	assert.NoError(suite.T(), err)

	record := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, "/api/v1/users/admin/create", bytes.NewBuffer(mockPayloadjson))
	assert.NoError(suite.T(), err)

	req.Header.Set("Authorization", "Bearer "+token)
	ctx, _ := gin.CreateTestContext(record)
	ctx.Request = req

	suite.userServiceMock.On("CreateAdmin", payloadUser).Return(payloadUser, nil)
	suite.userController.CreateAdminHandler(ctx)
	assert.Equal(suite.T(), http.StatusOK, record.Code)
}

func (suite *UserControllerTestSuite) TestCreateAdminHandler_FailedBinding() {
	record := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/", nil)

	req.Header.Set("Authorization", "Bearer "+token)
	ctx, _ := gin.CreateTestContext(record)
	ctx.Request = req
	suite.userController.CreateAdminHandler(ctx)
	assert.Equal(suite.T(), http.StatusBadRequest, record.Code)
}
func (suite *UserControllerTestSuite) TestCreateAdminHandler_Failed() {
	mockPayloadjson, err := json.Marshal(payloadUser)
	assert.NoError(suite.T(), err)

	record := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, "/api/v1/users/admin/create", bytes.NewBuffer(mockPayloadjson))
	assert.NoError(suite.T(), err)

	req.Header.Set("Authorization", "Bearer "+token)
	ctx, _ := gin.CreateTestContext(record)
	ctx.Request = req

	suite.userServiceMock.On("CreateAdmin", payloadUser).Return(model.User{}, errors.New("error"))
	suite.userController.CreateAdminHandler(ctx)
	assert.Equal(suite.T(), http.StatusInternalServerError, record.Code)
}

func (suite *UserControllerTestSuite) TestCreateCustomerHandler_Success() {
	mockPayloadjson, err := json.Marshal(payloadUser)
	assert.NoError(suite.T(), err)

	record := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, "/api/v1/users/customer/create", bytes.NewBuffer(mockPayloadjson))
	assert.NoError(suite.T(), err)

	req.Header.Set("Content-Type", "application/json")
	ctx, _ := gin.CreateTestContext(record)
	ctx.Request = req

	suite.userServiceMock.On("CreateCustomer", payloadUser).Return(payloadUser, nil)
	suite.userController.CreateCustomerHandler(ctx)
	assert.Equal(suite.T(), http.StatusOK, record.Code)
}

func (suite *UserControllerTestSuite) TestCreateCustomerHandler_FailedBinding() {
	record := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/", nil)

	req.Header.Set("Authorization", "Bearer "+token)
	ctx, _ := gin.CreateTestContext(record)
	ctx.Request = req
	suite.userController.CreateCustomerHandler(ctx)
	assert.Equal(suite.T(), http.StatusBadRequest, record.Code)
}

func (suite *UserControllerTestSuite) TestCreateCustomerHandler_Failed() {
	mockPayloadjson, err := json.Marshal(payloadUser)
	assert.NoError(suite.T(), err)

	record := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, "/api/v1/users/customer/create", bytes.NewBuffer(mockPayloadjson))
	assert.NoError(suite.T(), err)

	req.Header.Set("Content-Type", "application/json")
	ctx, _ := gin.CreateTestContext(record)
	ctx.Request = req

	suite.userServiceMock.On("CreateCustomer", payloadUser).Return(model.User{}, errors.New("error"))
	suite.userController.CreateCustomerHandler(ctx)
	assert.Equal(suite.T(), http.StatusInternalServerError, record.Code)
}

func (suite *UserControllerTestSuite) TestCreateEmployeeHandler_Success() {
	mockPayloadjson, err := json.Marshal(payloadUser)
	assert.NoError(suite.T(), err)

	record := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, "/api/v1/users/employee/create", bytes.NewBuffer(mockPayloadjson))
	assert.NoError(suite.T(), err)

	req.Header.Set("Authorization", "Bearer "+token)
	ctx, _ := gin.CreateTestContext(record)
	ctx.Request = req

	suite.userServiceMock.On("CreateEmployee", payloadUser).Return(payloadUser, nil)
	suite.userController.CreateEmployeeHandler(ctx)
	assert.Equal(suite.T(), http.StatusOK, record.Code)
}

func (suite *UserControllerTestSuite) TestCreateEmployeeHandler_FailedBinding() {
	record := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/", nil)

	req.Header.Set("Authorization", "Bearer "+token)
	ctx, _ := gin.CreateTestContext(record)
	ctx.Request = req
	suite.userController.CreateEmployeeHandler(ctx)
	assert.Equal(suite.T(), http.StatusBadRequest, record.Code)
}
func (suite *UserControllerTestSuite) TestCreateEmployeeHandler_Failed() {
	mockPayloadjson, err := json.Marshal(payloadUser)
	assert.NoError(suite.T(), err)

	record := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, "/api/v1/users/employee/create", bytes.NewBuffer(mockPayloadjson))
	assert.NoError(suite.T(), err)

	req.Header.Set("Authorization", "Bearer "+token)
	ctx, _ := gin.CreateTestContext(record)
	ctx.Request = req

	suite.userServiceMock.On("CreateEmployee", payloadUser).Return(model.User{}, errors.New("error"))
	suite.userController.CreateEmployeeHandler(ctx)
	assert.Equal(suite.T(), http.StatusInternalServerError, record.Code)
}

func (suite *UserControllerTestSuite) TestFindUserByRole_Success() {
	Paginate := dto.Paginate{
		Page:       1,
		Size:       10,
		TotalRows:  0,
		TotalPages: 0,
	}

	payloadUser2 := []model.User{
		{
			Id:          "1",
			Name:        "lala",
			PhoneNumber: "08989",
			Email:       "lala@mail",
			Username:    "lalalele",
			Password:    "password",
			Point:       0,
			Role:        "admin",
		},
	}
	record := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/api/v1/users/role/"+payloadUser.Role, nil)
	assert.NoError(suite.T(), err)

	req.Header.Set("Authorization", "Bearer "+token)
	ctx, router := gin.CreateTestContext(record)
	router.GET("/api/v1/users/role/:role", suite.userController.FindUserByRoleHandler)
	ctx.Request = req

	suite.userServiceMock.On("FindUserByRole", payloadUser.Role, Paginate.Page, Paginate.Size).Return(payloadUser2, Paginate, nil)

	router.ServeHTTP(record, req)
	assert.Equal(suite.T(), http.StatusOK, record.Code)
	suite.userServiceMock.AssertExpectations(suite.T())
}
func (suite *UserControllerTestSuite) TestFindUserByRole_Failed() {
	Paginate := dto.Paginate{
		Page: 1,
		Size: 10,
	}

	record := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/api/v1/users/role/"+payloadUser.Role, nil)
	assert.NoError(suite.T(), err)

	req.Header.Set("Authorization", "Bearer "+token)
	ctx, router := gin.CreateTestContext(record)
	router.GET("/api/v1/users/role/:role", suite.userController.FindUserByRoleHandler)
	ctx.Request = req

	suite.userServiceMock.On("FindUserByRole", payloadUser.Role, Paginate.Page, Paginate.Size).Return([]model.User{}, Paginate, errors.New("not found"))

	router.ServeHTTP(record, req)
	assert.Equal(suite.T(), http.StatusInternalServerError, record.Code)
	suite.userServiceMock.AssertExpectations(suite.T())
}
func (suite *UserControllerTestSuite) TestFindUserByRole_InvalidPageSize() {
	record := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/api/v1/users/role/"+payloadUser.Role+"?page=invalid&size=invalid", nil)
	assert.NoError(suite.T(), err)

	req.Header.Set("Authorization", "Bearer "+token)
	ctx, router := gin.CreateTestContext(record)
	router.GET("/api/v1/users/role/:role", suite.userController.FindUserByRoleHandler)
	ctx.Request = req

	router.ServeHTTP(record, req)
	assert.Equal(suite.T(), http.StatusBadRequest, record.Code)
}
func (suite *UserControllerTestSuite) TestFindUserByUsername_Success() {
	record := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/api/v1/users/username/"+payload.Username, nil)
	assert.NoError(suite.T(), err)

	req.Header.Set("Authorization", "Bearer "+token)
	ctx, router := gin.CreateTestContext(record)
	router.GET("/api/v1/users/username/:username", suite.userController.FindUserByUsernameHandler)
	ctx.Request = req

	suite.userServiceMock.On("FindUserByUsername", payloadUser.Username).Return(payloadUser, nil)
	router.ServeHTTP(record, req)
	assert.Equal(suite.T(), http.StatusOK, record.Code)
}

func (suite *UserControllerTestSuite) TestFindUserByUsername_Failed() {
	record := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/api/v1/users/username/"+payload.Username, nil)
	assert.NoError(suite.T(), err)

	req.Header.Set("Authorization", "Bearer "+token)
	ctx, router := gin.CreateTestContext(record)
	router.GET("/api/v1/users/username/:username", suite.userController.FindUserByUsernameHandler)
	ctx.Request = req

	suite.userServiceMock.On("FindUserByUsername", payloadUser.Username).Return(model.User{}, errors.New("not found"))
	router.ServeHTTP(record, req)
	assert.Equal(suite.T(), http.StatusNotFound, record.Code)
}

func (suite *UserControllerTestSuite) TestFindUserByIdHandler_Success() {
	record := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/api/v1/users/id/"+payloadUser.Id, nil)
	assert.NoError(suite.T(), err)

	req.Header.Set("Authorization", "Bearer "+token)
	ctx, router := gin.CreateTestContext(record)
	router.GET("/api/v1/users/id/:id", suite.userController.FindUserByIdHandler)
	ctx.Request = req

	suite.userServiceMock.On("FindUserById", payloadUser.Id).Return(payloadUser, nil)
	router.ServeHTTP(record, req)
	assert.Equal(suite.T(), http.StatusOK, record.Code)
}

func (suite *UserControllerTestSuite) TestFindUserByIdHandler_Failed() {
	record := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/api/v1/users/id/"+payloadUser.Id, nil)
	assert.NoError(suite.T(), err)

	req.Header.Set("Authorization", "Bearer "+token)
	ctx, router := gin.CreateTestContext(record)
	router.GET("/api/v1/users/id/:id", suite.userController.FindUserByIdHandler)
	ctx.Request = req

	suite.userServiceMock.On("FindUserById", payloadUser.Id).Return(model.User{}, errors.New("not found"))
	router.ServeHTTP(record, req)
	assert.Equal(suite.T(), http.StatusNotFound, record.Code)
}

func (suite *UserControllerTestSuite) TestDeleteUserHandler_Success() {
	suite.userServiceMock.On("DeletedUser", "1").Return(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/users/1", nil)
	req.Header.Set("Content-Type", "application/json")

	ctx, _ := gin.CreateTestContext(w)
	ctx.Params = gin.Params{{Key: "id", Value: "1"}}
	ctx.Request = req

	suite.userController.DeleteUserHandler(ctx)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func (suite *UserControllerTestSuite) TestDeleteUserHandler_Fail() {
	suite.userServiceMock.On("DeletedUser", "1").Return(errors.New("user not found"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/users/1", nil)
	req.Header.Set("Content-Type", "application/json")

	ctx, _ := gin.CreateTestContext(w)
	ctx.Params = gin.Params{{Key: "id", Value: "1"}}
	ctx.Request = req

	suite.userController.DeleteUserHandler(ctx)

	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
}

//

func (suite *UserControllerTestSuite) TestUpdateUserHandler_Success() {
	suite.userServiceMock.On("UpdatedUser", "1", payloadUser).Return(payloadUser, nil)

	w := httptest.NewRecorder()
	body, _ := json.Marshal(payloadUser)
	req, _ := http.NewRequest(http.MethodPut, "/api/v1/users/1", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+token)

	ctx, _ := gin.CreateTestContext(w)
	ctx.Params = gin.Params{{Key: "id", Value: "1"}}
	ctx.Request = req
	ctx.Set("role", "admin")
	ctx.Set("userId", "1")

	suite.userController.UpdateUserHandler(ctx)
	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func (suite *UserControllerTestSuite) TestUpdateUserHandler_FailedBind() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/api/v1/users/1", bytes.NewBuffer([]byte("{invalid json}")))
	req.Header.Set("Authorization", "Bearer "+token)

	ctx, _ := gin.CreateTestContext(w)
	ctx.Params = gin.Params{{Key: "id", Value: "1"}}
	ctx.Request = req
	ctx.Set("role", "admin")
	ctx.Set("userId", "1")

	suite.userController.UpdateUserHandler(ctx)
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}

func (suite *UserControllerTestSuite) TestUpdateUserHandler_Failed() {
	suite.userServiceMock.On("UpdatedUser", "1", payloadUser).Return(model.User{}, errors.New("update error"))

	w := httptest.NewRecorder()
	body, _ := json.Marshal(payloadUser)
	req, _ := http.NewRequest(http.MethodPut, "/api/v1/users/1", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+token)

	ctx, _ := gin.CreateTestContext(w)
	ctx.Params = gin.Params{{Key: "id", Value: "1"}}
	ctx.Request = req
	ctx.Set("role", "admin")
	ctx.Set("userId", "1")

	suite.userController.UpdateUserHandler(ctx)
	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
}

func (suite *UserControllerTestSuite) TestUpdateUserHandler_Fail() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/api/v1/users/1", bytes.NewBuffer([]byte("{invalid json}")))
	req.Header.Set("Authorization", "Bearer "+token)

	ctx, _ := gin.CreateTestContext(w)
	ctx.Params = gin.Params{{Key: "id", Value: "1"}}
	ctx.Request = req

	suite.userController.UpdateUserHandler(ctx)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
}
