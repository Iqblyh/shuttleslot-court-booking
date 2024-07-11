package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"math"
	"team2/shuttleslot/model"
	"testing"
	"time"

	"team2/shuttleslot/model/dto"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var mockBooking = model.Booking{
	Id: "1",
	Customer: model.User{
		Id:   "1",
		Name: "jojo",
	},
	Court: model.Court{
		Id:    "1",
		Name:  "field1",
		Price: 30000,
	},
	Employee: model.User{
		Id:   "2",
		Name: "lala",
	},
	BookingDate:   time.Time{},
	StartTime:     time.Time{},
	EndTime:       time.Time{},
	Total_Payment: 30000,
	Status:        "done",
	PaymentDetails: []model.Payment{
		{
			Id:            "1",
			BookingId:     "1",
			OrderId:       "23",
			Description:   "booking_field1",
			PaymentMethod: "gopay",
			Price:         30000,
			Qty:           1,
			Status:        "paid",
			PaymentURL:    "https://app.sandbox.midtrans.com/snap/v4/redirection/76655bb9-438d-42e1-b12b-2067044c200d",
		},
	},
	CreatedAt: time.Time{},
	UpdatedAt: time.Time{},
}
var mockPayment = model.Payment{
	Id:            "1",
	BookingId:     "1",
	User:          model.User{},
	Court:         model.Court{},
	OrderId:       "1",
	Description:   "booking_field1",
	PaymentMethod: "gopay",
	Price:         30000,
	Qty:           1,
	Status:        "paid",
	PaymentURL:    "https://app.sandbox.midtrans.com/snap/v4/redirection/76655bb9-438d-42e1-b12b-2067044c200d",
}

type BookingRepositoryTestSuite struct {
	suite.Suite
	mockDb  *sql.DB
	mockSql sqlmock.Sqlmock
	repo    BookingRepository
}

func (suite *BookingRepositoryTestSuite) SetupTest() {
	db, mock, err := sqlmock.New()
	assert.NoError(suite.T(), err)
	suite.mockDb = db
	suite.mockSql = mock
	suite.repo = NewBookingRepository(suite.mockDb)
}

func TestBookingRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(BookingRepositoryTestSuite))
}

func (suite *BookingRepositoryTestSuite) TestCreateBooking_Success() {
	suite.mockSql.ExpectBegin()

	rows := sqlmock.NewRows([]string{"id", "customer_id", "court_id", "booking_date", "start_time", "end_time", "total_payment", "status"}).AddRow(mockBooking.Id, mockBooking.Customer.Id, mockBooking.Court.Id, mockBooking.BookingDate, mockBooking.StartTime, mockBooking.EndTime, mockBooking.Total_Payment, mockBooking.Status)
	suite.mockSql.ExpectQuery("INSERT INTO bookings").WillReturnRows(rows)

	for _, mb := range mockBooking.PaymentDetails {
		rows := sqlmock.NewRows([]string{"id", "booking_id", "order_id", "description", "payment_method", "price", "status", "payment_url"}).AddRow(mb.Id, mb.BookingId, mb.OrderId, mb.Description, mb.PaymentMethod, mb.Price, mb.Status, mb.PaymentURL)
		suite.mockSql.ExpectQuery("INSERT INTO payments").WillReturnRows(rows)
	}

	suite.mockSql.ExpectCommit()
	actual, err := suite.repo.Create(mockBooking)
	assert.Nil(suite.T(), err)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), mockBooking.Id, actual.Id)

}

func (suite *BookingRepositoryTestSuite) TestCreatePayment_Failed() {
	suite.mockSql.ExpectBegin()

	rows := sqlmock.NewRows([]string{"id", "customer_id", "court_id", "booking_date", "start_time", "end_time", "total_payment", "status"}).AddRow(mockBooking.Id, mockBooking.Customer.Id, mockBooking.Court.Id, mockBooking.BookingDate, mockBooking.StartTime, mockBooking.EndTime, mockBooking.Total_Payment, mockBooking.Status)
	suite.mockSql.ExpectQuery("INSERT INTO bookings").WillReturnRows(rows)

	for _, mb := range mockBooking.PaymentDetails {
		fmt.Print(mb)
		suite.mockSql.ExpectQuery("INSERT INTO payments").WillReturnError(errors.New("Insert payments failed"))
		_, err := suite.repo.Create(mockBooking)
		assert.Error(suite.T(), err)
	}
}

func (suite *BookingRepositoryTestSuite) TestCreate_Failed() {
	suite.mockSql.ExpectBegin()

	sqlmock.NewRows([]string{"id", "customer_id", "court_id", "booking_date", "start_time", "end_time", "total_payment", "status"}).AddRow(mockBooking.Id, mockBooking.Customer.Id, mockBooking.Court.Id, mockBooking.BookingDate, mockBooking.StartTime, mockBooking.EndTime, mockBooking.Total_Payment, mockBooking.Status)
	suite.mockSql.ExpectQuery("INSERT INTO bookings").WillReturnError(errors.New("Insert payments failed"))

	_, err := suite.repo.Create(mockBooking)
	assert.Error(suite.T(), err)
}

