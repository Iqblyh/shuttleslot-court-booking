package repomock

import (
	"team2/shuttleslot/model"
	"team2/shuttleslot/model/dto"
	"time"

	"github.com/stretchr/testify/mock"
)

type BookingRepositoryMock struct {
	mock.Mock
}

func (b *BookingRepositoryMock) Create(payload model.Booking) (model.Booking, error) {
	args := b.Called(payload)
	return args.Get(0).(model.Booking), args.Error(1)
}
func (b *BookingRepositoryMock) FindAll(page int, size int) ([]model.Booking, dto.Paginate, error) {
	args := b.Called(page, size)
	return args.Get(0).([]model.Booking), args.Get(1).(dto.Paginate), args.Error(2)
}
func (b *BookingRepositoryMock) FindByDate(bookingDate time.Time) ([]model.Booking, error) {
	args := b.Called(bookingDate)
	return args.Get(0).([]model.Booking), args.Error(1)
}
func (b *BookingRepositoryMock) FindById(bookingId string) (model.Booking, error) {
	args := b.Called(bookingId)
	return args.Get(0).(model.Booking), args.Error(1)
}
func (b *BookingRepositoryMock) FindTotal(customerId string) (int, error) {
	args := b.Called(customerId)
	return args.Int(0), args.Error(1)
}
func (b *BookingRepositoryMock) FindPaymentByOrderId(order_id string) (model.Payment, error) {
	args := b.Called(order_id)
	return args.Get(0).(model.Payment), args.Error(1)
}
func (b *BookingRepositoryMock) UpdateStatus(payload model.Payment) error {
	args := b.Called(payload)
	return args.Error(0)
}
func (b *BookingRepositoryMock) CreateRepay(payload model.Payment) (model.Payment, error) {
	args := b.Called(payload)
	return args.Get(0).(model.Payment), args.Error(1)
}
func (b *BookingRepositoryMock) UpdateRepaymentStatus(payload model.Payment) error {
	args := b.Called(payload)
	return args.Error(0)
}
func (b *BookingRepositoryMock) FindBooked(bookingDate time.Time, page int, size int) ([]model.Booking, dto.Paginate, error) {
	args := b.Called(bookingDate, page, size)
	return args.Get(0).([]model.Booking), args.Get(1).(dto.Paginate), args.Error(2)

}
func (b *BookingRepositoryMock) FindEnding(bookingDate time.Time, page int, size int) ([]model.Booking, dto.Paginate, error) {
	args := b.Called(bookingDate, page, size)
	return args.Get(0).([]model.Booking), args.Get(1).(dto.Paginate), args.Error(2)
}

func (b *BookingRepositoryMock) FindPaymentReport(day, month, year, page, size int, filterType string) ([]model.Payment, dto.Paginate, int64, error) {
	args := b.Called(day, month, year, page, size, filterType)
	return args.Get(0).([]model.Payment), args.Get(1).(dto.Paginate), args.Get(2).(int64), args.Error(3)
}
