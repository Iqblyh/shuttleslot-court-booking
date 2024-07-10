package service

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"team2/shuttleslot/model"
	"team2/shuttleslot/model/dto"
	"team2/shuttleslot/repository"
	"team2/shuttleslot/util"
	"time"
)

type BookingService interface {
	Create(payload dto.CreateBookingRequest) (model.Booking, error)
	UpdatePayment(payload dto.PaymentNotificationInput) error
	CreateRepay(payload dto.CreateRepayRequest) (model.Payment, error)
	FindAllBookings(page int, size int) ([]model.Booking, dto.Paginate, error)
	FindBookedCourt(bookingDate time.Time, page int, size int) ([]model.Booking, dto.Paginate, error)
	FindEndingBookings(bookingDate time.Time, page int, size int) ([]model.Booking, dto.Paginate, error)
	FindPaymentReport(day, month, year, page, size int, filterType string) ([]model.Payment, dto.Paginate, int64, error)
	// UpdateCancel(orderId string) error //Update Cancel Masih ku coba, skip dulu unit testingnya
}

type bookingService struct {
	bookingRepository repository.BookingRepository
	userServ          UserService
	courtServ         CourtService
	payGate           PaymentGateService
}

func (s *bookingService) Create(payload dto.CreateBookingRequest) (model.Booking, error) {
	var newPayload model.Booking

	existBooking, err := s.bookingRepository.FindByDate(util.StringToDate(payload.BookingDate))
	if err != nil {
		return model.Booking{}, err
	}

	endTime := util.StringToTime(payload.StartTime).Add(time.Hour * time.Duration(payload.Hour))

	for _, val := range existBooking {
		if val.Customer.Id == payload.CustomerId && val.Status == "pending" {
			return model.Booking{}, errors.New("cannot book, there still payment to complete")
		}
		if val.Court.Id == payload.CourtId && util.DateToString(val.BookingDate) == payload.BookingDate && (val.Status == "pending" || val.Status == "booked" || val.Status == "done") {
			if util.InTimeSpanStart(val.StartTime, val.EndTime, util.StringToTime(payload.StartTime)) {
				err = errors.New("cannot book court in that time")

			} else if util.InTimeSpanEnd(val.StartTime, val.EndTime, endTime) {
				err = errors.New("cannot book that long, because the schedule collides with another schedule")

			}
		}
	}
	if err != nil {
		return model.Booking{}, err
	}

	customer, err := s.userServ.FindUserById(payload.CustomerId)
	if err != nil {
		return model.Booking{}, err
	}

	court, err := s.courtServ.FindCourtById(payload.CourtId)
	if err != nil {
		return model.Booking{}, err
	}

	totalBooking, err := s.bookingRepository.FindTotal(payload.CustomerId)
	if err != nil {
		return model.Booking{}, err
	}

	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	orderId := fmt.Sprintf("Booking%s-%d", fmt.Sprintf("%05d", totalBooking+1), random.Int())
	desc := fmt.Sprintf("Pembayaran Booking %s", court.Name)

	realCourtPrice := court.Price

	court.Price = court.Price / 2

	payment := model.Payment{
		OrderId:     orderId,
		Description: desc,
		Court:       court,
		User:        customer,
		Price:       (realCourtPrice * payload.Hour) / 2,
		Qty:         payload.Hour,
	}

	paymentURL, err := s.payGate.GetPaymentURL(payment)
	if err != nil {
		return model.Booking{}, err
	}

	newPayload = model.Booking{
		Customer:      customer,
		Court:         court,
		Total_Payment: realCourtPrice * payload.Hour,
		BookingDate:   util.StringToDate(payload.BookingDate),
		StartTime:     util.StringToTime(payload.StartTime),
		EndTime:       endTime,
		PaymentDetails: []model.Payment{
			{
				OrderId:     payment.OrderId,
				Description: payment.Description,
				PaymentURL:  paymentURL,
			},
		},
	}

	booking, err := s.bookingRepository.Create(newPayload)
	if err != nil {
		return model.Booking{}, err
	}

	booking.PaymentDetails[0].PaymentURL = paymentURL
	booking.Court = court
	booking.Customer = customer

	return booking, nil
}

func (s *bookingService) UpdatePayment(payload dto.PaymentNotificationInput) error {
	payment, err := s.payGate.PaymentProcess(payload)
	if err != nil {
		return err
	}

	if strings.Contains(payload.OrderId, "Booking") {
		err = s.bookingRepository.UpdateStatus(payment)
		if err != nil {
			return err
		}

		return nil
	}

	fmt.Println("==================== PAYMENT >>>> ", payment)

	err = s.bookingRepository.UpdateRepaymentStatus(payment)
	if err != nil {
		fmt.Println("==================== ERROR >>>> ", err.Error())

		return err
	}

	return nil
}

