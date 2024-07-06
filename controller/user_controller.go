package controller

import (
	"net/http"
	"strconv"
	"team2/shuttleslot/model"
	"team2/shuttleslot/service"
	"team2/shuttleslot/util"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService service.UserService
	rg          *gin.RouterGroup
}

func (c *UserController) CreateAdminHandler(ctx *gin.Context) {
	payload := model.User{}
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		util.SendErrorResponse(ctx, err.Error(), http.StatusBadRequest)
		return
	}

	data, err := c.userService.CreateAdmin(payload)
	if err != nil {
		util.SendErrorResponse(ctx, err.Error(), http.StatusInternalServerError)
		return
	}

	util.SendSingleResponse(ctx, "Admin created successfully", data, http.StatusOK)
}

func (c *UserController) CreateCustomerHandler(ctx *gin.Context) {
	payload := model.User{}
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		util.SendErrorResponse(ctx, err.Error(), http.StatusBadRequest)
		return
	}

	data, err := c.userService.CreateCustomer(payload)
	if err != nil {
		util.SendErrorResponse(ctx, err.Error(), http.StatusInternalServerError)
		return
	}

	util.SendSingleResponse(ctx, "customer created successfully", data, http.StatusOK)
}

func (c *UserController) CreateEmployeeHandler(ctx *gin.Context) {
	payload := model.User{}
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		util.SendErrorResponse(ctx, err.Error(), http.StatusBadRequest)
		return
	}

	data, err := c.userService.CreateEmployee(payload)
	if err != nil {
		util.SendErrorResponse(ctx, err.Error(), http.StatusInternalServerError)
		return
	}

	util.SendSingleResponse(ctx, "employee created successfully", data, http.StatusOK)
}

func (c *UserController) FindUserByRoleHandler(ctx *gin.Context) {
	role := ctx.Param("role")
	page, err1 := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, err2 := strconv.Atoi(ctx.DefaultQuery("size", "10"))
	if err1 != nil || err2 != nil {
		util.SendErrorResponse(ctx, "invalid page or size", http.StatusBadRequest)
		return
	}

	data, paginate, err := c.userService.FindUserByRole(role, page, size)

	if err != nil {
		util.SendErrorResponse(ctx, "role "+role+" is invalid", http.StatusInternalServerError)
		return
	}

	var listData []any
	for _, v := range data {
		listData = append(listData, v)
	}

	util.SendPaginateResponse(ctx, "success get data by role "+role, listData, paginate, http.StatusOK)
}

func (c *UserController) FindUserByUsernameHandler(ctx *gin.Context) {
	username := ctx.Param("username")
	data, err := c.userService.FindUserByUsername(username)

	if err != nil {
		util.SendErrorResponse(ctx, "user with username "+username+" not found", http.StatusNotFound)
		return
	}

	util.SendSingleResponse(ctx, "success", data, http.StatusOK)
}

func (c *UserController) FindUserByIdHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	data, err := c.userService.FindUserById(id)

	if err != nil {
		util.SendErrorResponse(ctx, "user with id "+id+" not found", http.StatusNotFound)
		return
	}

	util.SendSingleResponse(ctx, "success", data, http.StatusOK)
}

func (c *UserController) UpdateUserHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	var payload model.User
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		util.SendErrorResponse(ctx, "error in payload", http.StatusBadRequest)
		return
	}

	data, err := c.userService.UpdatedUser(id, payload)
	if err != nil {
		util.SendErrorResponse(ctx, err.Error(), http.StatusInternalServerError)
		return
	}

	util.SendSingleResponse(ctx, "user updated successfully", data, http.StatusOK)
}

func (c *UserController) DeleteUserHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	err := c.userService.DeletedUser(id)
	if err != nil {
		util.SendErrorResponse(ctx, err.Error(), http.StatusInternalServerError)
		return
	}

	util.SendSingleResponse(ctx, "user deleted successfully", nil, http.StatusOK)
}

func (c *UserController) Route() {
	router := c.rg.Group("users")
	router.POST("/register", c.CreateCustomerHandler)
	router.POST("/admin/create", c.CreateAdminHandler)
	router.POST("/employee/create", c.CreateEmployeeHandler)
	router.GET("/role/:role", c.FindUserByRoleHandler)
	router.GET("/username/:username", c.FindUserByUsernameHandler)
	router.GET("/:id", c.FindUserByIdHandler)
	router.PUT("/:id", c.UpdateUserHandler)
	router.DELETE("/:id", c.DeleteUserHandler)
}

func NewUserController(userService service.UserService, rg *gin.RouterGroup) *UserController {
	return &UserController{
		userService: userService,
		rg:          rg,
	}
}