func (suite *BookingRepositoryTestSuite) TestFindAll_Success() {
	page := 1
	size := 10
	offset := (page - 1) * size

	rows := sqlmock.NewRows([]string{"id", "customer_id", "court_id", "employee_id", "booking_date", "start_time", "end_time", "total_payment", "status", "created_at", "updated_at"}).
		AddRow(mockBooking.Id, mockBooking.Customer.Id, mockBooking.Court.Id, mockBooking.Employee.Id, mockBooking.BookingDate, mockBooking.StartTime, mockBooking.EndTime, mockBooking.Total_Payment, mockBooking.Status, mockBooking.CreatedAt, mockBooking.UpdatedAt)

	suite.mockSql.ExpectQuery("SELECT id, customer_id, court_id, employee_id, booking_date, start_time, end_time, total_payment, status, created_at, updated_at FROM bookings LIMIT \\$1 OFFSET \\$2").
		WithArgs(size, offset).
		WillReturnRows(rows)

	customerRows := sqlmock.NewRows([]string{"id", "name", "phone_number", "email", "username", "points", "role"}).
		AddRow(mockBooking.Customer.Id, mockBooking.Customer.Name, "08123456789", "jojo@example.com", "jojo123", 100, "customer")

	suite.mockSql.ExpectQuery("SELECT id, name, phone_number, email, username, points, role FROM users WHERE id = \\$1").
		WithArgs(mockBooking.Customer.Id).
		WillReturnRows(customerRows)

	employeeRows := sqlmock.NewRows([]string{"id", "name", "phone_number", "email", "username", "points", "role"}).
		AddRow(mockBooking.Employee.Id, mockBooking.Employee.Name, "08123456789", "lala@example.com", "lala123", 100, "employee")

	suite.mockSql.ExpectQuery("SELECT id, name, phone_number, email, username, points, role FROM users WHERE id = \\$1").
		WithArgs(mockBooking.Employee.Id).
		WillReturnRows(employeeRows)

	courtRows := sqlmock.NewRows([]string{"id", "name", "price"}).
		AddRow(mockBooking.Court.Id, mockBooking.Court.Name, mockBooking.Court.Price)

	suite.mockSql.ExpectQuery("SELECT id, name, price FROM courts WHERE id = \\$1").
		WithArgs(mockBooking.Court.Id).
		WillReturnRows(courtRows)

	actualBookings, paginate, err := suite.repo.FindAll(page, size)
	assert.Nil(suite.T(), err)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 1, len(actualBookings))
	assert.Equal(suite.T(), mockBooking.Id, actualBookings[0].Id)
	assert.Equal(suite.T(), dto.Paginate{
		Page:       page,
		Size:       size,
		TotalRows:  1,
		TotalPages: int(math.Ceil(float64(1) / float64(size))),
	}, paginate)
}

func (suite *BookingRepositoryTestSuite) TestFindAll_QueryError() {
	page := 1
	size := 10
	offset := (page - 1) * size

	suite.mockSql.ExpectQuery("SELECT id, customer_id, court_id, employee_id, booking_date, start_time, end_time, total_payment, status, created_at, updated_at FROM bookings LIMIT \\$1 OFFSET \\$2").
		WithArgs(size, offset).
		WillReturnError(errors.New("query error"))

	_, _, err := suite.repo.FindAll(page, size)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "query error", err.Error())
}

func (suite *BookingRepositoryTestSuite) TestFindAll_ScanError() {
	page := 1
	size := 10
	offset := (page - 1) * size

	rows := sqlmock.NewRows([]string{"id", "customer_id", "court_id", "employee_id", "booking_date", "start_time", "end_time", "total_payment", "status", "created_at", "updated_at"}).
		AddRow("wrong_id_type", mockBooking.Customer.Id, mockBooking.Court.Id, mockBooking.Employee.Id, mockBooking.BookingDate, mockBooking.StartTime, mockBooking.EndTime, mockBooking.Total_Payment, mockBooking.Status, mockBooking.CreatedAt, mockBooking.UpdatedAt)

	suite.mockSql.ExpectQuery("SELECT id, customer_id, court_id, employee_id, booking_date, start_time, end_time, total_payment, status, created_at, updated_at FROM bookings LIMIT \\$1 OFFSET \\$2").
		WithArgs(size, offset).
		WillReturnRows(rows)

	_, _, err := suite.repo.FindAll(page, size)
	assert.Error(suite.T(), err)
}

func (suite *BookingRepositoryTestSuite) TestFindAll_CustomerError() {
	page := 1
	size := 10
	offset := (page - 1) * size

	rows := sqlmock.NewRows([]string{"id", "customer_id", "court_id", "employee_id", "booking_date", "start_time", "end_time", "total_payment", "status", "created_at", "updated_at"}).
		AddRow(mockBooking.Id, mockBooking.Customer.Id, mockBooking.Court.Id, mockBooking.Employee.Id, mockBooking.BookingDate, mockBooking.StartTime, mockBooking.EndTime, mockBooking.Total_Payment, mockBooking.Status, mockBooking.CreatedAt, mockBooking.UpdatedAt)

	suite.mockSql.ExpectQuery("SELECT id, customer_id, court_id, employee_id, booking_date, start_time, end_time, total_payment, status, created_at, updated_at FROM bookings LIMIT \\$1 OFFSET \\$2").
		WithArgs(size, offset).
		WillReturnRows(rows)

	suite.mockSql.ExpectQuery("SELECT id, name, phone_number, email, username, points, role FROM users WHERE id = \\$1").
		WithArgs(mockBooking.Customer.Id).
		WillReturnError(errors.New("customer query error"))

	_, _, err := suite.repo.FindAll(page, size)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "customer query error", err.Error())
}

