package util

import (
	"net/http"
	"team2/shuttleslot/model/dto"

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

func SendErrorResponse(c *gin.Context, message string, code int) {
	c.JSON(code, dto.SingleResponse{
		Status: dto.Status{
			Code:    code,
			Message: message,
		},
	})
}
