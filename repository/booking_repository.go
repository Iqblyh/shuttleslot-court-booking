package repository

import (
	"database/sql"
	"math"
	"team2/shuttleslot/model"
	"team2/shuttleslot/model/dto"
	"time"
)

type bookingRepository struct {
	DB *sql.DB
}

type BookingRepository interface {
	Create(payload model.Booking) (model.Booking, error)
	FindAll(page int, size int) ([]model.Booking, dto.Paginate, error)
	FindByDate(bookingDate time.Time) ([]model.Booking, error)
	FindById(bookingId string) (model.Booking, error)
	FindTotal(customerId string) (int, error)
	FindPaymentByOrderId(order_id string) (model.Payment, error)
	UpdateStatus(payload model.Payment) error
	CreateRepay(payload model.Payment) (model.Payment, error)
	UpdateRepaymentStatus(payload model.Payment) error
	FindBooked(bookingDate time.Time, page int, size int) ([]model.Booking, dto.Paginate, error)
	FindEnding(bookingDate time.Time, page int, size int) ([]model.Booking, dto.Paginate, error)
	FindPaymentReport(day, month, year, page, size int, filterType string) ([]model.Payment, dto.Paginate, int64, error)
	// UpdateCancel(orderId string) error // Update Cancel Masih ku coba, skip dulu unit testingnya
}

func (r *bookingRepository) Create(payload model.Booking) (model.Booking, error) {
	transaction, _ := r.DB.Begin()

	var booking model.Booking
	query := "INSERT INTO bookings (customer_id, court_id, booking_date, start_time, end_time, total_payment, status) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, customer_id, court_id, booking_date, start_time, end_time, total_payment, status"

	err := transaction.QueryRow(query, payload.Customer.Id, payload.Court.Id, payload.BookingDate, payload.StartTime, payload.EndTime, payload.Total_Payment, "pending").Scan(
		&booking.Id,
		&booking.Customer.Id,
		&booking.Court.Id,
		&booking.BookingDate,
		&booking.StartTime,
		&booking.EndTime,
		&booking.Total_Payment,
		&booking.Status,
	)

	if err != nil {
		transaction.Rollback()
		return booking, err
	}

	var payment model.Payment

	depoPrice := booking.Total_Payment / 2

	query = "INSERT INTO payments (booking_id, order_id, description, payment_method, price, status, payment_url) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, booking_id, order_id, description, payment_method, price, status, payment_url"
	err = transaction.QueryRow(
		query,
		booking.Id,
		payload.PaymentDetails[0].OrderId,
		payload.PaymentDetails[0].Description,
		"mid",
		depoPrice,
		"unpaid",
		payload.PaymentDetails[0].PaymentURL,
	).Scan(
		&payment.Id,
		&payment.BookingId,
		&payment.OrderId,
		&payment.Description,
		&payment.PaymentMethod,
		&payment.Price,
		&payment.Status,
		&payment.PaymentURL,
	)

	if err != nil {
		transaction.Rollback()
		return booking, err
	}

	booking.PaymentDetails = append(booking.PaymentDetails, payment)

	transaction.Commit()
	return booking, nil
}

func (r *bookingRepository) FindAll(page int, size int) ([]model.Booking, dto.Paginate, error) {
	var bookings []model.Booking

	offset := (page - 1) * size

	query := "SELECT id, customer_id, court_id, employee_id, booking_date, start_time, end_time, total_payment, status, created_at, updated_at FROM bookings LIMIT $1 OFFSET $2"

	rows, err := r.DB.Query(query, size, offset)
	if err != nil {
		return []model.Booking{}, dto.Paginate{}, err
	}

	totalRows := 0
	var nullEmployee sql.NullString

	for rows.Next() {
		var b model.Booking
		if err := rows.Scan(
			&b.Id,
			&b.Customer.Id,
			&b.Court.Id,
			&nullEmployee,
			&b.BookingDate,
			&b.StartTime,
			&b.EndTime,
			&b.Total_Payment,
			&b.Status,
			&b.CreatedAt,
			&b.UpdatedAt,
		); err != nil {
			return []model.Booking{}, dto.Paginate{}, err
		}
		var customer model.User
		var employee model.User
		var court model.Court

		err = r.DB.QueryRow("SELECT id, name, phone_number, email, username, points, role FROM users WHERE id = $1", b.Customer.Id).Scan(
			&customer.Id,
			&customer.Name,
			&customer.PhoneNumber,
			&customer.Email,
			&customer.Username,
			&customer.Point,
			&customer.Role,
		)
		if err != nil {
			return []model.Booking{}, dto.Paginate{}, err
		}

		if nullEmployee.Valid {
			b.Employee.Id = nullEmployee.String

			err = r.DB.QueryRow("SELECT id, name, phone_number, email, username, points, role FROM users WHERE id = $1", b.Employee.Id).Scan(
				&employee.Id,
				&employee.Name,
				&employee.PhoneNumber,
				&employee.Email,
				&employee.Username,
				&employee.Point,
				&employee.Role,
			)
			if err != nil {
				return []model.Booking{}, dto.Paginate{}, err
			}
		}

		err := r.DB.QueryRow("SELECT id, name, price FROM courts WHERE id = $1", b.Court.Id).Scan(
			&court.Id,
			&court.Name,
			&court.Price,
		)
		if err != nil {
			return []model.Booking{}, dto.Paginate{}, err
		}

		b.Customer = customer
		b.Employee = employee
		b.Court = court

		bookings = append(bookings, b)
		totalRows++
	}

	paginate := dto.Paginate{
		Page:       page,
		Size:       size,
		TotalRows:  totalRows,
		TotalPages: int(math.Ceil(float64(totalRows) / float64(size))),
	}

	return bookings, paginate, nil
}

