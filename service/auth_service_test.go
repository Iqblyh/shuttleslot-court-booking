package service

import (
	"team2/shuttleslot/config"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AuthServiceTestSuite struct {
	suite.Suite
	authConfig config.SecurityConfig
	aS         AuthService
}

func (suite *AuthServiceTestSuite) SetupTest() {
	suite.authConfig = config.SecurityConfig{
		Issuer:   "testIssuer",
		Key:      "testKey",
		Duration: 1,
	}
	suite.aS = NewAuthService(suite.authConfig)
}

func TestAuthServiceTestSuite(t *testing.T) {
	suite.Run(t, new(AuthServiceTestSuite))
}

func (suite *AuthServiceTestSuite) TestGenerateToken_Success() {
	loginResponse, err := suite.aS.GenerateToken(mockUser)

	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), loginResponse.Token)
}

func (suite *AuthServiceTestSuite) TestVerifyToken_Success() {
	token, err := suite.aS.GenerateToken(mockUser)
	assert.NoError(suite.T(), err)

	claims, err := suite.aS.VerifyToken(token.Token)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), nil, claims["UserId"])
	assert.Equal(suite.T(), nil, claims["Role"])
	assert.Equal(suite.T(), suite.authConfig.Issuer, claims["iss"])
}

func (suite *AuthServiceTestSuite) TestVerifyToken_Fail_InvalidToken() {
	_, err := suite.aS.VerifyToken("invalidToken")
	assert.Error(suite.T(), err)
	assert.EqualError(suite.T(), err, "failed to verify token! ")
}

func (suite *AuthServiceTestSuite) TestVerifyToken_Fail_InvalidIssuer() {
	token, err := suite.aS.GenerateToken(mockUser)
	assert.NoError(suite.T(), err)

	suite.aS = NewAuthService(config.SecurityConfig{Issuer: "invalidIssuer", Key: suite.authConfig.Key})
	_, err = suite.aS.VerifyToken(token.Token)
	assert.Error(suite.T(), err)
	assert.EqualError(suite.T(), err, "invalid issuer or claim")
}
