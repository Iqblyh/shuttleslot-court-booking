package service

import (
	"errors"
	repomock "team2/shuttleslot/mock/repo_mock"
	servicemock "team2/shuttleslot/mock/service_mock"
	"team2/shuttleslot/model"
	"team2/shuttleslot/model/dto"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type BookingServiceTestSuite struct {
	suite.Suite
	repoMock *repomock.BookingRepositoryMock
	bS       BookingService
	uS       *servicemock.UserServiceMock
	cS       *servicemock.CourtServiceMock
	pS       *servicemock.PaymentGateServiceMock
}

type UserServiceMock struct {
	mock.Mock
}

func (suite *BookingServiceTestSuite) SetupTest() {
	suite.repoMock = new(repomock.BookingRepositoryMock)
	suite.uS = new(servicemock.UserServiceMock)
	suite.cS = new(servicemock.CourtServiceMock)
	suite.pS = new(servicemock.PaymentGateServiceMock)
	suite.bS = NewBookingService(suite.repoMock, suite.uS, suite.cS, suite.pS)
}

func TestBookingServiceTestSuite(t *testing.T) {
	suite.Run(t, new(BookingServiceTestSuite))
}

var payload = dto.CreateBookingRequest{
	CourtId:     "court_id",
	BookingDate: "2006-01-02",
	StartTime:   "10:00",
	Hour:        2,
	CustomerId:  "customer_id",
}
var paymentNotif = dto.PaymentNotificationInput{
	TransactionStatus: "done",
	OrderId:           "Booking_1",
	PaymentType:       "gopay",
	FraudStatus:       "",
}

var createRepayRequest = dto.CreateRepayRequest{
	BookingId:     "1",
	PaymentMethod: "mid",
	EmployeeId:    "employee_id",
}
var payment = model.Payment{
	Id:            "1",
	BookingId:     "1",
	User:          user,
	Court:         court,
	OrderId:       "Booking_1",
	Description:   "paid",
	PaymentMethod: "gopay",
	Price:         30000,
	Qty:           1,
	Status:        "paid",
	PaymentURL:    "http://test-payment-url.com",
}
var booking = model.Booking{
	Id:            "1",
	Customer:      model.User{Id: "customer_id"},
	Court:         model.Court{Id: "court_id", Name: "Test Court", Price: 60000},
	Employee:      model.User{Id: "employee_id"},
	Total_Payment: 30000,
	Status:        "booked",
}

var expectedPayment = model.Payment{
	BookingId:     createRepayRequest.BookingId,
	OrderId:       "Repayment00001-12345",
	Description:   "Pelunasan Booking Test Court",
	PaymentMethod: createRepayRequest.PaymentMethod,
	Price:         15000,
	Qty:           booking.Total_Payment / court.Price,
	Court:         court,
	User:          model.User{Id: createRepayRequest.EmployeeId},
	PaymentURL:    "http://test-payment-url.com",
}
var user = model.User{
	Id:   "customer_id",
	Name: "Test Customer",
}

var court = model.Court{
	Id:    "court_id",
	Name:  "Test Court",
	Price: 60000,
}

func (suite *BookingServiceTestSuite) TestCreate_Success() {
	suite.repoMock.On("FindByDate", mock.Anything).Return([]model.Booking{}, nil)
	suite.uS.On("FindUserById", payload.CustomerId).Return(user, nil)
	suite.cS.On("FindCourtById", payload.CourtId).Return(court, nil)
	suite.repoMock.On("FindTotal", payload.CustomerId).Return(1, nil)
	suite.pS.On("GetPaymentURL", mock.Anything).Return("http://test-payment-url.com", nil)
	suite.repoMock.On("Create", mock.Anything).Return(model.Booking{
		PaymentDetails: []model.Payment{{PaymentURL: "http://test-payment-url.com"}},
	}, nil)

	_, err := suite.bS.Create(payload)
	assert.NoError(suite.T(), err, "Expected no error")
}

func (suite *BookingServiceTestSuite) TestCreate_Failure() {
	suite.repoMock.On("FindByDate", mock.Anything).Return([]model.Booking{}, errors.New("error"))
	_, err := suite.bS.Create(payload)
	assert.Error(suite.T(), err)
}

func (suite *BookingServiceTestSuite) TestCreate_Failure2() {
	suite.repoMock.On("FindByDate", mock.Anything).Return([]model.Booking{}, nil)
	suite.uS.On("FindUserById", mock.Anything).Return(model.User{}, errors.New("error"))
	_, err := suite.bS.Create(payload)
	assert.Error(suite.T(), err)
}

func (suite *BookingServiceTestSuite) TestCreate_Failure3() {
	suite.repoMock.On("FindByDate", mock.Anything).Return([]model.Booking{}, nil)
	suite.uS.On("FindUserById", payload.CustomerId).Return(user, nil)
	suite.cS.On("FindCourtById", payload.CourtId).Return(model.Court{}, errors.New("err"))
	_, err := suite.bS.Create(payload)
	assert.Error(suite.T(), err)
}

func (suite *BookingServiceTestSuite) TestCreate_Failure4() {
	suite.repoMock.On("FindByDate", mock.Anything).Return([]model.Booking{}, nil)
	suite.uS.On("FindUserById", payload.CustomerId).Return(user, nil)
	suite.cS.On("FindCourtById", payload.CourtId).Return(court, nil)
	suite.repoMock.On("FindTotal", payload.CustomerId).Return(0, errors.New("err"))
	_, err := suite.bS.Create(payload)
	assert.Error(suite.T(), err)
}

func (suite *BookingServiceTestSuite) TestCreate_Failure5() {
	suite.repoMock.On("FindByDate", mock.Anything).Return([]model.Booking{}, nil)

	suite.uS.On("FindUserById", payload.CustomerId).Return(user, nil)
	suite.cS.On("FindCourtById", payload.CourtId).Return(court, nil)
	suite.repoMock.On("FindTotal", payload.CustomerId).Return(1, nil)

	suite.pS.On("GetPaymentURL", mock.Anything).Return("http://test-payment-url.com", errors.New("err"))
	_, err := suite.bS.Create(payload)
	assert.Error(suite.T(), err)
}

func (suite *BookingServiceTestSuite) TestCreate_Failure6() {
	suite.repoMock.On("FindByDate", mock.Anything).Return([]model.Booking{}, nil)

	suite.uS.On("FindUserById", payload.CustomerId).Return(user, nil)
	suite.cS.On("FindCourtById", payload.CourtId).Return(court, nil)
	suite.repoMock.On("FindTotal", payload.CustomerId).Return(1, nil)

	suite.pS.On("GetPaymentURL", mock.Anything).Return("http://test-payment-url.com", nil)
	suite.repoMock.On("Create", mock.Anything).Return(model.Booking{
		PaymentDetails: []model.Payment{{PaymentURL: "http://test-payment-url.com"}},
	}, errors.New("err"))

	_, err := suite.bS.Create(payload)
	assert.Error(suite.T(), err)
}

func (suite *BookingServiceTestSuite) TestUpdatePayment_Success_Booking() {
	suite.pS.On("PaymentProcess", paymentNotif).Return(payment, nil)
	suite.repoMock.On("UpdateStatus", payment).Return(nil)

	err := suite.bS.UpdatePayment(paymentNotif)
	suite.NoError(err)
	suite.pS.AssertExpectations(suite.T())
	suite.repoMock.AssertExpectations(suite.T())
}
func (suite *BookingServiceTestSuite) TestUpdatePayment_Success_Repayment() {
	paymentNotif := dto.PaymentNotificationInput{
		TransactionStatus: "done",
		OrderId:           "Repayment_1",
		PaymentType:       "gopay",
		FraudStatus:       "",
	}

	payment := model.Payment{
		Id:            "1",
		BookingId:     "1",
		User:          user,
		Court:         court,
		OrderId:       "Repayment_1",
		Description:   "paid",
		PaymentMethod: "gopay",
		Price:         30000,
		Qty:           1,
		Status:        "paid",
		PaymentURL:    "http://test-payment-url.com",
	}

	suite.pS.On("PaymentProcess", paymentNotif).Return(payment, nil)
	suite.repoMock.On("UpdateRepaymentStatus", payment).Return(nil)

	err := suite.bS.UpdatePayment(paymentNotif)
	suite.NoError(err)
	suite.pS.AssertExpectations(suite.T())
	suite.repoMock.AssertExpectations(suite.T())
}
func (suite *BookingServiceTestSuite) TestUpdatePayment_Failed_Booking() {
	suite.pS.On("PaymentProcess", paymentNotif).Return(model.Payment{}, errors.New("error"))
	suite.repoMock.On("UpdateStatus", payment).Return(errors.New("error"))

	err := suite.bS.UpdatePayment(paymentNotif)
	suite.Error(err)
}
func (suite *BookingServiceTestSuite) TestUpdatePayment_Failed_Booking2() {
	suite.pS.On("PaymentProcess", paymentNotif).Return(payment, nil)
	suite.repoMock.On("UpdateStatus", payment).Return(errors.New("error"))

	err := suite.bS.UpdatePayment(paymentNotif)
	suite.Error(err)
}

func (suite *BookingServiceTestSuite) TestUpdatePayment_Failed_Repayment() {
	paymentNotif := dto.PaymentNotificationInput{
		TransactionStatus: "done",
		OrderId:           "Repayment_1",
		PaymentType:       "gopay",
		FraudStatus:       "",
	}

	payment := model.Payment{
		Id:            "1",
		BookingId:     "1",
		User:          user,
		Court:         court,
		OrderId:       "Repayment_1",
		Description:   "paid",
		PaymentMethod: "gopay",
		Price:         30000,
		Qty:           1,
		Status:        "paid",
		PaymentURL:    "http://test-payment-url.com",
	}

	suite.pS.On("PaymentProcess", paymentNotif).Return(payment, nil)
	suite.repoMock.On("UpdateRepaymentStatus", payment).Return(errors.New("error"))

	err := suite.bS.UpdatePayment(paymentNotif)
	suite.Error(err)
}

func (suite *BookingServiceTestSuite) TestCreateRepay_PaymentMethodNotMid() {
	suite.repoMock.On("FindById", createRepayRequest.BookingId).Return(booking, nil)
	suite.uS.On("FindUserById", booking.Customer.Id).Return(user, nil)
	suite.cS.On("FindCourtById", booking.Court.Id).Return(court, nil)
	suite.repoMock.On("FindTotal", booking.Customer.Id).Return(1, nil)
	suite.repoMock.On("CreateRepay", mock.AnythingOfType("model.Payment")).Return(expectedPayment, nil)

	createRepayRequest.PaymentMethod = "cash"

	payment, err := suite.bS.CreateRepay(createRepayRequest)

	suite.NoError(err)
	assert.Equal(suite.T(), expectedPayment, payment)
	suite.repoMock.AssertExpectations(suite.T())
	suite.uS.AssertExpectations(suite.T())
	suite.cS.AssertExpectations(suite.T())
}
func (suite *BookingServiceTestSuite) TestCreateRepay_PaymentMethodMid() {
	suite.repoMock.On("FindById", createRepayRequest.BookingId).Return(booking, nil)
	suite.uS.On("FindUserById", booking.Customer.Id).Return(user, nil)
	suite.cS.On("FindCourtById", booking.Court.Id).Return(court, nil)
	suite.repoMock.On("FindTotal", booking.Customer.Id).Return(1, nil)
	suite.pS.On("GetPaymentURL", mock.AnythingOfType("model.Payment")).Return("http://test-payment-url.com", nil)
	suite.repoMock.On("CreateRepay", mock.AnythingOfType("model.Payment")).Return(expectedPayment, nil)

	createRepayRequest.PaymentMethod = "mid"

	payment, err := suite.bS.CreateRepay(createRepayRequest)

	suite.NoError(err)
	assert.Equal(suite.T(), expectedPayment, payment)
	suite.repoMock.AssertExpectations(suite.T())
	suite.uS.AssertExpectations(suite.T())
	suite.cS.AssertExpectations(suite.T())
	suite.pS.AssertExpectations(suite.T())
}
func (suite *BookingServiceTestSuite) TestCreateRepay_Failed() {
	booking := model.Booking{
		Id:            "1",
		Customer:      model.User{Id: "customer_id"},
		Court:         model.Court{Id: "court_id", Name: "Test Court", Price: 60000},
		Employee:      model.User{Id: "employee_id"},
		Total_Payment: 30000,
		Status:        "not booked",
	}
	suite.repoMock.On("FindById", createRepayRequest.BookingId).Return(booking, nil)

	_, err := suite.bS.CreateRepay(createRepayRequest)
	suite.EqualError(err, "this booking still not booked")
	suite.repoMock.AssertExpectations(suite.T())
}

func (suite *BookingServiceTestSuite) TestCreateRepay_Failed1() {
	suite.repoMock.On("FindById", createRepayRequest.BookingId).Return(model.Booking{}, errors.New("error"))
	_, err := suite.bS.CreateRepay(createRepayRequest)
	suite.Error(err)
}
func (suite *BookingServiceTestSuite) TestCreateRepay_Failed2() {
	suite.repoMock.On("FindById", createRepayRequest.BookingId).Return(booking, nil)
	suite.uS.On("FindUserById", booking.Customer.Id).Return(model.User{}, errors.New("error"))
	_, err := suite.bS.CreateRepay(createRepayRequest)
	suite.Error(err)
}

func (suite *BookingServiceTestSuite) TestCreateRepay_Failed3() {
	suite.repoMock.On("FindById", createRepayRequest.BookingId).Return(booking, nil)
	suite.uS.On("FindUserById", booking.Customer.Id).Return(user, nil)
	suite.cS.On("FindCourtById", booking.Court.Id).Return(model.Court{}, errors.New("error"))
	_, err := suite.bS.CreateRepay(createRepayRequest)
	suite.Error(err)
}

func (suite *BookingServiceTestSuite) TestCreateRepay_Failed4() {
	suite.repoMock.On("FindById", createRepayRequest.BookingId).Return(booking, nil)
	suite.uS.On("FindUserById", booking.Customer.Id).Return(user, nil)
	suite.cS.On("FindCourtById", booking.Court.Id).Return(court, nil)
	suite.repoMock.On("FindTotal", booking.Customer.Id).Return(0, errors.New("error"))
	_, err := suite.bS.CreateRepay(createRepayRequest)
	suite.Error(err)
}

func (suite *BookingServiceTestSuite) TestCreateRepay_Failed5() {
	suite.repoMock.On("FindById", createRepayRequest.BookingId).Return(booking, nil)
	suite.uS.On("FindUserById", booking.Customer.Id).Return(user, nil)
	suite.cS.On("FindCourtById", booking.Court.Id).Return(court, nil)
	suite.repoMock.On("FindTotal", booking.Customer.Id).Return(1, nil)
	suite.pS.On("GetPaymentURL", mock.AnythingOfType("model.Payment")).Return("", errors.New("error"))
	_, err := suite.bS.CreateRepay(createRepayRequest)
	suite.Error(err)
}

func (suite *BookingServiceTestSuite) TestCreateRepay_Failed6() {
	suite.repoMock.On("FindById", createRepayRequest.BookingId).Return(booking, nil)
	suite.uS.On("FindUserById", booking.Customer.Id).Return(user, nil)
	suite.cS.On("FindCourtById", booking.Court.Id).Return(court, nil)
	suite.repoMock.On("FindTotal", booking.Customer.Id).Return(1, nil)
	suite.pS.On("GetPaymentURL", mock.AnythingOfType("model.Payment")).Return("http://test-payment-url.com", nil)
	suite.repoMock.On("CreateRepay", mock.AnythingOfType("model.Payment")).Return(model.Payment{}, errors.New("error"))
	_, err := suite.bS.CreateRepay(createRepayRequest)
	suite.Error(err)
}

func (suite *BookingServiceTestSuite) TestCreateRepay_Failed_ScheduleTooLong() {
	suite.repoMock.On("FindById", createRepayRequest.BookingId).Return(booking, nil)
	suite.uS.On("FindUserById", booking.Customer.Id).Return(user, nil)
	suite.cS.On("FindCourtById", booking.Court.Id).Return(court, nil)
	suite.repoMock.On("FindTotal", booking.Customer.Id).Return(1, nil)

	suite.repoMock.On("CreateRepay", mock.AnythingOfType("model.Payment")).Return(model.Payment{}, errors.New("schedule too long"))

	createRepayRequest.PaymentMethod = "cash"
	_, err := suite.bS.CreateRepay(createRepayRequest)

	suite.Error(err)
	assert.Contains(suite.T(), err.Error(), "schedule too long")
	suite.repoMock.AssertExpectations(suite.T())
	suite.uS.AssertExpectations(suite.T())
	suite.cS.AssertExpectations(suite.T())
	suite.pS.AssertExpectations(suite.T())
}

func (suite *BookingServiceTestSuite) TestFindAllBookings_Success() {
	bookings := []model.Booking{
		{Id: "1", Customer: model.User{Id: "customer_id_1"}, Court: model.Court{Id: "court_id_1"}, Total_Payment: 100, Status: "booked"},
		{Id: "2", Customer: model.User{Id: "customer_id_2"}, Court: model.Court{Id: "court_id_2"}, Total_Payment: 200, Status: "booked"},
	}
	paginate := dto.Paginate{Page: 1, Size: 2, TotalRows: 2}

	suite.repoMock.On("FindAll", 1, 2).Return(bookings, paginate, nil)

	result, pag, err := suite.bS.FindAllBookings(1, 2)

	suite.NoError(err)
	suite.Equal(bookings, result)
	suite.Equal(paginate, pag)
	suite.repoMock.AssertExpectations(suite.T())
}

func (suite *BookingServiceTestSuite) TestFindAllBookings_Failed() {
	suite.repoMock.On("FindAll", 1, 2).Return([]model.Booking{}, dto.Paginate{}, errors.New("error"))

	_, _, err := suite.bS.FindAllBookings(1, 2)

	suite.Error(err)
	suite.repoMock.AssertExpectations(suite.T())
}

func (suite *BookingServiceTestSuite) TestFindBookedCourt_Success() {
	page := 1
	size := 10
	bookingDate := time.Time{}
	expectedBookings := []model.Booking{booking}
	expectedPaginate := dto.Paginate{Page: page, Size: size, TotalRows: 1, TotalPages: 1}

	suite.repoMock.On("FindBooked", bookingDate, page, size).Return(expectedBookings, expectedPaginate, nil)
	suite.cS.On("FindCourtById", booking.Court.Id).Return(court, nil)

	bookings, paginate, err := suite.bS.FindBookedCourt(bookingDate, page, size)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedBookings, bookings)
	assert.Equal(suite.T(), expectedPaginate, paginate)

	suite.repoMock.AssertExpectations(suite.T())
	suite.cS.AssertExpectations(suite.T())
}
func (suite *BookingServiceTestSuite) TestFindBookedCourt_Failed() {
	page := 1
	size := 10
	bookingDate := time.Time{}

	expectedPaginate := dto.Paginate{Page: 0, Size: 0, TotalRows: 0, TotalPages: 0}

	suite.repoMock.On("FindBooked", bookingDate, page, size).Return([]model.Booking{}, expectedPaginate, errors.New("error"))

	_, paginate, err := suite.bS.FindBookedCourt(bookingDate, page, size)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), expectedPaginate, paginate)

	suite.repoMock.AssertExpectations(suite.T())
	suite.cS.AssertExpectations(suite.T())
}