func (s *bookingService) CreateRepay(payload dto.CreateRepayRequest) (model.Payment, error) {
	var newPayload model.Payment

	booking, err := s.bookingRepository.FindById(payload.BookingId)
	if err != nil {
		return model.Payment{}, err
	}

	if booking.Status != "booked" {
		return model.Payment{}, errors.New("this booking still not booked")
	}

	customer, err := s.userServ.FindUserById(booking.Customer.Id)
	if err != nil {
		return model.Payment{}, err
	}

	court, err := s.courtServ.FindCourtById(booking.Court.Id)
	if err != nil {
		return model.Payment{}, err
	}

	totalBooking, err := s.bookingRepository.FindTotal(booking.Customer.Id)
	if err != nil {
		return model.Payment{}, err
	}

	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	orderId := fmt.Sprintf("Repayment%s-%d", fmt.Sprintf("%05d", totalBooking), random.Int())
	desc := fmt.Sprintf("Pelunasan Booking %s", court.Name)

	realCourtPrice := court.Price
	court.Price = court.Price / 2

	newPayload = model.Payment{
		BookingId:     payload.BookingId,
		OrderId:       orderId,
		Description:   desc,
		PaymentMethod: payload.PaymentMethod,
		Price:         booking.Total_Payment / 2,
		Qty:           booking.Total_Payment / realCourtPrice,
		Court:         court,
	}

	if payload.PaymentMethod != "mid" {
		newPayload.User.Id = payload.EmployeeId
		payment, err := s.bookingRepository.CreateRepay(newPayload)
		if err != nil {
			return model.Payment{}, err
		}

		return payment, nil
	}

	newPayload.User = customer

	paymentUrl, err := s.payGate.GetPaymentURL(newPayload)
	if err != nil {
		return model.Payment{}, err
	}

	newPayload.User.Id = payload.EmployeeId
	newPayload.PaymentURL = paymentUrl

	payment, err := s.bookingRepository.CreateRepay(newPayload)
	if err != nil {
		return model.Payment{}, err
	}

	return payment, nil
}

func (s *bookingService) FindAllBookings(page int, size int) ([]model.Booking, dto.Paginate, error) {

	return s.bookingRepository.FindAll(page, size)
}

func (s *bookingService) FindBookedCourt(bookingDate time.Time, page int, size int) ([]model.Booking, dto.Paginate, error) {
	bookings, paginate, err := s.bookingRepository.FindBooked(bookingDate, page, size)
	if err != nil {
		return []model.Booking{}, dto.Paginate{}, err
	}

	for i, val := range bookings {
		court, err := s.courtServ.FindCourtById(val.Court.Id)
		if err != nil {
			return []model.Booking{}, dto.Paginate{}, err
		}

		bookings[i].Court = court
	}

	return bookings, paginate, nil
}

func (s *bookingService) FindEndingBookings(bookingDate time.Time, page int, size int) ([]model.Booking, dto.Paginate, error) {
	bookings, paginate, err := s.bookingRepository.FindEnding(bookingDate, page, size)
	if err != nil {
		return []model.Booking{}, dto.Paginate{}, err
	}

	for i, val := range bookings {
		customer, err := s.userServ.FindUserById(val.Customer.Id)
		if err != nil {
			return []model.Booking{}, dto.Paginate{}, err
		}

		court, err := s.courtServ.FindCourtById(val.Court.Id)
		if err != nil {
			return []model.Booking{}, dto.Paginate{}, err
		}

		bookings[i].Customer = customer
		bookings[i].Court = court
	}

	return bookings, paginate, nil
}

func (s *bookingService) FindPaymentReport(day, month, year, page, size int, filterType string) ([]model.Payment, dto.Paginate, int64, error) {
	return s.bookingRepository.FindPaymentReport(day, month, year, page, size, filterType)
}

// func (s *bookingService) UpdateCancel(orderId string) error {
// 	return s.bookingRepository.UpdateCancel(orderId)
// }

func NewBookingService(bookingRepository repository.BookingRepository, userService UserService, courtService CourtService, payGate PaymentGateService) BookingService {
	return &bookingService{
		bookingRepository: bookingRepository,
		userServ:          userService,
		courtServ:         courtService,
		payGate:           payGate,
	}
}
