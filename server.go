package main

import (
	"database/sql"
	"fmt"
	"team2/shuttleslot/config"
	"team2/shuttleslot/controller"
	"team2/shuttleslot/middleware"
	"team2/shuttleslot/repository"
	"team2/shuttleslot/service"
	"team2/shuttleslot/util"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Server struct {
	uS      service.UserService
	cS      service.CourtService
	bS      service.BookingService
	pGS     service.PaymentGateService
	auth    middleware.AuthMiddleware
	util    util.UtilInterface
	engine  *gin.Engine
	portApp string
}

func (s *Server) initiateRoute() {
	routerGroup := s.engine.Group("/api/v1")
	controller.NewUserController(s.uS, s.auth, routerGroup).Route()
	controller.NewCourtController(s.cS, s.auth, routerGroup).Route()
	controller.NewBookingController(s.bS, s.auth, routerGroup).Route()
}

func (s *Server) Start() {
	s.initiateRoute()
	s.engine.Run(s.portApp)
}

func NewServer() *Server {
	co, _ := config.NewConfig()

	urlConnection := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", co.Host, co.Port, co.User, co.Password, co.Name)

	db, err := sql.Open(co.Driver, urlConnection)
	if err != nil {
		return &Server{}
	}

	portApp := co.AppPort
	userRepository := repository.NewUserRepository(db)
	courtRepository := repository.NewCourtRepository(db)
	bookingRepository := repository.NewBookingRepository(db)

	utilService := util.NewUtilService()
	payGateService := service.NewPayGateService(co.PayGateConfig, bookingRepository)
	authService := service.NewAuthService(co.SecurityConfig)
	userService := service.NewUserService(userRepository, authService, utilService)
	courtService := service.NewCourtService(courtRepository)
	bookingService := service.NewBookingService(bookingRepository, userService, courtService, payGateService)

	authMiddleware := middleware.NewAuthMiddleware(authService)

	return &Server{
		uS:      userService,
		cS:      courtService,
		engine:  gin.Default(),
		bS:      bookingService,
		pGS:     payGateService,
		auth:    authMiddleware,
		portApp: portApp,
	}
}
