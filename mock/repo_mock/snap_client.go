package repomock

import (
	"github.com/midtrans/midtrans-go/snap"
	"github.com/stretchr/testify/mock"
)

type SnapClient struct {
	mock.Mock
}

func (m *SnapClient) CreateTransaction(req *snap.Request) (*snap.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*snap.Response), args.Error(1)
}