func (suite *BookingRepositoryTestSuite) TestFindByDate_Success() {
	bookingDate := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

	rows := sqlmock.NewRows([]string{"id", "customer_id", "court_id", "employee_id", "booking_date", "start_time", "end_time", "total_payment", "status", "created_at", "updated_at"}).
		AddRow(mockBooking.Id, mockBooking.Customer.Id, mockBooking.Court.Id, mockBooking.Employee.Id, mockBooking.BookingDate, mockBooking.StartTime, mockBooking.EndTime, mockBooking.Total_Payment, mockBooking.Status, mockBooking.CreatedAt, mockBooking.UpdatedAt)

	suite.mockSql.ExpectQuery("SELECT id, customer_id, court_id, employee_id, booking_date, start_time, end_time, total_payment, status, created_at, updated_at FROM bookings WHERE booking_date = \\$1").
		WithArgs(bookingDate).
		WillReturnRows(rows)

	actualBookings, err := suite.repo.FindByDate(bookingDate)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 1, len(actualBookings))
	assert.Equal(suite.T(), mockBooking.Id, actualBookings[0].Id)
}

func (suite *BookingRepositoryTestSuite) TestFindByDate_QueryError() {
	bookingDate := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

	suite.mockSql.ExpectQuery("SELECT id, customer_id, court_id, employee_id, booking_date, start_time, end_time, total_payment, status, created_at, updated_at FROM bookings WHERE booking_date = \\$1").
		WithArgs(bookingDate).
		WillReturnError(errors.New("query error"))

	_, err := suite.repo.FindByDate(bookingDate)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "query error", err.Error())
}

func (suite *BookingRepositoryTestSuite) TestFindByDate_ScanError() {
	bookingDate := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

	rows := sqlmock.NewRows([]string{"id", "customer_id", "court_id", "employee_id", "booking_date", "start_time", "end_time", "total_payment", "status", "created_at", "updated_at"}).
		AddRow("wrong_id_type", mockBooking.Customer.Id, mockBooking.Court.Id, mockBooking.Employee.Id, mockBooking.BookingDate, mockBooking.StartTime, mockBooking.EndTime, mockBooking.Total_Payment, mockBooking.Status, mockBooking.CreatedAt, mockBooking.UpdatedAt)

	suite.mockSql.ExpectQuery("SELECT id, customer_id, court_id, employee_id, booking_date, start_time, end_time, total_payment, ss, created_at, updated_at FROM bookings WHERE booking_date = \\$1").
		WithArgs(bookingDate).
		WillReturnRows(rows)

	_, err := suite.repo.FindByDate(bookingDate)
	assert.Error(suite.T(), err)
}

func (suite *BookingRepositoryTestSuite) TestFindById_Success() {
	bookingId := "1"

	row := sqlmock.NewRows([]string{"id", "customer_id", "court_id", "employee_id", "booking_date", "start_time", "end_time", "total_payment", "status", "created_at", "updated_at"}).
		AddRow(mockBooking.Id, mockBooking.Customer.Id, mockBooking.Court.Id, mockBooking.Employee.Id, mockBooking.BookingDate, mockBooking.StartTime, mockBooking.EndTime, mockBooking.Total_Payment, mockBooking.Status, mockBooking.CreatedAt, mockBooking.UpdatedAt)

	suite.mockSql.ExpectQuery("SELECT id, customer_id, court_id, employee_id, booking_date, start_time, end_time, total_payment, status, created_at, updated_at FROM bookings WHERE id = \\$1").
		WithArgs(bookingId).
		WillReturnRows(row)

	actualBooking, err := suite.repo.FindById(bookingId)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), mockBooking.Id, actualBooking.Id)
}

func (suite *BookingRepositoryTestSuite) TestFindById_QueryError() {
	bookingId := "1"

	suite.mockSql.ExpectQuery("SELECT id, customer_id, court_id, employee_id, booking_date, start_time, end_time, total_payment, status, created_at, updated_at FROM bookings WHERE id = \\$1").
		WithArgs(bookingId).
		WillReturnError(errors.New("query error"))

	_, err := suite.repo.FindById(bookingId)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "query error", err.Error())
}

func (suite *BookingRepositoryTestSuite) TestFindById_ScanError() {
	bookingId := "1"

	row := sqlmock.NewRows([]string{"id", "customer_id", "court_id", "employee_id", "booking_date", "start_time", "end_time", "total_payment", "status", "created_at", "updated_at"}).
		AddRow("wrong_id_type", mockBooking.Customer.Id, mockBooking.Court.Id, mockBooking.Employee.Id, mockBooking.BookingDate, mockBooking.StartTime, mockBooking.EndTime, mockBooking.Total_Payment, mockBooking.Status, mockBooking.CreatedAt, mockBooking.UpdatedAt)

	suite.mockSql.ExpectQuery("SELECT id, customer_id, court_id, employee_id, booking_date, start_time, end_time, total_payment, ss, created_at, updated_at FROM bookings WHERE id = \\$1").
		WithArgs(bookingId).
		WillReturnRows(row)

	_, err := suite.repo.FindById(bookingId)
	assert.Error(suite.T(), err)
}

func (suite *BookingRepositoryTestSuite) TestFindTotal_Success() {
	customerId := "1"
	totalBooking := 5

	row := sqlmock.NewRows([]string{"total_booking"}).AddRow(totalBooking)

	suite.mockSql.ExpectQuery("SELECT COUNT \\(\\*\\) AS total_booking FROM bookings WHERE customer_id = \\$1").
		WithArgs(customerId).
		WillReturnRows(row)

	actualTotal, err := suite.repo.FindTotal(customerId)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), totalBooking, actualTotal)
}

