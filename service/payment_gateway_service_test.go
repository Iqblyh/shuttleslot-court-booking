package service

import (
	"errors"
	"os"
	"team2/shuttleslot/config"
	repomock "team2/shuttleslot/mock/repo_mock"
	"team2/shuttleslot/model"
	"team2/shuttleslot/model/dto"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func mockEnvironment() {
	_ = os.Setenv("MIDTRANS_SB_SERVER_KEY", "SB-Mid-server-8O6BSnPU3-4fB_W7rG6QfdyF")
}

type PaymentServiceTestSuite struct {
	suite.Suite
	repoMock   *repomock.BookingRepositoryMock
	pS         PaymentGateService
	config     *config.PayGateConfig
	mockClient *repomock.SnapClient
}

func (suite *PaymentServiceTestSuite) SetupSuite() {
	mockEnvironment()
}

func (suite *PaymentServiceTestSuite) SetupTest() {
	suite.repoMock = new(repomock.BookingRepositoryMock)
	suite.config = &config.PayGateConfig{
		ServerKey: os.Getenv("MIDTRANS_SB_SERVER_KEY"),
	}

	suite.mockClient = new(repomock.SnapClient)

	suite.pS = NewPayGateService(*suite.config, suite.repoMock)
}

func TestPaymentServiceTestSuite(t *testing.T) {
	suite.Run(t, new(PaymentServiceTestSuite))
}

func (suite *PaymentServiceTestSuite) TestGetPaymentURL_Success() {
	payment := model.Payment{
		OrderId: "order123",
		Price:   10000,
		User: model.User{
			Name:        "Lala Test 2",
			PhoneNumber: "1234567890",
		},
		Court: model.Court{
			Name:  "Court A",
			Price: 5000,
		},
		Qty: 2,
	}

	suite.repoMock.On("FindPaymentByOrderId", payment.OrderId).Return(nil, nil)
	suite.repoMock.On("FindById", "").Return(nil, nil)

	url, err := suite.pS.GetPaymentURL(payment)

	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), url)
}

func (suite *PaymentServiceTestSuite) TestGetPaymentURL_Failed() {
	payment := model.Payment{
		OrderId: "order123",
		Price:   10000,
		User: model.User{
			Name:        "Lala Test 2",
			PhoneNumber: "1234567890",
		},
		Court: model.Court{
			Name:  "Court A",
			Price: 5000,
		},
		Qty: 1,
	}

	suite.repoMock.On("FindPaymentByOrderId", payment.OrderId).Return(nil, nil)
	suite.repoMock.On("FindById", "").Return(nil, nil)

	_, err := suite.pS.GetPaymentURL(payment)

	assert.Error(suite.T(), err)
}

func (suite *PaymentServiceTestSuite) TestPaymentProcess_Succes() {
	payload := dto.PaymentNotificationInput{
		OrderId:           "order123",
		PaymentType:       "gopay",
		TransactionStatus: "capture",
		FraudStatus:       "accept",
	}

	booking := model.Booking{
		Employee: model.User{
			Id: "employee123",
		},
	}
	suite.repoMock.On("FindPaymentByOrderId", payload.OrderId).Return(model.Payment{}, nil)
	suite.repoMock.On("FindById", "").Return(booking, nil)
	suite.repoMock.On("FindPaymentByOrderId", payload.OrderId).Return(model.Payment{}, nil)

	payment, err := suite.pS.PaymentProcess(payload)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "paid", payment.Status)
	assert.Equal(suite.T(), "credit-card", payment.PaymentMethod)
	assert.Equal(suite.T(), "employee123", payment.User.Id)
}

func (suite *PaymentServiceTestSuite) TestPaymentProcess_Cancel() {
	payload := dto.PaymentNotificationInput{
		OrderId:           "order123",
		TransactionStatus: "cancel",
	}

	suite.repoMock.On("FindPaymentByOrderId", payload.OrderId).Return(model.Payment{}, nil)
	suite.repoMock.On("FindById", "").Return(model.Booking{}, nil)

	payment, err := suite.pS.PaymentProcess(payload)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "cancel", payment.Status)
	assert.Equal(suite.T(), "", payment.PaymentMethod)
}

func (suite *PaymentServiceTestSuite) TestPaymentProcess_ErrorFindingPayment() {
	payload := dto.PaymentNotificationInput{
		OrderId: "order123",
	}

	suite.repoMock.On("FindPaymentByOrderId", payload.OrderId).Return(model.Payment{}, errors.New("error"))

	_, err := suite.pS.PaymentProcess(payload)

	assert.Error(suite.T(), err)
}
func (suite *PaymentServiceTestSuite) TestPaymentProcess_ErrorFindingPayment2() {

}

func (suite *PaymentServiceTestSuite) TestPaymentProcess_ErrorFindingId() {
	payload := dto.PaymentNotificationInput{
		OrderId:           "order123",
		PaymentType:       "credit-card",
		TransactionStatus: "capture",
		FraudStatus:       "accept",
	}

	suite.repoMock.On("FindPaymentByOrderId", payload.OrderId).Return(model.Payment{}, nil)
	suite.repoMock.On("FindById", "").Return(model.Booking{}, errors.New("error"))

	_, err := suite.pS.PaymentProcess(payload)

	assert.Error(suite.T(), err)
}
func (suite *PaymentServiceTestSuite) TestPaymentProcess_TransactionStatusSettlement() {
	payload := dto.PaymentNotificationInput{
		OrderId:           "order123",
		PaymentType:       "credit-card",
		TransactionStatus: "settlement",
		FraudStatus:       "accept",
	}

	booking := model.Booking{
		Employee: model.User{
			Id: "employee123",
		},
	}
	suite.repoMock.On("FindById", "").Return(booking, nil)
	suite.repoMock.On("FindPaymentByOrderId", payload.OrderId).Return(model.Payment{}, nil)

	payment, err := suite.pS.PaymentProcess(payload)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "paid", payment.Status)
	assert.Equal(suite.T(), "credit-card", payment.PaymentMethod)
	assert.Equal(suite.T(), "employee123", payment.User.Id)
}

func (suite *PaymentServiceTestSuite) TestPaymentProcess_TransactionStatuspending() {
	payload := dto.PaymentNotificationInput{
		OrderId:           "order123",
		PaymentType:       "credit-card",
		TransactionStatus: "pending",
		FraudStatus:       "accept",
	}

	booking := model.Booking{
		Employee: model.User{
			Id: "employee123",
		},
	}
	suite.repoMock.On("FindById", "").Return(booking, nil)
	suite.repoMock.On("FindPaymentByOrderId", payload.OrderId).Return(model.Payment{}, nil)

	payment, err := suite.pS.PaymentProcess(payload)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "paid", payment.Status)
	assert.Equal(suite.T(), "credit-card", payment.PaymentMethod)
	assert.Equal(suite.T(), "employee123", payment.User.Id)
}
