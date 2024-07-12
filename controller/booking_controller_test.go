package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	mock "team2/shuttleslot/mock"
	"team2/shuttleslot/util"

	servicemock "team2/shuttleslot/mock/service_mock"
	"team2/shuttleslot/model"
	"team2/shuttleslot/model/dto"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var (
	mockBooking = model.Booking{
		Id:            "1",
		Customer:      model.User{},
		Court:         model.Court{},
		Employee:      model.User{},
		BookingDate:   time.Time{},
		StartTime:     time.Time{},
		EndTime:       time.Time{},
		Total_Payment: 60000,
		Status:        "Status",
		PaymentDetails: []model.Payment{
			{
				Id:            "1",
				BookingId:     "1",
				User:          model.User{},
				Court:         model.Court{},
				OrderId:       "1",
				Description:   "Description",
				PaymentMethod: "gopay",
				Price:         60000,
				Qty:           1,
				Status:        "paid",
				PaymentURL:    "payment.com",
			},
		},
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}
)

type BookingControllerTestSuite struct {
	suite.Suite
	bookingServiceMock *servicemock.BookingServiceMock
	rg                 *gin.RouterGroup
	controller         *BookingController
	middlewareMock     *mock.AuthMiddlewareMock
}

func (suite *BookingControllerTestSuite) SetupTest() {
	suite.bookingServiceMock = new(servicemock.BookingServiceMock)
	suite.rg = gin.Default().Group("/api/v1")
	suite.controller = NewBookingController(suite.bookingServiceMock, suite.middlewareMock, suite.rg.Group("/bookings"))
	suite.controller.Route()
}

func TestBookingControllerTestSuite(t *testing.T) {
	suite.Run(t, new(BookingControllerTestSuite))
}

func (suite *BookingControllerTestSuite) TestCreateBookingHandler_Success() {
	now := time.Now()
	futureDate := now.AddDate(0, 0, 1).Format("02-01-2006")
	futureTime := now.Add(time.Hour).Format("15:04:05")

	payload := dto.CreateBookingRequest{
		CourtId:     "1",
		BookingDate: futureDate,
		StartTime:   futureTime,
		Hour:        1,
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/bookings", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	ctx, router := gin.CreateTestContext(rec)
	router.POST("/api/v1/bookings", suite.controller.CreateBookingHandler)
	ctx.Request = req

	ctx.Set("userId", "1")

	expectedBooking := mockBooking
	suite.bookingServiceMock.On("Create", payload).Return(expectedBooking, nil)

	router.ServeHTTP(rec, req)

	var responseBody map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &responseBody)
	if err != nil {
		suite.T().Logf("Failed to unmarshal response body: %v", err)
	} else {
		suite.T().Logf("Response body: %v", responseBody)
	}

	assert.Equal(suite.T(), http.StatusOK, rec.Code)
}

func (suite *BookingControllerTestSuite) TestCreateBookingHandler_InvalidJSON() {
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/bookings", strings.NewReader("{invalid json}"))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	ctx, router := gin.CreateTestContext(rec)
	router.POST("/api/v1/bookings", suite.controller.CreateBookingHandler)
	ctx.Request = req

	router.ServeHTTP(rec, req)
	assert.Equal(suite.T(), http.StatusBadRequest, rec.Code)
}

func (suite *BookingControllerTestSuite) TestCreateBookingHandler_InvalidDate() {
	payload := dto.CreateBookingRequest{
		CourtId:     "1",
		BookingDate: "invalid-date",
		StartTime:   "14:00:00",
		Hour:        1,
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/bookings", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	ctx, router := gin.CreateTestContext(rec)
	router.POST("/api/v1/bookings", suite.controller.CreateBookingHandler)
	ctx.Request = req

	router.ServeHTTP(rec, req)
	assert.Equal(suite.T(), http.StatusBadRequest, rec.Code)
	var response dto.SingleResponse
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "invalid date or time format, use 'dd-mm-yyyy for bookingDate and 'hh-mm-ss' for startTime", response.Status.Message)
}