func (suite *BookingRepositoryTestSuite) TestFindTotal_QueryError() {
	customerId := "1"

	suite.mockSql.ExpectQuery("SELECT COUNT \\(\\*\\) AS total_booking FROM bookings WHERE customer_id = \\$1").
		WithArgs(customerId).
		WillReturnError(errors.New("query error"))

	_, err := suite.repo.FindTotal(customerId)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "query error", err.Error())
}

func (suite *BookingRepositoryTestSuite) TestFindTotal_ScanError() {
	customerId := "1"

	row := sqlmock.NewRows([]string{"total_booking"}).AddRow("wrong_type")

	suite.mockSql.ExpectQuery("SELECT COUNT \\(\\*\\) AS total_booking FROM bookings WHERE customer_id = \\$1").
		WithArgs(customerId).
		WillReturnRows(row)

	_, err := suite.repo.FindTotal(customerId)
	assert.Error(suite.T(), err)
}

func (suite *BookingRepositoryTestSuite) TestFindPaymentByOrderId_Success() {
	orderId := "23"
	expectedPayment := model.Payment{
		Id:            "1",
		BookingId:     "1",
		OrderId:       orderId,
		Description:   "booking_field1",
		PaymentMethod: "gopay",
		Price:         30000,
		Status:        "paid",
		PaymentURL:    "https://app.sandbox.midtrans.com/snap/v4/redirection/76655bb9-438d-42e1-b12b-2067044c200d",
	}

	row := sqlmock.NewRows([]string{"id", "booking_id", "order_id", "description", "payment_method", "price", "status", "payment_url"}).
		AddRow(expectedPayment.Id, expectedPayment.BookingId, expectedPayment.OrderId, expectedPayment.Description, expectedPayment.PaymentMethod, expectedPayment.Price, expectedPayment.Status, expectedPayment.PaymentURL)

	suite.mockSql.ExpectQuery("SELECT id, booking_id, order_id, description, payment_method, price, status, payment_url FROM payments WHERE order_id = \\$1").
		WithArgs(orderId).
		WillReturnRows(row)

	actualPayment, err := suite.repo.FindPaymentByOrderId(orderId)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedPayment, actualPayment)
}

func (suite *BookingRepositoryTestSuite) TestFindPaymentByOrderId_QueryError() {
	orderId := "23"

	suite.mockSql.ExpectQuery("SELECT id, booking_id, order_id, description, payment_method, price, status, payment_url FROM payments WHERE order_id = \\$1").
		WithArgs(orderId).
		WillReturnError(errors.New("query error"))

	_, err := suite.repo.FindPaymentByOrderId(orderId)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "query error", err.Error())
}

func (suite *BookingRepositoryTestSuite) TestFindPaymentByOrderId_ScanError() {
	orderId := "23"

	row := sqlmock.NewRows([]string{"id", "booking_id", "order_id", "description", "payment_method", "price", "status", "payment_url"}).
		AddRow("wrong_type", "wrong_type", "wrong_type", "wrong_type", "wrong_type", "wrong_type", "wrong_type", "wrong_type")

	suite.mockSql.ExpectQuery("SELECT id, booking_id, order_id, description, payment_method, price, status, payment_url FROM payments WHERE order_id = \\$1").
		WithArgs(orderId).
		WillReturnRows(row)

	_, err := suite.repo.FindPaymentByOrderId(orderId)
	assert.Error(suite.T(), err)
}

func (suite *BookingRepositoryTestSuite) TestUpdateStatus_Success_Pending() {
	payload := model.Payment{
		BookingId:     "1",
		OrderId:       "23",
		PaymentMethod: "gopay",
		Status:        "pending",
	}

	suite.mockSql.ExpectBegin()
	suite.mockSql.ExpectExec(
		"UPDATE payments SET payment_method = \\$1, updated_at = \\$2 WHERE order_id = \\$3",
	).WithArgs(payload.PaymentMethod, sqlmock.AnyArg(), payload.OrderId).WillReturnResult(sqlmock.NewResult(0, 1))
	suite.mockSql.ExpectCommit()

	err := suite.repo.UpdateStatus(payload)
	assert.NoError(suite.T(), err)
	assert.NoError(suite.T(), suite.mockSql.ExpectationsWereMet())
}

func (suite *BookingRepositoryTestSuite) TestUpdateStatus_Success_Paid() {
	payload := model.Payment{
		BookingId:     "1",
		OrderId:       "23",
		PaymentMethod: "gopay",
		Status:        "paid",
	}

	suite.mockSql.ExpectBegin()
	suite.mockSql.ExpectExec(
		"UPDATE payments SET payment_method = \\$1, status = \\$2, payment_url = \\$3, updated_at = \\$4 WHERE order_id = \\$5",
	).WithArgs(payload.PaymentMethod, payload.Status, "", sqlmock.AnyArg(), payload.OrderId).WillReturnResult(sqlmock.NewResult(0, 1))
	suite.mockSql.ExpectExec(
		"UPDATE bookings SET status = \\$1, updated_at = \\$2 WHERE id = \\$3",
	).WithArgs("booked", sqlmock.AnyArg(), payload.BookingId).WillReturnResult(sqlmock.NewResult(0, 1))
	suite.mockSql.ExpectCommit()

	err := suite.repo.UpdateStatus(payload)
	assert.NoError(suite.T(), err)
	assert.NoError(suite.T(), suite.mockSql.ExpectationsWereMet())
}

