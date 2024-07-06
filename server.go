package main

import (
	"database/sql"
	"fmt"
	"team2/shuttleslot/config"
	"team2/shuttleslot/controller"
	"team2/shuttleslot/repository"
	"team2/shuttleslot/service"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Server struct {
	uS      service.UserService
	engine  *gin.Engine
	portApp string
}

func (s *Server) initiateRoute() {
	routerGroup := s.engine.Group("/api/v1")
	controller.NewUserController(s.uS, routerGroup).Route()
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

	userService := service.NewUserService(userRepository)

	return &Server{
		uS:      userService,
		engine:  gin.Default(),
		portApp: portApp,
	}
}
