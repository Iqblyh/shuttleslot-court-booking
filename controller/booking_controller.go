package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"team2/shuttleslot/model/dto"
	"team2/shuttleslot/service"
	"team2/shuttleslot/util"
	"time"

	"github.com/gin-gonic/gin"
)

type BookingController struct {
	service service.BookingService
	rg      *gin.RouterGroup
}

func (c *BookingController) Route() {
	router := c.rg.Group("bookings")
	{
		router.POST("/", c.CreateBookingHandler)
		router.POST("/payment/notif", c.NotificationHandler)
		router.POST("/repayment", c.CreateRepayHandler)
		router.GET("/", c.GetAllBookingsHandler)
		router.GET("/check", c.CheckBookingHandler)
		router.GET("/ending", c.CheckEndingHandler)
		router.GET("/report", c.PaymentReportHandler)
		// router.GET("/payment/cancel", c.GetCancel)
	}
}

func (c *BookingController) CreateBookingHandler(ctx *gin.Context) {
	var payload dto.CreateBookingRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		util.SendErrorResponse(ctx, err.Error(), http.StatusBadRequest)
		return
	}

	dateNow := util.StringToDate(time.Now().Format("02-01-2006"))
	timeNow := util.StringToTime(time.Now().Format("15:04:05"))
	bookingDate := util.StringToDate(payload.BookingDate)
	startTime := util.StringToTime(payload.StartTime)

	if bookingDate.Before(dateNow) {
		util.SendErrorResponse(ctx, "booking date cant in the past", http.StatusBadRequest)
		return
	}

	if bookingDate.Equal(dateNow) && startTime.Before(timeNow) {
		util.SendErrorResponse(ctx, "start time cant in the past", http.StatusBadRequest)
		return
	}

	data, err := c.service.Create(payload)
	if err != nil {
		if strings.Contains(err.Error(), "cannot book") {
			util.SendErrorResponse(ctx, err.Error(), http.StatusBadRequest)
			return
		}
		util.SendErrorResponse(ctx, err.Error(), http.StatusInternalServerError)
		return
	}

	response := util.CreateBookingResponse{}

	util.SendSingleResponse(ctx, "booking created successfully", response.FromModel(data), http.StatusCreated)
}

func (c *BookingController) NotificationHandler(ctx *gin.Context) {
	var payload dto.PaymentNotificationInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		util.SendErrorResponse(ctx, err.Error(), http.StatusBadRequest)
		return
	}

	err := c.service.UpdatePayment(payload)
	if err != nil {
		fmt.Println("=============== ERROR >>>>", err.Error())
		util.SendErrorResponse(ctx, err.Error(), http.StatusInternalServerError)
		return
	}

	util.SendSingleResponse(ctx, "booking created successfully", payload, http.StatusOK)
}

func (c *BookingController) CreateRepayHandler(ctx *gin.Context) {
	var payload dto.CreateRepayRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		util.SendErrorResponse(ctx, err.Error(), http.StatusBadRequest)
		return
	}

	data, err := c.service.CreateRepay(payload)
	if err != nil {
		util.SendErrorResponse(ctx, err.Error(), http.StatusInternalServerError)
		return
	}

	response := util.CreateRepaymentResponse{}
	util.SendSingleResponse(ctx, "repayment created successfully", response.FromModel(data), http.StatusCreated)
}

func (c *BookingController) GetAllBookingsHandler(ctx *gin.Context) {
	page, err1 := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, err2 := strconv.Atoi(ctx.DefaultQuery("size", "10"))
	if err1 != nil || err2 != nil {
		util.SendErrorResponse(ctx, "invalid page or size", http.StatusBadRequest)
		return
	}

	rows, paginate, err := c.service.FindAllBookings(page, size)
	if err != nil {
		util.SendErrorResponse(ctx, "Data not found", http.StatusNotFound)
		return
	}

	var listData []any
	var reponseTemplate util.GetBookingsResponse
	for _, v := range rows {
		listData = append(listData, reponseTemplate.FromModel(v))
	}

	util.SendPaginateResponse(ctx, "success get data", listData, paginate, http.StatusOK)
}