func (suite *BookingRepositoryTestSuite) TestUpdateStatus_Success_Cancel() {
	payload := model.Payment{
		BookingId: "1",
		OrderId:   "23",
		Status:    "cancel",
	}

	suite.mockSql.ExpectBegin()
	suite.mockSql.ExpectExec(
		"DELETE FROM payments WHERE order_id = \\$1",
	).WithArgs(payload.OrderId).WillReturnResult(sqlmock.NewResult(0, 1))
	suite.mockSql.ExpectExec(
		"UPDATE bookings SET status = \\$1, updated_at = \\$2 WHERE id = \\$3",
	).WithArgs("cancel", sqlmock.AnyArg(), payload.BookingId).WillReturnResult(sqlmock.NewResult(0, 1))
	suite.mockSql.ExpectCommit()

	err := suite.repo.UpdateStatus(payload)
	assert.NoError(suite.T(), err)
	assert.NoError(suite.T(), suite.mockSql.ExpectationsWereMet())
}

func (suite *BookingRepositoryTestSuite) TestUpdateStatus_Failure() {
	payload := model.Payment{
		BookingId:     "1",
		OrderId:       "23",
		PaymentMethod: "gopay",
		Status:        "paid",
	}

	suite.mockSql.ExpectBegin()
	suite.mockSql.ExpectExec(
		"UPDATE payments SET payment_method = \\$1, status = \\$2, payment_url = \\$3, updated_at = \\$4 WHERE order_id = \\$5",
	).WithArgs(payload.PaymentMethod, payload.Status, "", sqlmock.AnyArg(), payload.OrderId).WillReturnError(fmt.Errorf("update failed"))
	suite.mockSql.ExpectRollback()

	err := suite.repo.UpdateStatus(payload)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "update failed", err.Error())
	assert.NoError(suite.T(), suite.mockSql.ExpectationsWereMet())
}

func (suite *BookingRepositoryTestSuite) TestCreateRepay_Success() {
	payload := model.Payment{
		BookingId:     "1",
		OrderId:       "23",
		Description:   "booking_field1",
		PaymentMethod: "cash",
		Price:         30000,
		PaymentURL:    "https://app.sandbox.midtrans.com/snap/v4/redirection/76655bb9-438d-42e1-b12b-2067044c200d",
		User: model.User{
			Id:   "2",
			Name: "lala",
		},
	}

	suite.mockSql.ExpectBegin()
	suite.mockSql.ExpectQuery(
		"INSERT INTO payments \\(booking_id, order_id, description, payment_method, price, status, payment_url\\) VALUES \\(\\$1, \\$2, \\$3, \\$4, \\$5, \\$6, \\$7\\) RETURNING id, booking_id, order_id, description, payment_method, price, status, payment_url",
	).WithArgs(payload.BookingId, payload.OrderId, payload.Description, payload.PaymentMethod, payload.Price, "unpaid", payload.PaymentURL).WillReturnRows(sqlmock.NewRows([]string{"id", "booking_id", "order_id", "description", "payment_method", "price", "status", "payment_url"}).
		AddRow("1", payload.BookingId, payload.OrderId, payload.Description, payload.PaymentMethod, payload.Price, "unpaid", payload.PaymentURL))

	suite.mockSql.ExpectExec(
		"UPDATE bookings SET employee_id = \\$1, updated_at = \\$2 WHERE id = \\$3").
		WithArgs(payload.User.Id, sqlmock.AnyArg(), payload.BookingId).WillReturnResult(sqlmock.NewResult(0, 1))

	suite.mockSql.ExpectExec("UPDATE payments SET status = \\$1, updated_at = \\$2 WHERE id = \\$3").WithArgs("paid", sqlmock.AnyArg(), "1"). // Assuming payment ID '1' returned from insert
																			WillReturnResult(sqlmock.NewResult(0, 1))

	suite.mockSql.ExpectExec(
		"UPDATE bookings SET status = \\$1, updated_at = \\$2 WHERE id = \\$3",
	).WithArgs(
		"done", sqlmock.AnyArg(), payload.BookingId,
	).WillReturnResult(sqlmock.NewResult(0, 1))

	suite.mockSql.ExpectCommit()

	payment, err := suite.repo.CreateRepay(payload)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), payload.BookingId, payment.BookingId)
	assert.Equal(suite.T(), "paid", payment.Status)

	assert.NoError(suite.T(), suite.mockSql.ExpectationsWereMet())
}

func (suite *BookingRepositoryTestSuite) TestUpdateRepaymentStatus_Pending() {
	payload := model.Payment{
		OrderId:       "23",
		PaymentMethod: "gopay",
		Status:        "pending",
	}

	suite.mockSql.ExpectBegin()

	updatePayment := "UPDATE payments SET payment_method = \\$1, updated_at = \\$2 WHERE order_id = \\$3"
	suite.mockSql.ExpectExec(updatePayment).
		WithArgs(payload.PaymentMethod, sqlmock.AnyArg(), payload.OrderId).
		WillReturnResult(sqlmock.NewResult(1, 1))

	suite.mockSql.ExpectCommit()

	err := suite.repo.UpdateRepaymentStatus(payload)
	assert.NoError(suite.T(), err)
}

