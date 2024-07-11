package servicemock

import (
	"team2/shuttleslot/model"
	"team2/shuttleslot/model/dto"

	"github.com/stretchr/testify/mock"
)

type PaymentGateServiceMock struct {
	mock.Mock
}

func (m *PaymentGateServiceMock) GetPaymentURL(payment model.Payment) (string, error) {
	args := m.Called(payment)
	return args.String(0), args.Error(1)
}

func (m *PaymentGateServiceMock) PaymentProcess(payload dto.PaymentNotificationInput) (model.Payment, error) {
	args := m.Called(payload)
	return args.Get(0).(model.Payment), args.Error(1)
}