func (suite *BookingServiceTestSuite) TestFindBookedCourt_Failed2() {
	page := 1
	size := 10
	bookingDate := time.Time{}
	expectedBookings := []model.Booking{booking}
	expectedPaginate := dto.Paginate{Page: 0, Size: 0, TotalRows: 0, TotalPages: 0}

	suite.repoMock.On("FindBooked", bookingDate, page, size).Return(expectedBookings, expectedPaginate, nil)
	suite.cS.On("FindCourtById", booking.Court.Id).Return(model.Court{}, errors.New("error"))

	_, paginate, err := suite.bS.FindBookedCourt(bookingDate, page, size)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), expectedPaginate, paginate)

	suite.repoMock.AssertExpectations(suite.T())
	suite.cS.AssertExpectations(suite.T())
}

func (suite *BookingServiceTestSuite) TestFindEndingBookings_Success() {
	page := 1
	size := 10
	bookingDate := time.Time{}
	expectedBookings := []model.Booking{booking}
	expectedPaginate := dto.Paginate{Page: page, Size: size, TotalRows: 1, TotalPages: 1}

	suite.repoMock.On("FindEnding", bookingDate, page, size).Return(expectedBookings, expectedPaginate, nil)
	suite.uS.On("FindUserById", booking.Customer.Id).Return(user, nil)
	suite.cS.On("FindCourtById", booking.Court.Id).Return(court, nil)

	bookings, paginate, err := suite.bS.FindEndingBookings(bookingDate, page, size)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedBookings, bookings)
	assert.Equal(suite.T(), expectedPaginate, paginate)

	suite.repoMock.AssertExpectations(suite.T())
	suite.cS.AssertExpectations(suite.T())
}

