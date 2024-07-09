package utilmock

import (
	"github.com/stretchr/testify/mock"
)

type MockUtil struct {
	mock.Mock
}

func (m *MockUtil) EncryptPassword(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *MockUtil) ComparePasswordHash(passwordHash string, passwordDb string) error {
	args := m.Called(passwordHash, passwordDb)
	return args.Error(0)
}