func (suite *BookingControllerTestSuite) TestCreateBookingHandler_PastBookingDate() {
	payload := dto.CreateBookingRequest{
		CourtId:     "1",
		BookingDate: "10-07-2020",
		StartTime:   "14:00:00",
		Hour:        1,
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/bookings", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	ctx, router := gin.CreateTestContext(rec)
	router.POST("/api/v1/bookings", suite.controller.CreateBookingHandler)
	ctx.Request = req

	router.ServeHTTP(rec, req)
	assert.Equal(suite.T(), http.StatusBadRequest, rec.Code)
	var response dto.SingleResponse
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "booking date cant in the past", response.Status.Message)
}

func (suite *BookingControllerTestSuite) TestCreateBookingHandler_PastStartTime() {
	now := time.Now()
	payload := dto.CreateBookingRequest{
		CourtId:     "1",
		BookingDate: now.Format("02-01-2006"),
		StartTime:   now.Add(-time.Hour).Format("15:04:05"),
		Hour:        1,
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/bookings", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	ctx, router := gin.CreateTestContext(rec)
	router.POST("/api/v1/bookings", suite.controller.CreateBookingHandler)
	ctx.Request = req

	router.ServeHTTP(rec, req)
	assert.Equal(suite.T(), http.StatusBadRequest, rec.Code)
	var response dto.SingleResponse
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "start time cant in the past", response.Status.Message)
}

func (suite *BookingControllerTestSuite) TestCreateBookingHandler_ServiceError() {
	now := time.Now()
	futureDate := now.AddDate(0, 0, 1).Format("02-01-2006")
	futureTime := now.Add(time.Hour).Format("15:04:05")

	payload := dto.CreateBookingRequest{
		CourtId:     "1",
		BookingDate: futureDate,
		StartTime:   futureTime,
		Hour:        1,
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/bookings", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	ctx, router := gin.CreateTestContext(rec)
	router.POST("/api/v1/bookings", suite.controller.CreateBookingHandler)
	ctx.Request = req

	ctx.Set("userId", "1")

	expectedError := errors.New("cannot book: court not available")
	suite.bookingServiceMock.On("Create", payload).Return(model.Booking{}, expectedError)

	router.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusBadRequest, rec.Code)
	suite.T().Logf("Response body: %s", rec.Body.String())
}

func (suite *BookingControllerTestSuite) TestCreateBookingHandler_InternalServerError() {
	now := time.Now()
	futureDate := now.AddDate(0, 0, 1).Format("02-01-2006")
	futureTime := now.Add(time.Hour).Format("15:04:05")

	payload := dto.CreateBookingRequest{
		CourtId:     "1",
		BookingDate: futureDate,
		StartTime:   futureTime,
		Hour:        1,
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/bookings", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	ctx, router := gin.CreateTestContext(rec)
	router.POST("/api/v1/bookings", suite.controller.CreateBookingHandler)
	ctx.Request = req

	ctx.Set("userId", "1")

	expectedError := errors.New("unexpected error")
	suite.bookingServiceMock.On("Create", payload).Return(model.Booking{}, expectedError)

	router.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, rec.Code)
	suite.T().Logf("Response body: %s", rec.Body.String())
}

// -
func (suite *BookingControllerTestSuite) TestNotificationHandler_Success() {
	mockNotification := dto.PaymentNotificationInput{
		OrderId: "order123",
	}

	mockPayloadJSON, err := json.Marshal(mockNotification)
	assert.NoError(suite.T(), err)

	record := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, "/api/v1/bookings/payment/notif", bytes.NewBuffer(mockPayloadJSON))
	assert.NoError(suite.T(), err)

	ctx, _ := gin.CreateTestContext(record)
	ctx.Request = req

	suite.bookingServiceMock.On("UpdatePayment", mockNotification).Return(nil)

	suite.controller.NotificationHandler(ctx)
	assert.Equal(suite.T(), http.StatusOK, record.Code)
}

func (suite *BookingControllerTestSuite) TestNotificationHandler_BindingError() {
	record := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/bookings/payment/notif", strings.NewReader("invalid JSON"))

	ctx, _ := gin.CreateTestContext(record)
	ctx.Request = req

	suite.controller.NotificationHandler(ctx)
	assert.Equal(suite.T(), http.StatusBadRequest, record.Code)
}