func (suite *BookingServiceTestSuite) TestFindEndingBookings_Failed() {
	page := 1
	size := 10
	bookingDate := time.Time{}
	expectedBookings := []model.Booking{booking}
	expectedPaginate := dto.Paginate{Page: 0, Size: 0, TotalRows: 0, TotalPages: 0}

	suite.repoMock.On("FindEnding", bookingDate, page, size).Return(expectedBookings, expectedPaginate, nil)
	suite.uS.On("FindUserById", booking.Customer.Id).Return(user, nil)
	suite.cS.On("FindCourtById", booking.Court.Id).Return(model.Court{}, errors.New("error"))

	_, paginate, err := suite.bS.FindEndingBookings(bookingDate, page, size)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), expectedPaginate, paginate)

	suite.repoMock.AssertExpectations(suite.T())
	suite.cS.AssertExpectations(suite.T())
}

func (suite *BookingServiceTestSuite) TestFindEndingBookings_Failed2() {
	page := 1
	size := 10
	bookingDate := time.Time{}
	expectedBookings := []model.Booking{booking}
	expectedPaginate := dto.Paginate{Page: 0, Size: 0, TotalRows: 0, TotalPages: 0}

	suite.repoMock.On("FindEnding", bookingDate, page, size).Return(expectedBookings, expectedPaginate, nil)
	suite.uS.On("FindUserById", booking.Customer.Id).Return(model.User{}, errors.New("error"))
	suite.cS.On("FindCourtById", booking.Court.Id).Return(court, nil)

	_, paginate, err := suite.bS.FindEndingBookings(bookingDate, page, size)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), expectedPaginate, paginate)

	suite.repoMock.AssertExpectations(suite.T())
}