func (r *bookingRepository) FindByDate(bookingDate time.Time) ([]model.Booking, error) {
	var bookings []model.Booking

	query := "SELECT id, customer_id, court_id, employee_id, booking_date, start_time, end_time, total_payment, status, created_at, updated_at FROM bookings WHERE booking_date = $1"
	rows, err := r.DB.Query(query, bookingDate)
	if err != nil {
		return []model.Booking{}, err
	}

	var nullEmployee sql.NullString
	for rows.Next() {
		var b model.Booking
		if err := rows.Scan(
			&b.Id,
			&b.Customer.Id,
			&b.Court.Id,
			&nullEmployee,
			&b.BookingDate,
			&b.StartTime,
			&b.EndTime,
			&b.Total_Payment,
			&b.Status,
			&b.CreatedAt,
			&b.UpdatedAt,
		); err != nil {
			return []model.Booking{}, err
		}
		if nullEmployee.Valid {
			b.Employee.Id = nullEmployee.String
		}
		bookings = append(bookings, b)
	}

	return bookings, nil
}

func (r *bookingRepository) FindById(bookingId string) (model.Booking, error) {
	var booking model.Booking

	query := "SELECT id, customer_id, court_id, employee_id, booking_date, start_time, end_time, total_payment, status, created_at, updated_at FROM bookings WHERE id = $1"

	var employeeId sql.NullString

	err := r.DB.QueryRow(query, bookingId).Scan(
		&booking.Id,
		&booking.Customer.Id,
		&booking.Court.Id,
		&employeeId,
		&booking.BookingDate,
		&booking.StartTime,
		&booking.EndTime,
		&booking.Total_Payment,
		&booking.Status,
		&booking.CreatedAt,
		&booking.UpdatedAt,
	)

	if err != nil {
		return booking, err
	}

	if employeeId.Valid {
		booking.Employee.Id = employeeId.String
	}

	return booking, nil
}

func (r *bookingRepository) FindTotal(customerId string) (int, error) {
	var totalBooking int
	query := "SELECT COUNT (*) AS total_booking FROM bookings WHERE customer_id = $1"

	err := r.DB.QueryRow(query, customerId).Scan(&totalBooking)
	if err != nil {
		return 0, err
	}

	return totalBooking, nil
}

func (r *bookingRepository) FindPaymentByOrderId(order_id string) (model.Payment, error) {
	var payment model.Payment

	query := "SELECT id, booking_id, order_id, description, payment_method, price, status, payment_url FROM payments WHERE order_id = $1"

	err := r.DB.QueryRow(query, order_id).Scan(
		&payment.Id,
		&payment.BookingId,
		&payment.OrderId,
		&payment.Description,
		&payment.PaymentMethod,
		&payment.Price,
		&payment.Status,
		&payment.PaymentURL,
	)
	if err != nil {
		return model.Payment{}, err
	}

	return payment, nil
}

