package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/abdeel07/backend-go-cars/model"
	"github.com/abdeel07/backend-go-cars/server"
	"github.com/gorilla/mux"
)

func ListCars(s *server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var cars []model.Car
		s.ParkingLotService.DB.Find(&cars)

		json.NewEncoder(w).Encode(cars)
	}
}

func AddCar(s *server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var car model.Car
		if err := json.NewDecoder(r.Body).Decode(&car); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request payload"})
			return
		}

		exists, _ := s.ParkingLotService.IsExist(car.Registration)

		if exists {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(map[string]string{"error": "Car already exists"})
			return
		}

		car.Available = true
		if err := s.ParkingLotService.DB.Create(&car).Error; err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create car"})
			return
		}

		json.NewEncoder(w).Encode(car)
	}
}

func RentCar(s *server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r)
		registration := params["registration"]

		exists, car := s.ParkingLotService.IsExist(registration)

		if !exists {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "Car not found"})
			return
		}

		if !car.Available {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(map[string]string{"error": "Car is not available"})
			return
		}

		car.Available = false
		s.ParkingLotService.DB.Save(&car)

		json.NewEncoder(w).Encode("The Car with registration " + registration + " is rented!")
	}
}

func ReturnCar(s *server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r)
		registration := params["registration"]

		var kilometersParam string
		json.NewDecoder(r.Body).Decode(&kilometersParam)

		exists, car := s.ParkingLotService.IsExist(registration)

		if !exists {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "Car not found"})
			return
		}

		if car.Available {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(map[string]string{"error": "Car is already available"})
			return
		}

		kilometers, err := strconv.Atoi(kilometersParam)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid kilometers parameter"})
			return
		}

		car.Available = true
		car.Mileage += int(kilometers)
		s.ParkingLotService.DB.Save(&car)

		json.NewEncoder(w).Encode(car)
	}
}

func GetCar(s *server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r)
		registration := params["registration"]

		exists, car := s.ParkingLotService.IsExist(registration)

		if !exists {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "Car not found"})
			return
		}

		json.NewEncoder(w).Encode(car)
	}
}

func DeleteCar(s *server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r)
		registration := params["registration"]

		exists, car := s.ParkingLotService.IsExist(registration)

		if !exists {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "Car not found"})
			return
		}

		s.ParkingLotService.DB.Delete(&car)

		json.NewEncoder(w).Encode("The Car with registration " + registration + " is deleted!")
	}
}
