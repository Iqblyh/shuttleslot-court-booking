package servicemock

import (
	"team2/shuttleslot/model"
	"team2/shuttleslot/model/dto"
	"time"

	"github.com/stretchr/testify/mock"
)

type BookingServiceMock struct {
	mock.Mock
}

func (b *BookingServiceMock) Create(payload dto.CreateBookingRequest) (model.Booking, error) {
	args := b.Called(payload)
	return args.Get(0).(model.Booking), args.Error(1)
}
func (b *BookingServiceMock) UpdatePayment(payload dto.PaymentNotificationInput) error {
	args := b.Called(payload)
	return args.Error(0)
}
func (b *BookingServiceMock) CreateRepay(payload dto.CreateRepayRequest) (model.Payment, error) {
	args := b.Called(payload)
	return args.Get(0).(model.Payment), args.Error(1)
}
func (b *BookingServiceMock) FindAllBookings(page int, size int) ([]model.Booking, dto.Paginate, error) {
	args := b.Called(page, size)
	return args.Get(0).([]model.Booking), args.Get(1).(dto.Paginate), args.Error(2)
}
func (b *BookingServiceMock) FindBookedCourt(bookingDate time.Time, page int, size int) ([]model.Booking, dto.Paginate, error) {
	args := b.Called(bookingDate, page, size)
	return args.Get(0).([]model.Booking), args.Get(1).(dto.Paginate), args.Error(2)
}
func (b *BookingServiceMock) FindEndingBookings(bookingDate time.Time, page int, size int) ([]model.Booking, dto.Paginate, error) {
	args := b.Called(bookingDate, page, size)
	return args.Get(0).([]model.Booking), args.Get(1).(dto.Paginate), args.Error(2)
}
func (b *BookingServiceMock) FindPaymentReport(day, month, year, page, size int, filterType string) ([]model.Payment, dto.Paginate, int64, error) {
	args := b.Called(day, month, year, page, size, filterType)
	return args.Get(0).([]model.Payment), args.Get(1).(dto.Paginate), args.Get(2).(int64), args.Error(3)
}
