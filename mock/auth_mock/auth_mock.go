package authmock

import (
	"team2/shuttleslot/model"
	"team2/shuttleslot/model/dto"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/mock"
)

type AuthServiceMock struct {
	mock.Mock
}

func (a *AuthServiceMock) GenerateToken(payload model.User) (dto.LoginResponse, error) {
	args := a.Called(payload)
	return args.Get(0).(dto.LoginResponse), args.Error(1)
}

func (a *AuthServiceMock) VerifyToken(tokenString string) (jwt.MapClaims, error) {
	args := a.Called(tokenString)
	return args.Get(0).(jwt.MapClaims), args.Error(1)
}
