package controller

import (
	"net/http"
	"strconv"
	"strings"
	"team2/shuttleslot/model/dto"
	"team2/shuttleslot/service"
	"team2/shuttleslot/util"

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
		// router.GET("/payment/cancel", c.GetCancel)
	}
}

func (c *BookingController) CreateBookingHandler(ctx *gin.Context) {
	var payload dto.CreateBookingRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		util.SendErrorResponse(ctx, err.Error(), http.StatusBadRequest)
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
