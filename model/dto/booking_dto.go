package dto

type CreateBookingRequest struct {
	CourtId     string `json:"courtId"`
	BookingDate string `json:"bookingDate"`
	StartTime   string `json:"startTime"`
	Hour        int    `json:"hour"`
	CustomerId  string `json:"customerId"`
}

type CreateRepayRequest struct {
	BookingId     string `json:"bookingId"`
	EmployeeId    string `json:"employeeId"`
	PaymentMethod string `json:"paymentMethod"`
}

type PaymentNotificationInput struct {
	TransactionStatus string `json:"transaction_status"`
	OrderId           string `json:"order_id"`
	PaymentType       string `json:"payment_type"`
	FraudStatus       string `json:"fraud_status"`
}

func (c CreateRepayRequest) IsValidMethod() bool {
	return c.PaymentMethod == "mid" || c.PaymentMethod == "cash"
}