func (r *bookingRepository) UpdateStatus(payload model.Payment) error {
	transaction, _ := r.DB.Begin()

	if payload.Status == "pending" {
		updatePayment := "UPDATE payments SET payment_method = $1, updated_at = $2 WHERE order_id = $3"

		_, err := transaction.Exec(updatePayment, payload.PaymentMethod, time.Now(), payload.OrderId)
		if err != nil {
			transaction.Rollback()
			return err
		}
	}

	if payload.Status == "paid" {
		updatePayment := "UPDATE payments SET payment_method = $1, status = $2, payment_url = $3, updated_at = $4 WHERE order_id = $5"

		_, err := transaction.Exec(updatePayment, payload.PaymentMethod, payload.Status, "", time.Now(), payload.OrderId)
		if err != nil {
			transaction.Rollback()
			return err
		}

		updateBooking := "UPDATE bookings SET status = $1, updated_at = $2 WHERE id = $3"

		_, err = transaction.Exec(updateBooking, "booked", time.Now(), payload.BookingId)

		if err != nil {
			transaction.Rollback()
			return err
		}
	}

	if payload.Status == "cancel" {
		updatePayment := "DELETE FROM payments WHERE order_id = $1"

		_, err := transaction.Exec(updatePayment, payload.OrderId)

		if err != nil {
			transaction.Rollback()
			return err
		}

		updateBooking := "UPDATE bookings SET status = $1, updated_at = $2 WHERE id = $3"

		_, err = transaction.Exec(updateBooking, "cancel", time.Now(), payload.BookingId)

		if err != nil {
			transaction.Rollback()
			return err
		}
	}

	transaction.Commit()
	return nil
}

func (r *bookingRepository) CreateRepay(payload model.Payment) (model.Payment, error) {
	transaction, _ := r.DB.Begin()

	var payment model.Payment

	query := "INSERT INTO payments (booking_id, order_id, description, payment_method, price, status, payment_url) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, booking_id, order_id, description, payment_method, price, status, payment_url"

	err := transaction.QueryRow(
		query,
		payload.BookingId,
		payload.OrderId,
		payload.Description,
		payload.PaymentMethod,
		payload.Price,
		"unpaid",
		payload.PaymentURL,
	).Scan(
		&payment.Id,
		&payment.BookingId,
		&payment.OrderId,
		&payment.Description,
		&payment.PaymentMethod,
		&payment.Price,
		&payment.Status,
		&payment.PaymentURL,
	)
	if err != nil {
		transaction.Rollback()
		return payment, err
	}

	updateBooking := "UPDATE bookings SET employee_id = $1, updated_at = $2 WHERE id = $3"

	_, err = transaction.Exec(updateBooking, payload.User.Id, time.Now(), payload.BookingId)
	if err != nil {
		transaction.Rollback()
		return model.Payment{}, nil
	}

	if payment.PaymentMethod == "cash" {
		updatePayment := "UPDATE payments SET status = $1, updated_at = $2 WHERE id = $3"

		_, err := transaction.Exec(updatePayment, "paid", time.Now(), payment.Id)
		if err != nil {
			transaction.Rollback()
			return model.Payment{}, nil
		}

		updateBooking := "UPDATE bookings SET status = $1, updated_at = $2 WHERE id = $3"

		_, err = transaction.Exec(updateBooking, "done", time.Now(), payload.BookingId)
		if err != nil {
			transaction.Rollback()
			return model.Payment{}, nil
		}
	}

	transaction.Commit()
	return payment, nil
}

func (r *bookingRepository) UpdateRepaymentStatus(payload model.Payment) error {
	transaction, _ := r.DB.Begin()

	if payload.Status == "pending" {
		updatePayment := "UPDATE payments SET payment_method = $1, updated_at = $2 WHERE order_id = $3"

		_, err := transaction.Exec(updatePayment, payload.PaymentMethod, time.Now(), payload.OrderId)
		if err != nil {
			transaction.Rollback()
			return err
		}
	}

	if payload.Status == "paid" {
		updatePayment := "UPDATE payments SET payment_method = $1, status = $2, payment_url = $3, updated_at = $4 WHERE order_id = $5"

		_, err := transaction.Exec(updatePayment, payload.PaymentMethod, payload.Status, "", time.Now(), payload.OrderId)
		if err != nil {
			transaction.Rollback()
			return err
		}

		updateBooking := "UPDATE bookings SET employee_id = $1, status = $2, updated_at = $3 WHERE id = $4"

		_, err = transaction.Exec(updateBooking, payload.User.Id, "done", time.Now(), payload.BookingId)
		if err != nil {
			transaction.Rollback()
			return err
		}
	}

	transaction.Commit()
	return nil
}

func (r *bookingRepository) FindBooked(bookingDate time.Time, page int, size int) ([]model.Booking, dto.Paginate, error) {
	var bookings []model.Booking

	query := "SELECT court_id, booking_date, start_time, end_time, status FROM bookings WHERE booking_date = $1 AND status IN ('pending', 'booked') LIMIT $2 OFFSET $3"

	offset := (page - 1) * size
	rows, err := r.DB.Query(query, bookingDate, size, offset)
	if err != nil {
		return []model.Booking{}, dto.Paginate{}, err
	}

	totalRows := 0

	for rows.Next() {
		var b model.Booking
		if err := rows.Scan(
			&b.Court.Id,
			&b.BookingDate,
			&b.StartTime,
			&b.EndTime,
			&b.Status,
		); err != nil {
			return []model.Booking{}, dto.Paginate{}, err
		}

		bookings = append(bookings, b)
		totalRows++
	}

	paginate := dto.Paginate{
		Page:       page,
		Size:       size,
		TotalRows:  totalRows,
		TotalPages: int(math.Ceil(float64(totalRows) / float64(size))),
	}

	return bookings, paginate, nil
}

