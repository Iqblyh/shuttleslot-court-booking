package model

import "time"

type Booking struct {
	Id             string    `json:"id"`
	Customer       User      `json:"customer"`
	Court          Court     `json:"court"`
	Employee       User      `json:"employee"`
	BookingDate    time.Time `json:"bookingDate"`
	StartTime      time.Time `json:"startTime"`
	EndTime        time.Time `json:"endTime"`
	Total_Payment  int       `json:"totalPayment"`
	Status         string    `json:"status"`
	PaymentDetails []Payment `json:"paymentDetails"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}
