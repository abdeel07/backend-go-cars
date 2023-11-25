package server

import (
	"github.com/abdeel07/backend-go-cars/service"
	"gorm.io/gorm"
)

type Server struct {
	ParkingLotService *service.ParkingLotService
}

func NewServer(db *gorm.DB) *Server {
	return &Server{
		ParkingLotService: service.NewParkingLotService(db),
	}
}