func (suite *BookingControllerTestSuite) TestNotificationHandler_ServiceError() {
	mockNotification := dto.PaymentNotificationInput{
		OrderId: "order123",
	}

	mockPayloadJSON, err := json.Marshal(mockNotification)
	assert.NoError(suite.T(), err)

	record := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, "/api/v1/bookings/payment/notif", bytes.NewBuffer(mockPayloadJSON))
	assert.NoError(suite.T(), err)

	ctx, _ := gin.CreateTestContext(record)
	ctx.Request = req

	suite.bookingServiceMock.On("UpdatePayment", mockNotification).Return(errors.New("update error"))

	suite.controller.NotificationHandler(ctx)
	assert.Equal(suite.T(), http.StatusInternalServerError, record.Code)
}

// -
func (suite *BookingControllerTestSuite) TestCreateRepayHandler_Success() {
	payload := dto.CreateRepayRequest{
		BookingId:     "1",
		EmployeeId:    "",
		PaymentMethod: "mid",
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/bookings/repayment", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	ctx, router := gin.CreateTestContext(rec)
	router.POST("/api/v1/bookings/repayment", suite.controller.CreateRepayHandler)
	ctx.Request = req

	ctx.Set("userId", "1")

	expectedRepayment := model.Payment{
		Id:            "1",
		BookingId:     "1",
		User:          model.User{},
		Court:         model.Court{},
		OrderId:       "1",
		Description:   "Description",
		PaymentMethod: "PaymentMethod",
		Price:         100000,
		Qty:           1,
		Status:        "Status",
		PaymentURL:    "PaymentURL",
	}

	suite.bookingServiceMock.On("CreateRepay", payload).Return(expectedRepayment, nil)

	router.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusOK, rec.Code)

	var response map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), nil, nil)

	data := response["data"].(map[string]interface{})
	assert.Equal(suite.T(), expectedRepayment.BookingId, data["bookingId"])

}

func (suite *BookingControllerTestSuite) TestCreateRepayHandler_BindingError() {
	record := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/bookings/repayment", strings.NewReader("invalid JSON"))

	ctx, _ := gin.CreateTestContext(record)
	ctx.Request = req

	suite.controller.CreateRepayHandler(ctx)
	assert.Equal(suite.T(), http.StatusBadRequest, record.Code)
}

func (suite *BookingControllerTestSuite) TestCreateRepayHandler_PaymentError() {
	mockPayload := dto.CreateRepayRequest{
		BookingId:     "1",
		EmployeeId:    "1",
		PaymentMethod: "gopay",
	}

	mockPayloadJSON, err := json.Marshal(mockPayload)
	assert.NoError(suite.T(), err)

	record := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, "/api/v1/bookings/repayment", bytes.NewBuffer(mockPayloadJSON))
	assert.NoError(suite.T(), err)

	ctx, _ := gin.CreateTestContext(record)
	ctx.Request = req

	suite.bookingServiceMock.On("CreateRepay", mockPayload).Return(model.Payment{}, errors.New("create error"))

	suite.controller.CreateRepayHandler(ctx)
	assert.Equal(suite.T(), http.StatusInternalServerError, record.Code)
}
func (suite *BookingControllerTestSuite) TestCreateRepayHandler_ServiceError() {

	payload := dto.CreateRepayRequest{
		BookingId:     "1",
		EmployeeId:    "",
		PaymentMethod: "mid",
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/bookings/repayment", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	ctx, router := gin.CreateTestContext(rec)
	router.POST("/api/v1/bookings/repayment", suite.controller.CreateRepayHandler)
	ctx.Request = req

	ctx.Set("userId", "1")

	suite.bookingServiceMock.On("CreateRepay", payload).Return(model.Payment{}, errors.New("error"))

	router.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, rec.Code)

	var response map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), nil, nil)
}

// --
func (suite *BookingControllerTestSuite) TestGetAllBookingsHandler_Success() {
	record := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/bookings?page=1&size=10", nil)

	ctx, router := gin.CreateTestContext(record)
	router.GET("/api/v1/bookings", suite.controller.GetAllBookingsHandler)
	ctx.Request = req

	suite.bookingServiceMock.On("FindAllBookings", 1, 10).Return([]model.Booking{mockBooking}, dto.Paginate{}, nil)

	router.ServeHTTP(record, req)
	assert.Equal(suite.T(), http.StatusOK, record.Code)
}