func (suite *BookingRepositoryTestSuite) TestUpdateRepaymentStatus_Paid() {
	payload := model.Payment{
		OrderId:       "23",
		PaymentMethod: "gopay",
		Status:        "paid",
		BookingId:     "1",
		User:          model.User{Id: "2"},
	}

	suite.mockSql.ExpectBegin()

	updatePayment := "UPDATE payments SET payment_method = \\$1, status = \\$2, payment_url = \\$3, updated_at = \\$4 WHERE order_id = \\$5"
	suite.mockSql.ExpectExec(updatePayment).
		WithArgs(payload.PaymentMethod, payload.Status, "", sqlmock.AnyArg(), payload.OrderId).
		WillReturnResult(sqlmock.NewResult(1, 1))

	updateBooking := "UPDATE bookings SET employee_id = \\$1, status = \\$2, updated_at = \\$3 WHERE id = \\$4"
	suite.mockSql.ExpectExec(updateBooking).
		WithArgs(payload.User.Id, "done", sqlmock.AnyArg(), payload.BookingId).
		WillReturnResult(sqlmock.NewResult(1, 1))

	suite.mockSql.ExpectCommit()

	err := suite.repo.UpdateRepaymentStatus(payload)
	assert.NoError(suite.T(), err)
}

func (suite *BookingRepositoryTestSuite) TestUpdateRepaymentStatus_PendingError() {
	payload := model.Payment{
		OrderId:       "23",
		PaymentMethod: "gopay",
		Status:        "pending",
	}

	suite.mockSql.ExpectBegin()

	updatePayment := "UPDATE payments SET payment_method = \\$1, updated_at = \\$2 WHERE order_id = \\$3"
	suite.mockSql.ExpectExec(updatePayment).
		WithArgs(payload.PaymentMethod, sqlmock.AnyArg(), payload.OrderId).
		WillReturnError(errors.New("update payment error"))

	suite.mockSql.ExpectRollback()

	err := suite.repo.UpdateRepaymentStatus(payload)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "update payment error", err.Error())
}

func (suite *BookingRepositoryTestSuite) TestUpdateRepaymentStatus_PaidPaymentError() {
	payload := model.Payment{
		OrderId:       "23",
		PaymentMethod: "gopay",
		Status:        "paid",
		BookingId:     "1",
		User:          model.User{Id: "2"},
	}

	suite.mockSql.ExpectBegin()

	updatePayment := "UPDATE payments SET payment_method = \\$1, status = \\$2, payment_url = \\$3, updated_at = \\$4 WHERE order_id = \\$5"
	suite.mockSql.ExpectExec(updatePayment).
		WithArgs(payload.PaymentMethod, payload.Status, "", sqlmock.AnyArg(), payload.OrderId).
		WillReturnError(errors.New("update payment error"))

	suite.mockSql.ExpectRollback()

	err := suite.repo.UpdateRepaymentStatus(payload)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "update payment error", err.Error())
}

func (suite *BookingRepositoryTestSuite) TestUpdateRepaymentStatus_PaidBookingError() {
	payload := model.Payment{
		OrderId:       "23",
		PaymentMethod: "gopay",
		Status:        "paid",
		BookingId:     "1",
		User:          model.User{Id: "2"},
	}

	suite.mockSql.ExpectBegin()

	updatePayment := "UPDATE payments SET payment_method = \\$1, status = \\$2, payment_url = \\$3, updated_at = \\$4 WHERE order_id = \\$5"
	suite.mockSql.ExpectExec(updatePayment).
		WithArgs(payload.PaymentMethod, payload.Status, "", sqlmock.AnyArg(), payload.OrderId).
		WillReturnResult(sqlmock.NewResult(1, 1))

	updateBooking := "UPDATE bookings SET employee_id = \\$1, status = \\$2, updated_at = \\$3 WHERE id = \\$4"
	suite.mockSql.ExpectExec(updateBooking).
		WithArgs(payload.User.Id, "done", sqlmock.AnyArg(), payload.BookingId).
		WillReturnError(errors.New("update booking error"))

	suite.mockSql.ExpectRollback()

	err := suite.repo.UpdateRepaymentStatus(payload)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "update booking error", err.Error())
}

func (suite *BookingRepositoryTestSuite) TestFindBooked_Success() {
	page := 1
	size := 10
	offset := (page - 1) * size
	bookingDate := time.Now()

	rows := sqlmock.NewRows([]string{"court_id", "booking_date", "start_time", "end_time", "status"}).
		AddRow(mockBooking.Court.Id, bookingDate, mockBooking.StartTime, mockBooking.EndTime, mockBooking.Status).
		AddRow(mockBooking.Court.Id, bookingDate, mockBooking.StartTime, mockBooking.EndTime, mockBooking.Status)

	suite.mockSql.ExpectQuery("SELECT court_id, booking_date, start_time, end_time, status FROM bookings WHERE booking_date = \\$1 AND status IN \\('pending', 'booked'\\) LIMIT \\$2 OFFSET \\$3").
		WithArgs(bookingDate, size, offset).
		WillReturnRows(rows)

	actualBookings, paginate, err := suite.repo.FindBooked(bookingDate, page, size)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 2, len(actualBookings))
	assert.Equal(suite.T(), mockBooking.Court.Id, actualBookings[0].Court.Id)
	assert.Equal(suite.T(), dto.Paginate{
		Page:       page,
		Size:       size,
		TotalRows:  2,
		TotalPages: int(math.Ceil(float64(2) / float64(size))),
	}, paginate)
}

