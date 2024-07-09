package util

import "golang.org/x/crypto/bcrypt"

type UtilInterface interface {
	EncryptPassword(password string) (string, error)
	ComparePasswordHash(passwordHash string, passwordDb string) error
}

type UtilService struct{}

func (u *UtilService) EncryptPassword(password string) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(passwordHash), nil
}

func (u *UtilService) ComparePasswordHash(passwordHash string, passwordDb string) error {
	return bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(passwordDb))
}

func NewUtilService() UtilInterface {
	return &UtilService{}
}