func (suite *BookingServiceTestSuite) TestFindEndingBookings_Failed3() {
	page := 1
	size := 10
	bookingDate := time.Time{}
	expectedPaginate := dto.Paginate{Page: 0, Size: 0, TotalRows: 0, TotalPages: 0}

	suite.repoMock.On("FindEnding", bookingDate, page, size).Return([]model.Booking{}, expectedPaginate, errors.New("error"))

	_, paginate, err := suite.bS.FindEndingBookings(bookingDate, page, size)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), expectedPaginate, paginate)

	suite.repoMock.AssertExpectations(suite.T())
	suite.cS.AssertExpectations(suite.T())
}

func (suite *BookingServiceTestSuite) TestFindPaymentReport_Success() {
	payments := []model.Payment{}
	paginate := dto.Paginate{Page: 1, Size: 2, TotalRows: 2}
	totalCount := int64(2)

	suite.repoMock.On("FindPaymentReport", 1, 2, 3, 4, 5, "6").Return(payments, paginate, totalCount, nil)

	resultPayments, resultPaginate, resultTotalCount, err := suite.bS.FindPaymentReport(1, 2, 3, 4, 5, "6")

	suite.NoError(err)
	suite.Equal(payments, resultPayments)
	suite.Equal(paginate, resultPaginate)
	suite.Equal(totalCount, resultTotalCount)
	suite.repoMock.AssertExpectations(suite.T())
}

func (suite *BookingServiceTestSuite) TestFindPaymentReport_Failed() {
	totalCount := int64(2)
	suite.repoMock.On("FindPaymentReport", 1, 2, 3, 4, 5, "6").Return([]model.Payment{}, dto.Paginate{}, totalCount, errors.New("error"))

	_, _, _, err := suite.bS.FindPaymentReport(1, 2, 3, 4, 5, "6")

	suite.Error(err)
	suite.repoMock.AssertExpectations(suite.T())
}
