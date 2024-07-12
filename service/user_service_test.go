package service

import (
	"errors"
	authmock "team2/shuttleslot/mock/auth_mock"
	repomock "team2/shuttleslot/mock/repo_mock"
	utilmock "team2/shuttleslot/mock/util_mock"
	"team2/shuttleslot/model"
	"team2/shuttleslot/model/dto"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

var mockUser = model.User{
	Id:          "1",
	Name:        "lala",
	PhoneNumber: "123",
	Email:       "@mail",
	Username:    "lalalele",
	Password:    "$2a$10$8TshMQnYvNo..jD6mePld.p5C2mCHfCF8mMTc7.phtuUhYen0Ny3G",
	Point:       0,
	Role:        "customer",
	CreatedAt:   time.Now(),
	UpdatedAt:   time.Now(),
}
var loginPayload = dto.LoginRequest{
	Username: "lalalele",
	Password: "password",
}

type UserServiceTestSuite struct {
	suite.Suite
	repoUserMock *repomock.UserRepositoryMock
	uS           UserService
	aU           *authmock.AuthServiceMock
	uM           *utilmock.MockUtil
}

func (suite *UserServiceTestSuite) SetupTest() {
	suite.repoUserMock = new(repomock.UserRepositoryMock)
	suite.uM = new(utilmock.MockUtil)
	suite.aU = new(authmock.AuthServiceMock)
	suite.uS = NewUserService(suite.repoUserMock, suite.aU, suite.uM)
}

func TestUserServiceTestSuite(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}

func (suite *UserServiceTestSuite) TestLogin_Success() {
	expectedResponse := dto.LoginResponse{
		Token: "mocked-jwt-token",
	}

	suite.repoUserMock.On("FindUserByUsername", loginPayload.Username).Return(mockUser, nil)
	suite.uM.On("ComparePasswordHash", mockUser.Password, loginPayload.Password).Return(nil)
	mockUser.Password = ""

	suite.aU.On("GenerateToken", mockUser).Return(expectedResponse, nil)
	result, err := suite.uS.Login(loginPayload)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedResponse, result)
}

func (suite *UserServiceTestSuite) TestLogin_Failed() {
	suite.repoUserMock.On("FindUserByUsername", loginPayload.Username).Return(model.User{}, errors.New("error"))
	_, err := suite.uS.Login(loginPayload)
	assert.Error(suite.T(), err)
}

func (suite *UserServiceTestSuite) TestLogin_Failed2() {
	suite.repoUserMock.On("FindUserByUsername", loginPayload.Username).Return(mockUser, nil)
	suite.uM.On("ComparePasswordHash", mockUser.Password, loginPayload.Password).Return(errors.New("error hash"))

	_, err := suite.uS.Login(loginPayload)
	assert.Error(suite.T(), err)
}

func (suite *UserServiceTestSuite) TestLogin_Failed3() {
	suite.repoUserMock.On("FindUserByUsername", loginPayload.Username).Return(mockUser, nil)

	suite.uM.On("ComparePasswordHash", mockUser.Password, loginPayload.Password).Return(nil)
	mockUser.Password = ""
	suite.aU.On("GenerateToken", mockUser).Return(dto.LoginResponse{}, errors.New("failed to create token"))

	_, err := suite.uS.Login(loginPayload)
	assert.Error(suite.T(), err)
}

func (suite *UserServiceTestSuite) TestCreateAdmin_Success() {
	role := "admin"
	suite.uM.On("EncryptPassword", mock.AnythingOfType("string")).Return(mockUser.Password, nil)
	suite.repoUserMock.On("CreateAdmin", mock.Anything).Return(mockUser, nil)

	createdUser, err := suite.uS.CreateAdmin(mockUser)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), mockUser.Password, createdUser.Password)
	assert.Equal(suite.T(), "admin", role)
}

func (suite *UserServiceTestSuite) TestCreateAdmin_Fail() {
	suite.uM.On("EncryptPassword", mock.AnythingOfType("string")).Return("", errors.New("err"))
	suite.repoUserMock.On("CreateAdmin", mock.Anything).Return(mockUser, nil)

	_, err := suite.uS.CreateAdmin(model.User{})
	assert.Error(suite.T(), err)
}

func (suite *UserServiceTestSuite) TestCreateCustomer_Success() {
	role := "customer"
	suite.uM.On("EncryptPassword", mock.AnythingOfType("string")).Return(mockUser.Password, nil)
	suite.repoUserMock.On("CreateCustomer", mock.Anything).Return(mockUser, nil)

	createdUser, err := suite.uS.CreateCustomer(mockUser)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), mockUser.Password, createdUser.Password)
	assert.Equal(suite.T(), "customer", role)
}

func (suite *UserServiceTestSuite) TestCreateCustomer_Fail() {
	suite.uM.On("EncryptPassword", mock.AnythingOfType("string")).Return("", errors.New("err"))
	suite.repoUserMock.On("CreateCustomer", mock.Anything).Return(mockUser, nil)

	_, err := suite.uS.CreateCustomer(mockUser)

	assert.Error(suite.T(), err)
}

