package util

import "golang.org/x/crypto/bcrypt"

func EncryptPassword(password string) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}
	return string(passwordHash), nil
}

func ComparePasswordHash(passwordHash string, passwordDb string) error {
	return bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(passwordDb))
}