func (suite *BookingRepositoryTestSuite) TestFindBooked_QueryError() {
	page := 1
	size := 10
	bookingDate := time.Now()

	suite.mockSql.ExpectQuery("SELECT court_id, booking_date, start_time, end_time, status FROM bookings WHERE booking_date = \\$1 AND status IN \\('pending', 'booked'\\) LIMIT \\$2 OFFSET \\$3").
		WithArgs(bookingDate, size, (page-1)*size).
		WillReturnError(errors.New("query error"))

	_, _, err := suite.repo.FindBooked(bookingDate, page, size)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "query error", err.Error())
}

func (suite *BookingRepositoryTestSuite) TestFindBooked_ScanError() {
	page := 1
	size := 10
	bookingDate := time.Now()

	rows := sqlmock.NewRows([]string{"court_id", "booking_date", "start_time", "end_time", "status"}).
		AddRow("invalid_id", bookingDate, mockBooking.StartTime, mockBooking.EndTime, mockBooking.Status)

	suite.mockSql.ExpectQuery("SELECT invalid_id, booking_date, start_time, end_time, status FROM bookings WHERE booking_date = \\$1 AND status IN \\('pending', 'booked'\\) LIMIT \\$2 OFFSET \\$3").
		WithArgs(bookingDate, size, (page-1)*size).
		WillReturnRows(rows)

	_, _, err := suite.repo.FindBooked(bookingDate, page, size)
	assert.Error(suite.T(), err)
}

func (suite *BookingRepositoryTestSuite) TestFindEnding_Success() {
	page := 1
	size := 10
	offset := (page - 1) * size
	bookingDate := time.Now()

	rows := sqlmock.NewRows([]string{"id", "customer_id", "court_id", "booking_date", "start_time", "end_time", "total_payment", "status"}).
		AddRow(mockBooking.Id, mockBooking.Customer.Id, mockBooking.Court.Id, bookingDate, mockBooking.StartTime, mockBooking.EndTime, mockBooking.Total_Payment, "booked").
		AddRow(mockBooking.Id, mockBooking.Customer.Id, mockBooking.Court.Id, bookingDate, mockBooking.StartTime, mockBooking.EndTime, mockBooking.Total_Payment, "booked")

	suite.mockSql.ExpectQuery("SELECT id, customer_id, court_id, booking_date, start_time, end_time, total_payment, status FROM bookings WHERE booking_date = \\$1 AND status = 'booked' LIMIT \\$2 OFFSET \\$3").
		WithArgs(bookingDate, size, offset).
		WillReturnRows(rows)

	actualBookings, paginate, err := suite.repo.FindEnding(bookingDate, page, size)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 2, len(actualBookings))
	assert.Equal(suite.T(), mockBooking.Id, actualBookings[0].Id)
	assert.Equal(suite.T(), dto.Paginate{
		Page:       page,
		Size:       size,
		TotalRows:  2,
		TotalPages: int(math.Ceil(float64(2) / float64(size))),
	}, paginate)
}

func (suite *BookingRepositoryTestSuite) TestFindEnding_QueryError() {
	page := 1
	size := 10
	bookingDate := time.Now()

	suite.mockSql.ExpectQuery("SELECT id, customer_id, court_id, booking_date, start_time, end_time, total_payment, status FROM bookings WHERE booking_date = \\$1 AND status = 'booked' LIMIT \\$2 OFFSET \\$3").
		WithArgs(bookingDate, size, (page-1)*size).
		WillReturnError(errors.New("query error"))

	_, _, err := suite.repo.FindEnding(bookingDate, page, size)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "query error", err.Error())
}

func (suite *BookingRepositoryTestSuite) TestFindEnding_ScanError() {
	page := 1
	size := 10
	bookingDate := time.Now()

	rows := sqlmock.NewRows([]string{"id", "customer_id", "court_id", "booking_date", "start_time", "end_time", "total_payment", "status"}).
		AddRow("invalid_id", mockBooking.Customer.Id, mockBooking.Court.Id, bookingDate, mockBooking.StartTime, mockBooking.EndTime, mockBooking.Total_Payment, "booked")

	suite.mockSql.ExpectQuery("SELECT invalid_id, customer_id, court_id, booking_date, start_time, end_time, total_payment, status FROM bookings WHERE booking_date = \\$1 AND status = 'booked' LIMIT \\$2 OFFSET \\$3").
		WithArgs(bookingDate, size, (page-1)*size).
		WillReturnRows(rows)

	_, _, err := suite.repo.FindEnding(bookingDate, page, size)
	assert.Error(suite.T(), err)
}

func (suite *BookingRepositoryTestSuite) TestFindPaymentReport_Daily_Success() {
	day := 1
	month := 7
	year := 2024
	page := 1
	size := 10
	filterType := "daily"
	offset := (page - 1) * size

	rows := sqlmock.NewRows([]string{"id", "booking_id", "order_id", "description", "payment_method", "price"}).
		AddRow(mockPayment.Id, mockPayment.BookingId, mockPayment.OrderId, mockPayment.Description, mockPayment.PaymentMethod, mockPayment.Price)

	expectedQuery := `SELECT id, booking_id, order_id, description, payment_method, price FROM payments WHERE EXTRACT\(DAY FROM created_at\) = \$1  AND EXTRACT\(MONTH FROM created_at\) = \$2 AND EXTRACT\(YEAR FROM created_at\) = \$3 LIMIT \$4 OFFSET \$5`

	suite.mockSql.ExpectQuery(expectedQuery).
		WithArgs(day, month, year, size, offset).
		WillReturnRows(rows)

	actualPayments, paginate, totalIncome, err := suite.repo.FindPaymentReport(day, month, year, page, size, filterType)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 1, len(actualPayments))
	assert.Equal(suite.T(), mockPayment.Id, actualPayments[0].Id)
	assert.Equal(suite.T(), dto.Paginate{
		Page:       page,
		Size:       size,
		TotalRows:  1,
		TotalPages: int(math.Ceil(float64(1) / float64(size))),
	}, paginate)
	assert.Equal(suite.T(), int64(mockPayment.Price), totalIncome)
}

