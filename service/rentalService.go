package service

import (
	"sync"

	"github.com/abdeel07/backend-go-cars/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ParkingLotService struct {
	DB        *gorm.DB
	CarsMutex *sync.Mutex
}

func (s *ParkingLotService) IsExist(registration string) (bool, model.Car) {
	var car model.Car
	result := s.DB.First(&car, "registration = ?", registration)

	if result.Error != nil {
		return false, model.Car{}
	}

	return true, car
}

func NewParkingLotService(db *gorm.DB) *ParkingLotService {
	return &ParkingLotService{
		DB:        db,
		CarsMutex: &sync.Mutex{},
	}
}

func MigrateDB(db *gorm.DB) {
	db.AutoMigrate(&model.Car{})
}

func InitializeDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	return db, nil
}