func (c *BookingController) CheckBookingHandler(ctx *gin.Context) {
	page, err1 := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, err2 := strconv.Atoi(ctx.DefaultQuery("size", "10"))
	if err1 != nil || err2 != nil {
		util.SendErrorResponse(ctx, "invalid page or size", http.StatusBadRequest)
		return
	}

	var bookingDate time.Time

	if ctx.Query("bookingDate") == "" {
		bookingDate = util.StringToDate(time.Now().String())
	} else {
		bookingDate = util.StringToDate(ctx.Query("bookingDate"))
	}

	rows, paginate, err := c.service.FindBookedCourt(bookingDate, page, size)
	if err != nil {
		util.SendErrorResponse(ctx, err.Error(), http.StatusInternalServerError)
		return
	}

	var listData []any
	var responseTemplate util.CheckBookingResponse

	for _, val := range rows {
		listData = append(listData, responseTemplate.FromModel(val))
	}

	util.SendPaginateResponse(ctx, "success get data", listData, paginate, http.StatusOK)
}

func (c *BookingController) CheckEndingHandler(ctx *gin.Context) {
	page, err1 := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, err2 := strconv.Atoi(ctx.DefaultQuery("size", "10"))
	if err1 != nil || err2 != nil {
		util.SendErrorResponse(ctx, "invalid page or size", http.StatusBadRequest)
		return
	}

	bookingDate := util.StringToDate(time.Now().Format("02-01-2006"))

	rows, paginate, err := c.service.FindEndingBookings(bookingDate, page, size)
	if err != nil {
		util.SendErrorResponse(ctx, err.Error(), http.StatusInternalServerError)
		return
	}

	var listData []any
	var responseTemplate util.GetEndingResponse

	for _, val := range rows {
		listData = append(listData, responseTemplate.FromModel(val))
	}

	util.SendPaginateResponse(ctx, "success get data", listData, paginate, http.StatusOK)
}

func (c *BookingController) PaymentReportHandler(ctx *gin.Context) {
	page, err1 := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, err2 := strconv.Atoi(ctx.DefaultQuery("size", "10"))
	if err1 != nil || err2 != nil {
		util.SendErrorResponse(ctx, "invalid page or size", http.StatusBadRequest)
		return
	}

	defaultDay := strconv.Itoa(time.Now().Day())
	defaultMonth := strconv.Itoa(int(time.Now().Month()))
	defaultYear := strconv.Itoa(time.Now().Year())

	filter := ctx.DefaultQuery("filter", "daily")
	day, _ := strconv.Atoi(ctx.DefaultQuery("day", defaultDay))
	month, _ := strconv.Atoi(ctx.DefaultQuery("month", defaultMonth))
	year, _ := strconv.Atoi(ctx.DefaultQuery("year", defaultYear))

	if !util.IsValidFilter(filter) {
		util.SendErrorResponse(ctx, "invalid filter, use 'daily', 'monthly', 'yearly'", http.StatusBadRequest)
		return
	}

	rows, paginate, totalIncome, err := c.service.FindPaymentReport(day, month, year, page, size, filter)

	if err != nil {
		util.SendErrorResponse(ctx, err.Error(), http.StatusInternalServerError)
		return
	}

	var listData []any
	var responseTemplate util.GetPaymentReportResponse

	for _, val := range rows {
		listData = append(listData, responseTemplate.FromModel(val))
	}

	util.SendReportPaginateResponse(ctx, "success get data", listData, totalIncome, paginate, http.StatusOK)

}

// func (c *BookingController) GetCancel(ctx *gin.Context) {
// 	orderId := ctx.Query("order_id")

// 	fmt.Println("================ Payload >>>> ", orderId)

// 	err := c.service.UpdateCancel(orderId)
// 	if err != nil {
// 		util.SendErrorResponse(ctx, "Data with that id not found", http.StatusInternalServerError)
// 		return
// 	}

// 	util.SendSingleResponse(ctx, "success update data", orderId, http.StatusOK)

// }

func NewBookingController(bookingService service.BookingService, rg *gin.RouterGroup) *BookingController {
	return &BookingController{
		service: bookingService,
		rg:      rg,
	}
}
