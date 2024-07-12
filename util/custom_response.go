package util

import (
	"net/http"
	"team2/shuttleslot/model"
	"team2/shuttleslot/model/dto"
	"time"

	"github.com/gin-gonic/gin"
)

func SendSingleResponse(c *gin.Context, message string, data any, code int) {
	c.JSON(http.StatusOK, dto.SingleResponse{
		Status: dto.Status{
			Code:    code,
			Message: message,
		},
		Data: data,
	})
}

func SendPaginateResponse(c *gin.Context, message string, data []any, paginate dto.Paginate, code int) {
	c.JSON(http.StatusOK, dto.PaginateResponse{
		Status: dto.Status{
			Code:    code,
			Message: message,
		},
		Data:     data,
		Paginate: paginate,
	})
}

func SendReportPaginateResponse(c *gin.Context, message string, data []any, totalIncome int64, paginate dto.Paginate, code int) {
	c.JSON(http.StatusOK, dto.ReportPaginateResponse{
		Status: dto.Status{
			Code:    code,
			Message: message,
		},
		Data:        data,
		TotalIncome: totalIncome,
		Paginate:    paginate,
	})
}

func SendErrorResponse(c *gin.Context, message string, code int) {
	c.JSON(code, dto.SingleResponse{
		Status: dto.Status{
			Code:    code,
			Message: message,
		},
	})
}

func SendPaymentResponse(c *gin.Context, data dto.PaymentResponse, code int) {
	c.JSON(code, dto.PaymentResponse{
		OrderId:           data.OrderId,
		TransactionStatus: data.TransactionStatus,
		Status: dto.Status{
			Code:    code,
			Message: data.TransactionStatus,
		},
	})
}

type GetUserByRoleResponse struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phoneNumber"`
	Email       string `json:"email"`
	Username    string `json:"username"`
	Point       int    `json:"point"`
	Role        string `json:"role"`
}

func (*GetUserByRoleResponse) FromModel(payload model.User) *GetUserByRoleResponse {
	return &GetUserByRoleResponse{
		Id:          payload.Id,
		Name:        payload.Name,
		PhoneNumber: payload.PhoneNumber,
		Email:       payload.Email,
		Username:    payload.Username,
		Point:       payload.Point,
		Role:        payload.Role,
	}
}

type CreateBookingResponse struct {
	BookingId    string          `json:"bookingId"`
	BookingDate  string          `json:"bookingDate"`
	CustomerName string          `json:"customerName"`
	CourtName    string          `json:"courtName"`
	StartTime    string          `json:"startTime"`
	EndTime      string          `json:"endTime"`
	TotalPayment int             `json:"totalPayment"`
	Payment      PaymentResponse `json:"payment"`
}

type PaymentResponse struct {
	OrderId     string `json:"orderId"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	PaymentUrl  string `json:"paymentUrl"`
}

func (*CreateBookingResponse) FromModel(payload model.Booking) *CreateBookingResponse {
	return &CreateBookingResponse{
		BookingId:    payload.Id,
		BookingDate:  DateToString(payload.BookingDate),
		CustomerName: payload.Customer.Name,
		CourtName:    payload.Court.Name,
		StartTime:    TimeToString(payload.StartTime),
		EndTime:      TimeToString(payload.EndTime),
		TotalPayment: payload.Total_Payment,
		Payment: PaymentResponse{
			OrderId:     payload.PaymentDetails[0].OrderId,
			Description: payload.PaymentDetails[0].Description,
			Price:       payload.PaymentDetails[0].Price,
			PaymentUrl:  payload.PaymentDetails[0].PaymentURL,
		},
	}
}

type CreateRepaymentResponse struct {
	BookingId     string `json:"bookingId"`
	OrderId       string `json:"orderId"`
	Description   string `json:"description"`
	Price         int    `json:"price"`
	PaymentMethod string `json:"paymentMethod"`
	PaymentUrl    string `json:"paymentUrl"`
}

func (*CreateRepaymentResponse) FromModel(payload model.Payment) *CreateRepaymentResponse {
	return &CreateRepaymentResponse{
		BookingId:     payload.BookingId,
		OrderId:       payload.OrderId,
		Description:   payload.Description,
		Price:         payload.Price,
		PaymentMethod: payload.PaymentMethod,
		PaymentUrl:    payload.PaymentURL,
	}
}

type GetBookingsResponse struct {
	Id          string      `json:"id"`
	Customer    UserBooking `json:"customer"`
	Court       CourBooking `json:"court"`
	Employee    UserBooking `json:"employee"`
	BookingDate string      `json:"bookingDate"`
	StartTime   string      `json:"startTime"`
	EndTime     string      `json:"endTime"`
	Status      string      `json:"status"`
	CreatedAt   time.Time   `json:"createdAt"`
	UpdatedAt   time.Time   `json:"updatedAt"`
}

type UserBooking struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phoneNumber"`
	Email       string `json:"email"`
}

