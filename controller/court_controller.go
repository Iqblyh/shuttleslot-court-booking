package controller

import (
	"net/http"
	"strconv"
	"team2/shuttleslot/model"
	"team2/shuttleslot/service"
	"team2/shuttleslot/util"

	"github.com/gin-gonic/gin"
)

type CourtController struct {
	courtService service.CourtService
	rg           *gin.RouterGroup
}

func (c *CourtController) CreateCourtHandler(ctx *gin.Context) {
	var payload model.Court
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		util.SendErrorResponse(ctx, err.Error(), http.StatusBadRequest)
		return
	}

	data, err := c.courtService.CreateCourt(payload)
	if err != nil {
		util.SendErrorResponse(ctx, err.Error(), http.StatusInternalServerError)
		return
	}

	util.SendSingleResponse(ctx, "court created successfully", data, http.StatusCreated)
}

func (c *CourtController) FindAllCourtsHandler(ctx *gin.Context) {
	page, err1 := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, err2 := strconv.Atoi(ctx.DefaultQuery("size", "10"))
	if err1 != nil || err2 != nil {
		util.SendErrorResponse(ctx, "invalid page or size", http.StatusBadRequest)
		return
	}

	rows, paginate, err := c.courtService.FindAllCourts(page, size)
	if err != nil {
		util.SendErrorResponse(ctx, "Data not found", http.StatusNotFound)
		return
	}

	var listData []any
	for _, v := range rows {
		listData = append(listData, v)
	}

	util.SendPaginateResponse(ctx, "success get data", listData, paginate, http.StatusOK)
}

func (c *CourtController) FindCourtByIdHandler(ctx *gin.Context) {
	id := ctx.Param("id")

	data, err := c.courtService.FindCourtById(id)
	if err != nil {
		util.SendErrorResponse(ctx, err.Error(), http.StatusNotFound)
	}

	util.SendSingleResponse(ctx, "success get data", data, http.StatusOK)
}

func (c *CourtController) UpdateCourtHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	var payload model.Court

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		util.SendErrorResponse(ctx, err.Error(), http.StatusBadRequest)
		return
	}

	courtUpdate, err := c.courtService.UpdateCourt(id, payload)
	if err != nil {
		util.SendErrorResponse(ctx, err.Error(), http.StatusInternalServerError)
		return
	}

	util.SendSingleResponse(ctx, "court updated successfully", courtUpdate, http.StatusOK)
}

func (c *CourtController) DeleteCourtHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	err := c.courtService.DeleteCourt(id)
	if err != nil {
		util.SendErrorResponse(ctx, err.Error(), http.StatusNotFound)
		return
	}
	util.SendSingleResponse(ctx, "court deleted successfully", nil, http.StatusOK)
}

func (c *CourtController) Route() {
	router := c.rg.Group("courts")
	router.POST("/", c.CreateCourtHandler)
	router.GET("/", c.FindAllCourtsHandler)
	router.GET("/:id", c.FindCourtByIdHandler)
	router.PUT("/:id", c.UpdateCourtHandler)
	router.DELETE("/:id", c.DeleteCourtHandler)
}

func NewCourtController(courtService service.CourtService, rg *gin.RouterGroup) *CourtController {
	return &CourtController{
		courtService: courtService,
		rg:           rg,
	}
}
