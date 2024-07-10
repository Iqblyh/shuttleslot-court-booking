package service

import (
	"team2/shuttleslot/config"
	"team2/shuttleslot/model"
	"team2/shuttleslot/model/dto"
	"team2/shuttleslot/repository"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

type paymentGateService struct {
	config      config.PayGateConfig
	bookingRepo repository.BookingRepository
}

var s snap.Client

type PaymentGateService interface {
	GetPaymentURL(payment model.Payment) (string, error)
	PaymentProcess(payload dto.PaymentNotificationInput) (model.Payment, error)
}

func (p *paymentGateService) GetPaymentURL(payment model.Payment) (string, error) {
	s.New(p.config.ServerKey, midtrans.Sandbox)

	snapReq := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  payment.OrderId,
			GrossAmt: int64(payment.Price),
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: payment.User.Name,
			Email: "test@mail.com",
			Phone: payment.User.PhoneNumber,
		},
		Items: &[]midtrans.ItemDetails{
			{
				Name:  payment.Court.Name,
				Price: int64(payment.Court.Price),
				Qty:   int32(payment.Qty),
			},
		},
	}

	resp, err := s.CreateTransaction(snapReq)
	if err != nil {
		return "", err
	}

	return resp.RedirectURL, nil
}

func (p *paymentGateService) PaymentProcess(payload dto.PaymentNotificationInput) (model.Payment, error) {
	newPayload, err := p.bookingRepo.FindPaymentByOrderId(payload.OrderId)
	if err != nil {
		return model.Payment{}, err
	}

	booking, err := p.bookingRepo.FindById(newPayload.BookingId)
	if err != nil {
		return model.Payment{}, err
	}

	payment, err := p.bookingRepo.FindPaymentByOrderId(payload.OrderId)
	if err != nil {
		return model.Payment{}, err
	}

	if booking.Employee.Id != "" {
		payment.User.Id = booking.Employee.Id
	}

	if payload.PaymentType == "credit-card" && payload.TransactionStatus == "capture" && payload.FraudStatus == "accept" {
		payment.Status = "paid"

	} else if payload.TransactionStatus == "settlement" {
		payment.Status = "paid"

	} else if payload.TransactionStatus == "cancel" || payload.TransactionStatus == "expire" {
		payment.Status = "cancel"

	} else if payload.TransactionStatus == "pending" {
		payment.Status = "pending"
	}

	payment.PaymentMethod = payload.PaymentType

	return payment, nil
}

func NewPayGateService(payGateConfig config.PayGateConfig, bookingRepository repository.BookingRepository) PaymentGateService {
	return &paymentGateService{
		config:      payGateConfig,
		bookingRepo: bookingRepository,
	}
}