func (suite *BookingControllerTestSuite) TestGetAllBookingsHandler_InvalidPageOrSize() {
	record := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/bookings?page=invalid&size=invalid", nil)

	ctx, router := gin.CreateTestContext(record)
	router.GET("/api/v1/bookings", suite.controller.GetAllBookingsHandler)
	ctx.Request = req

	router.ServeHTTP(record, req)
	assert.Equal(suite.T(), http.StatusBadRequest, record.Code)
}

func (suite *BookingControllerTestSuite) TestGetAllBookingsHandler_ServiceError() {
	record := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/bookings?page=1&size=10", nil)

	ctx, router := gin.CreateTestContext(record)
	router.GET("/api/v1/bookings", suite.controller.GetAllBookingsHandler)
	ctx.Request = req

	suite.bookingServiceMock.On("FindAllBookings", 1, 10).Return([]model.Booking{}, dto.Paginate{}, errors.New("not found"))

	router.ServeHTTP(record, req)
	assert.Equal(suite.T(), http.StatusNotFound, record.Code)
}

func (suite *BookingControllerTestSuite) TestCheckBookingHandler_Succes() {
	record := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/bookings?page=1&size=10", nil)

	ctx, router := gin.CreateTestContext(record)
	router.GET("/api/v1/bookings", suite.controller.CheckBookingHandler)
	ctx.Request = req
	bookingDate := time.Time{}
	suite.bookingServiceMock.On("FindBookedCourt", bookingDate, 1, 10).Return([]model.Booking{mockBooking}, dto.Paginate{}, nil)

	router.ServeHTTP(record, req)
	assert.Equal(suite.T(), http.StatusOK, record.Code)
}

func (suite *BookingControllerTestSuite) TestCheckBookingHandler_Failed() {
	record := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/bookings?page=invalid&size=invalid", nil)

	ctx, router := gin.CreateTestContext(record)
	router.GET("/api/v1/bookings", suite.controller.CheckBookingHandler)
	ctx.Request = req

	router.ServeHTTP(record, req)
	assert.Equal(suite.T(), http.StatusBadRequest, record.Code)
}

func (suite *BookingControllerTestSuite) TestCheckBookingHandler_Failed2() {
	record := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/bookings?page=1&size=10", nil)

	ctx, router := gin.CreateTestContext(record)
	router.GET("/api/v1/bookings", suite.controller.CheckBookingHandler)
	ctx.Request = req

	bookingDate := time.Time{}
	suite.bookingServiceMock.On("FindBookedCourt", bookingDate, 1, 10).Return([]model.Booking{}, dto.Paginate{}, errors.New("not found"))

	router.ServeHTTP(record, req)
	assert.Equal(suite.T(), http.StatusInternalServerError, record.Code)
}

func (suite *BookingControllerTestSuite) TestCheckEndingHandler_Succes() {
	record := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/bookings?page=1&size=10", nil)

	ctx, router := gin.CreateTestContext(record)
	router.GET("/api/v1/bookings", suite.controller.CheckBookingTodayHandler)
	ctx.Request = req
	bookingDate := util.StringToDate(time.Now().Format("02-01-2006"))
	suite.bookingServiceMock.On("FindEndingBookings", bookingDate, 1, 10).Return([]model.Booking{mockBooking}, dto.Paginate{}, nil)

	router.ServeHTTP(record, req)
	assert.Equal(suite.T(), http.StatusOK, record.Code)
}

func (suite *BookingControllerTestSuite) TestCheckEndingHandler_Failed() {
	record := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/bookings?page=invalid&size=invalid", nil)

	ctx, router := gin.CreateTestContext(record)
	router.GET("/api/v1/bookings", suite.controller.CheckBookingTodayHandler)
	ctx.Request = req

	router.ServeHTTP(record, req)
	assert.Equal(suite.T(), http.StatusBadRequest, record.Code)
}

