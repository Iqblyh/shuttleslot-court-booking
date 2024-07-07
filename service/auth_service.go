package service

import (
	"errors"
	"team2/shuttleslot/config"
	"team2/shuttleslot/model"
	"team2/shuttleslot/model/dto"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService interface {
	GenerateToken(payload model.User) (dto.LoginResponse, error)
	VerifyToken(token string) (jwt.MapClaims, error)
}

type authService struct {
	config config.SecurityConfig
}

// GenerateToken implements JwtService.
func (auth *authService) GenerateToken(payload model.User) (dto.LoginResponse, error) {
	claims := dto.JwtTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    auth.config.Issuer,
			Subject:   auth.config.Key,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(auth.config.Duration * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserId: payload.Id,
		Role:   payload.Role,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(auth.config.Key))
	if err != nil {
		return dto.LoginResponse{}, err
	}
	return dto.LoginResponse{Token: ss}, nil
}

// VerifyToken implements JwtService.
func (auth *authService) VerifyToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(auth.config.Key), nil
	})
	if err != nil {
		return nil, errors.New("failed to verify token! ")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !token.Valid || !ok || claims["iss"] != auth.config.Issuer {
		return nil, errors.New("invalid issuer or claim")
	}
	return claims, nil
}

func NewAuthService(authConfig config.SecurityConfig) AuthService {
	return &authService{
		config: authConfig,
	}
}