func (suite *BookingRepositoryTestSuite) TestFindPaymentReport_Monthly_Success() {
	month := 7
	year := 2024
	page := 1
	size := 10
	filterType := "monthly"
	offset := (page - 1) * size

	rows := sqlmock.NewRows([]string{"id", "booking_id", "order_id", "description", "payment_method", "price"}).
		AddRow(mockPayment.Id, mockPayment.BookingId, mockPayment.OrderId, mockPayment.Description, mockPayment.PaymentMethod, mockPayment.Price)

	expectedQuery := `SELECT id, booking_id, order_id, description, payment_method, price FROM payments WHERE EXTRACT\(MONTH FROM created_at\) = \$1 AND EXTRACT\(YEAR FROM created_at\) = \$2 LIMIT \$3 OFFSET \$4`

	suite.mockSql.ExpectQuery(expectedQuery).
		WithArgs(month, year, size, offset).
		WillReturnRows(rows)

	actualPayments, paginate, totalIncome, err := suite.repo.FindPaymentReport(0, month, year, page, size, filterType)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 1, len(actualPayments))
	assert.Equal(suite.T(), mockPayment.Id, actualPayments[0].Id)
	assert.Equal(suite.T(), dto.Paginate{
		Page:       page,
		Size:       size,
		TotalRows:  1,
		TotalPages: int(math.Ceil(float64(1) / float64(size))),
	}, paginate)
	assert.Equal(suite.T(), int64(mockPayment.Price), totalIncome)
}

func (suite *BookingRepositoryTestSuite) TestFindPaymentReport_Yearly_Success() {
	year := 2024
	page := 1
	size := 10
	filterType := "yearly"
	offset := (page - 1) * size

	rows := sqlmock.NewRows([]string{"id", "booking_id", "order_id", "description", "payment_method", "price"}).
		AddRow(mockPayment.Id, mockPayment.BookingId, mockPayment.OrderId, mockPayment.Description, mockPayment.PaymentMethod, mockPayment.Price)

	expectedQuery := `SELECT id, booking_id, order_id, description, payment_method, price FROM payments WHERE EXTRACT\(YEAR FROM created_at\) = \$1 LIMIT \$2 OFFSET \$3`

	suite.mockSql.ExpectQuery(expectedQuery).
		WithArgs(year, size, offset).
		WillReturnRows(rows)

	actualPayments, paginate, totalIncome, err := suite.repo.FindPaymentReport(0, 0, year, page, size, filterType)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 1, len(actualPayments))
	assert.Equal(suite.T(), mockPayment.Id, actualPayments[0].Id)
	assert.Equal(suite.T(), dto.Paginate{
		Page:       page,
		Size:       size,
		TotalRows:  1,
		TotalPages: int(math.Ceil(float64(1) / float64(size))),
	}, paginate)
	assert.Equal(suite.T(), int64(mockPayment.Price), totalIncome)
}

func (suite *BookingRepositoryTestSuite) TestFindPaymentReport_QueryError() {
	day := 1
	month := 7
	year := 2024
	page := 1
	size := 10
	filterType := "daily"
	offset := (page - 1) * size

	suite.mockSql.ExpectQuery(`SELECT id, booking_id, order_id, description, payment_method, price FROM payments WHERE EXTRACT\(DAY FROM created_at\) = \$1 AND EXTRACT\(MONTH FROM created_at\) = \$2 AND EXTRACT\(YEAR FROM created_at\) = \$3 LIMIT \$4 OFFSET \$5`).
		WithArgs(day, month, year, size, offset).
		WillReturnError(errors.New("query error"))

	_, _, _, err := suite.repo.FindPaymentReport(day, month, year, page, size, filterType)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "query error")
}

func (suite *BookingRepositoryTestSuite) TestFindPaymentReport_ScanError() {
	day := 1
	month := 7
	year := 2024
	page := 1
	size := 10
	filterType := "daily"
	offset := (page - 1) * size

	rows := sqlmock.NewRows([]string{"id", "booking_id", "order_id", "description", "payment_method", "price"}).
		AddRow("invalid_id", mockPayment.BookingId, mockPayment.OrderId, mockPayment.Description, mockPayment.PaymentMethod, mockPayment.Price)

	suite.mockSql.ExpectQuery("SELECT id, booking_id, order_id, description, payment_method, price FROM payments WHERE EXTRACT(DAY FROM created_at) = \\$1  AND EXTRACT(MONTH FROM created_at) = \\$2 AND EXTRACT(YEAR FROM created_at) = \\$3 LIMIT \\$4 OFFSET \\$5").
		WithArgs(day, month, year, size, offset).
		WillReturnRows(rows)

	_, _, _, err := suite.repo.FindPaymentReport(day, month, year, page, size, filterType)
	assert.Error(suite.T(), err)
}
