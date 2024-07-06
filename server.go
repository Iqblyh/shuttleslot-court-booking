package main

import (
	"database/sql"
	"fmt"
	"team2/shuttleslot/config"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Server struct {
	engine  *gin.Engine
	portApp string
}

func (s *Server) initiateRoute() {
	// routerGroup := s.engine.Group("/api/v1")
}

func (s *Server) Start() {
	s.initiateRoute()
	fmt.Println(s.portApp)
	s.engine.Run(s.portApp)
}

func NewServer() *Server {
	co, _ := config.NewConfig()

	urlConnection := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", co.Host, co.Port, co.User, co.Password, co.Name)

	db, err := sql.Open(co.Driver, urlConnection)

	fmt.Print(db) //Temp

	if err != nil {
		return &Server{}
	}

	portApp := co.AppPort
	
	return &Server{
		engine:  gin.Default(),
		portApp: portApp,
	}
}
