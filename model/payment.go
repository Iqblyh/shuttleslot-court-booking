package model

type Payment struct {
	Id            string `json:"id"`
	BookingId     string `json:"bookingId"`
	User          User
	Court         Court
	OrderId       string `json:"orderId"`
	Description   string `json:"description"`
	PaymentMethod string `json:"paymentMethod"`
	Price         int    `json:"price"`
	Qty           int    `json:"qty"`
	Status        string `json:"status"`
	PaymentURL    string `json:"paymentURL"`
}