type CourBooking struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

func (*GetBookingsResponse) FromModel(payload model.Booking) *GetBookingsResponse {
	return &GetBookingsResponse{
		Id: payload.Id,
		Customer: UserBooking{
			Id:          payload.Customer.Id,
			Name:        payload.Customer.Name,
			PhoneNumber: payload.Customer.PhoneNumber,
			Email:       payload.Customer.Email,
		},
		Court: CourBooking{
			Id:    payload.Court.Id,
			Name:  payload.Court.Name,
			Price: payload.Court.Price,
		},
		Employee: UserBooking{
			Id:          payload.Employee.Id,
			Name:        payload.Employee.Name,
			PhoneNumber: payload.Employee.PhoneNumber,
			Email:       payload.Employee.Email,
		},
		BookingDate: DateToString(payload.BookingDate),
		StartTime:   TimeToString(payload.StartTime),
		EndTime:     TimeToString(payload.EndTime),
		Status:      payload.Status,
		CreatedAt:   payload.CreatedAt,
		UpdatedAt:   payload.UpdatedAt,
	}
}

type CheckBookingResponse struct {
	CourtName   string `json:"courtName"`
	BookingDate string `json:"bookingDate"`
	StartTime   string `json:"startTime"`
	EndTime     string `json:"endTime"`
	Status      string `json:"status"`
}

func (*CheckBookingResponse) FromModel(payload model.Booking) *CheckBookingResponse {
	return &CheckBookingResponse{
		CourtName:   payload.Court.Name,
		BookingDate: DateToString(payload.BookingDate),
		StartTime:   TimeToString(payload.StartTime),
		EndTime:     TimeToString(payload.EndTime),
		Status:      payload.Status,
	}
}

type GetEndingResponse struct {
	BookingId    string `json:"bookingId"`
	CustomerName string `json:"customerName"`
	CourtName    string `json:"courtName"`
	BookingDate  string `json:"bookingDate"`
	StartTime    string `json:"startTime"`
	EndTime      string `json:"endTime"`
	TotalPayment int    `json:"totalPayment"`
	Status       string `json:"status"`
}

func (*GetEndingResponse) FromModel(payload model.Booking) *GetEndingResponse {
	return &GetEndingResponse{
		BookingId:    payload.Id,
		CustomerName: payload.Customer.Name,
		CourtName:    payload.Court.Name,
		BookingDate:  DateToString(payload.BookingDate),
		StartTime:    TimeToString(payload.StartTime),
		EndTime:      TimeToString(payload.EndTime),
		TotalPayment: payload.Total_Payment,
		Status:       payload.Status,
	}
}

type GetPaymentReportResponse struct {
	PaymentId     string `json:"paymentId"`
	BookingId     string `json:"bookingId"`
	OrderId       string `json:"orderId"`
	Description   string `json:"description"`
	PaymentMethod string `json:"paymentMethod"`
	Price         int    `json:"price"`
}

func (*GetPaymentReportResponse) FromModel(payload model.Payment) *GetPaymentReportResponse {
	return &GetPaymentReportResponse{
		PaymentId:     payload.Id,
		BookingId:     payload.BookingId,
		OrderId:       payload.OrderId,
		Description:   payload.Description,
		PaymentMethod: payload.PaymentMethod,
		Price:         payload.Price,
	}
}