func (r *bookingRepository) FindEnding(bookingDate time.Time, page int, size int) ([]model.Booking, dto.Paginate, error) {
	var bookings []model.Booking

	query := "SELECT id, customer_id, court_id, booking_date, start_time, end_time, total_payment, status FROM bookings WHERE booking_date = $1 AND status = 'booked' LIMIT $2 OFFSET $3"

	offset := (page - 1) * size

	rows, err := r.DB.Query(query, bookingDate, size, offset)
	if err != nil {
		return []model.Booking{}, dto.Paginate{}, err
	}

	totalRows := 0

	for rows.Next() {
		var b model.Booking
		if err := rows.Scan(
			&b.Id,
			&b.Customer.Id,
			&b.Court.Id,
			&b.BookingDate,
			&b.StartTime,
			&b.EndTime,
			&b.Total_Payment,
			&b.Status,
		); err != nil {
			return []model.Booking{}, dto.Paginate{}, err
		}

		bookings = append(bookings, b)
		totalRows++
	}

	paginate := dto.Paginate{
		Page:       page,
		Size:       size,
		TotalRows:  totalRows,
		TotalPages: int(math.Ceil(float64(totalRows) / float64(size))),
	}

	return bookings, paginate, nil
}

func (r *bookingRepository) FindPaymentReport(day, month, year, page, size int, filterType string) ([]model.Payment, dto.Paginate, int64, error) {
	var payments []model.Payment
	var rows *sql.Rows
	var err error

	query := "SELECT id, booking_id, order_id, description, payment_method, price FROM payments WHERE "

	offset := (page - 1) * size

	if filterType == "daily" {
		dailyQuery := "EXTRACT(DAY FROM created_at) = $1  AND EXTRACT(MONTH FROM created_at) = $2 AND EXTRACT(YEAR FROM created_at) = $3 LIMIT $4 OFFSET $5"
		query += dailyQuery

		rows, err = r.DB.Query(query, day, month, year, size, offset)

	} else if filterType == "monthly" {
		monthlyQuery := "EXTRACT(MONTH FROM created_at) = $1 AND EXTRACT(YEAR FROM created_at) = $2 LIMIT $3 OFFSET $4"
		query += monthlyQuery

		rows, err = r.DB.Query(query, month, year, size, offset)

	} else if filterType == "yearly" {
		yearlyQuery := "EXTRACT(YEAR FROM created_at) = $1 LIMIT $2 OFFSET $3"
		query += yearlyQuery

		rows, err = r.DB.Query(query, year, size, offset)
	}

	if err != nil {
		return []model.Payment{}, dto.Paginate{}, 0, err
	}

	totalRows := 0
	var totalIncome int64

	for rows.Next() {
		var p model.Payment
		if err := rows.Scan(
			&p.Id,
			&p.BookingId,
			&p.OrderId,
			&p.Description,
			&p.PaymentMethod,
			&p.Price,
		); err != nil {
			return []model.Payment{}, dto.Paginate{}, 0, err
		}

		payments = append(payments, p)
		totalRows++
		totalIncome += int64(p.Price)
	}

	paginate := dto.Paginate{
		Page:       page,
		Size:       size,
		TotalRows:  totalRows,
		TotalPages: int(math.Ceil(float64(totalRows) / float64(size))),
	}

	return payments, paginate, totalIncome, nil
}

// func (r *bookingRepository) UpdateCancel(orderId string) error {
// 	transaction, _ := r.DB.Begin()

// 	bookingId := ""
// 	err := transaction.QueryRow("SELECT booking_id FROM payments WHERE order_id = $1", orderId).Scan(&bookingId)
// 	if err != nil {
// 		transaction.Rollback()
// 		return err
// 	}

// 	_, err = transaction.Exec("UPDATE bookings SET status = $1 WHERE id = $2", "cancel", bookingId)
// 	if err != nil {
// 		transaction.Rollback()
// 		return err
// 	}

// 	transaction.Commit()
// 	return nil
// }

func NewBookingRepository(db *sql.DB) BookingRepository {
	return &bookingRepository{
		DB: db,
	}
}