func (suite *BookingControllerTestSuite) TestCheckEndingHandler_Failed2() {
	record := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/bookings?page=1&size=10", nil)

	ctx, router := gin.CreateTestContext(record)
	router.GET("/api/v1/bookings", suite.controller.CheckBookingTodayHandler)
	ctx.Request = req

	bookingDate := util.StringToDate(time.Now().Format("02-01-2006"))

	suite.bookingServiceMock.On("FindEndingBookings", bookingDate, 1, 10).Return([]model.Booking{}, dto.Paginate{}, errors.New("not found"))

	router.ServeHTTP(record, req)
	assert.Equal(suite.T(), http.StatusInternalServerError, record.Code)
}

func (suite *BookingControllerTestSuite) TestPaymentReportHandler_Success() {
	record := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/bookings/payment-report?page=1&size=2&filter=daily&day=1&month=1&year=2024", nil)

	ctx, router := gin.CreateTestContext(record)
	router.GET("/api/v1/bookings/payment-report", suite.controller.PaymentReportHandler)
	ctx.Request = req

	mockPayments := []model.Payment{
		{
			Id:            "1",
			BookingId:     "1",
			User:          model.User{},
			Court:         model.Court{},
			OrderId:       "1",
			Description:   "Description",
			PaymentMethod: "gopay",
			Price:         60000,
			Qty:           1,
			Status:        "paid",
			PaymentURL:    "payment.com",
		},
	}

	paginate := dto.Paginate{
		Page:       1,
		Size:       2,
		TotalRows:  2,
		TotalPages: 1,
	}

	suite.bookingServiceMock.On("FindPaymentReport", 1, 1, 2024, 1, 2, "daily").Return(mockPayments, paginate, int64(120000), nil)

	router.ServeHTTP(record, req)
	assert.Equal(suite.T(), http.StatusOK, record.Code)
	var response dto.ReportPaginateResponse
	err := json.Unmarshal(record.Body.Bytes(), &response)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "success get data", response.Status.Message)
	assert.Equal(suite.T(), int64(120000), response.TotalIncome)
	assert.Equal(suite.T(), 1, len(response.Data))
}

func (suite *BookingControllerTestSuite) TestPaymentReportHandler_InvalidPageSize() {
	record := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/bookings/payment-report?page=invalid&size=invalid", nil)

	ctx, router := gin.CreateTestContext(record)
	router.GET("/api/v1/bookings/payment-report", suite.controller.PaymentReportHandler)
	ctx.Request = req

	router.ServeHTTP(record, req)
	assert.Equal(suite.T(), http.StatusBadRequest, record.Code)
	var response dto.SingleResponse
	err := json.Unmarshal(record.Body.Bytes(), &response)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "invalid page or size", response.Status.Message)
}

func (suite *BookingControllerTestSuite) TestPaymentReportHandler_InvalidFilter() {
	record := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/bookings/payment-report?page=1&size=2&filter=invalid", nil)

	ctx, router := gin.CreateTestContext(record)
	router.GET("/api/v1/bookings/payment-report", suite.controller.PaymentReportHandler)
	ctx.Request = req

	router.ServeHTTP(record, req)
	assert.Equal(suite.T(), http.StatusBadRequest, record.Code)
	var response dto.SingleResponse
	err := json.Unmarshal(record.Body.Bytes(), &response)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "invalid filter, use 'daily', 'monthly', 'yearly'", response.Status.Message)
}

func (suite *BookingControllerTestSuite) TestPaymentReportHandler_InternalServerError() {
	record := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/bookings/payment-report?page=1&size=2&filter=daily&day=1&month=1&year=2024", nil)

	ctx, router := gin.CreateTestContext(record)
	router.GET("/api/v1/bookings/payment-report", suite.controller.PaymentReportHandler)
	ctx.Request = req

	suite.bookingServiceMock.On("FindPaymentReport", 1, 1, 2024, 1, 2, "daily").Return([]model.Payment{}, dto.Paginate{}, int64(0), errors.New("internal error"))

	router.ServeHTTP(record, req)
	assert.Equal(suite.T(), http.StatusInternalServerError, record.Code)
	var response dto.SingleResponse
	err := json.Unmarshal(record.Body.Bytes(), &response)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "internal error", response.Status.Message)
}

func (suite *BookingControllerTestSuite) TestRoute() {
	assert.NotNil(suite.T(), suite.rg)
}