func (suite *UserServiceTestSuite) TestCreateEmployee_Success() {
	role := "employee"
	suite.uM.On("EncryptPassword", mock.AnythingOfType("string")).Return(mockUser.Password, nil)
	suite.repoUserMock.On("CreateEmployee", mock.Anything).Return(mockUser, nil)

	createdUser, err := suite.uS.CreateEmployee(mockUser)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), mockUser.Password, createdUser.Password)
	assert.Equal(suite.T(), "employee", role)
}

func (suite *UserServiceTestSuite) TestCreateEmployee_Fail() {
	suite.uM.On("EncryptPassword", mock.AnythingOfType("string")).Return("", errors.New("err"))
	suite.repoUserMock.On("CreateEmployee", mock.Anything).Return(mockUser, nil)

	_, err := suite.uS.CreateEmployee(mockUser)

	assert.Error(suite.T(), err)
}

func (suite *UserServiceTestSuite) TestFindUserByRole_Success() {
	page := 1
	size := 10
	role := "customer"
	mockUsers := []model.User{mockUser}
	mockPaginate := dto.Paginate{
		Page:       page,
		Size:       size,
		TotalRows:  1,
		TotalPages: 1,
	}

	suite.repoUserMock.On("FindUserByRole", role, page, size).Return(mockUsers, mockPaginate, nil)

	users, paginate, err := suite.uS.FindUserByRole(role, page, size)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), mockUsers, users)
	assert.Equal(suite.T(), mockPaginate, paginate)
}

func (suite *UserServiceTestSuite) TestFindUserByRole_Fail() {
	page := 1
	size := 10
	role := "customer"
	mockError := errors.New("error finding users by role")

	suite.repoUserMock.On("FindUserByRole", role, page, size).Return([]model.User{}, dto.Paginate{}, mockError)

	_, paginate, err := suite.uS.FindUserByRole(role, page, size)

	assert.Error(suite.T(), err)
	assert.EqualError(suite.T(), err, "error finding users by role")
	assert.Equal(suite.T(), dto.Paginate{}, paginate)
}

func (suite *UserServiceTestSuite) TestFindUserByUsername_Success() {
	suite.repoUserMock.On("FindUserByUsername", mockUser.Username).Return(mockUser, nil)

	user, err := suite.uS.FindUserByUsername(mockUser.Username)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), mockUser, user)
}

func (suite *UserServiceTestSuite) TestFindUserByUsername_Fail() {
	suite.repoUserMock.On("FindUserByUsername", mockUser.Username).Return(model.User{}, errors.New("error"))

	user, err := suite.uS.FindUserByUsername(mockUser.Username)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), model.User{}, user)
}

func (suite *UserServiceTestSuite) TestFindUserById_Success() {
	suite.repoUserMock.On("FindUserById", mockUser.Id).Return(mockUser, nil)

	user, err := suite.uS.FindUserById(mockUser.Id)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), mockUser, user)
}

func (suite *UserServiceTestSuite) TestFindUserById_Fail() {
	suite.repoUserMock.On("FindUserById", mockUser.Id).Return(model.User{}, errors.New("error"))

	user, err := suite.uS.FindUserById(mockUser.Id)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), model.User{}, user)
}

func (suite *UserServiceTestSuite) TestUpdateUser_Success() {
	suite.repoUserMock.On("FindUserById", "user_id").Return(mockUser, nil)
	suite.uM.On("EncryptPassword", mockUser.Password).Return(mockUser.Password, nil)
	suite.repoUserMock.On("UpdateUser", "user_id", mockUser).Return(mockUser, nil)

	returnedUser, err := suite.uS.UpdatedUser("user_id", mockUser)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), returnedUser)
	suite.repoUserMock.AssertExpectations(suite.T())
}

func (suite *UserServiceTestSuite) TestUpdateUser_Fail() {
	suite.repoUserMock.On("FindUserById", mockUser.Id).Return(model.User{}, errors.New("user not found"))

	updatedUser, err := suite.uS.UpdatedUser(mockUser.Id, mockUser)

	assert.Error(suite.T(), err)
	assert.EqualError(suite.T(), err, "user not found")
	assert.Equal(suite.T(), model.User{}, updatedUser)
}
func (suite *UserServiceTestSuite) TestUpdatedUser_Fail_EncryptPassword() {
	id := "1"
	payload := model.User{
		Password: "newpassword",
	}

	suite.repoUserMock.On("FindUserById", id).Return(mockUser, nil)
	suite.uM.On("EncryptPassword", payload.Password).Return("", errors.New("error in encrypting password"))

	result, err := suite.uS.UpdatedUser(id, payload)

	assert.Error(suite.T(), err)
	assert.EqualError(suite.T(), err, "error in encrypting password")
	assert.Equal(suite.T(), model.User{}, result)
}

func (suite *UserServiceTestSuite) TestDeleteUser_Success() {
	suite.repoUserMock.On("FindUserById", mockUser.Id).Return(mockUser, nil)
	suite.repoUserMock.On("DeleteUser", mockUser.Id).Return(nil)

	err := suite.uS.DeletedUser(mockUser.Id)

	assert.NoError(suite.T(), err)
}

func (suite *UserServiceTestSuite) TestDeleteUser_Fail() {
	suite.repoUserMock.On("FindUserById", mockUser.Id).Return(model.User{}, errors.New("user not found"))

	err := suite.uS.DeletedUser(mockUser.Id)

	assert.Error(suite.T(), err)
	assert.EqualError(suite.T(), err, "user not found")
}
