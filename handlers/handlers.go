package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/abdeel07/backend-go-cars/model"
	"github.com/abdeel07/backend-go-cars/server"
	"github.com/gorilla/mux"
)

type CarResponse struct {
	Message string    `json:"message"`
	Car     model.Car `json:"car"`
}

func ListCars(s *server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var cars []model.Car
		s.ParkingLotService.DB.Find(&cars)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(cars)
	}
}

func AddCar(s *server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var car model.Car
		if err := json.NewDecoder(r.Body).Decode(&car); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{"Invalid request payload"})
			return
		}

		if car.Mileage < 0 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{"Mileage parameter must be positive"})
			return

		}
		if car.CarModel == "" || car.Registration == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{"Car model and registration are required"})
			return
		}

		exists, _ := s.ParkingLotService.IsExist(car.Registration)

		if exists {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(ErrorResponse{"Car already exists"})
			return

		}

		car.Available = true
		if err := s.ParkingLotService.DB.Create(&car).Error; err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{"Failed to create car"})
			return
		}

		w.WriteHeader(http.StatusCreated)
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
			json.NewEncoder(w).Encode(ErrorResponse{"Car not found"})
			return
		}

		if !car.Available {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(ErrorResponse{"Car is not available"})
			return
		}

		car.Available = false
		s.ParkingLotService.DB.Save(&car)

		response := CarResponse{
			Message: "The Car with registration " + registration + " is rented!",
			Car:     car,
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

func ReturnCar(s *server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r)
		registration := params["registration"]

		type KilometersPayload struct {
			Kilometers float64 `json:"kilometers"`
		}

		var payload KilometersPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{"Invalid kilometers payload"})
			return

		}
		if payload.Kilometers < 0 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{"Kilometers parameter must be positive"})
			return
		}

		exists, car := s.ParkingLotService.IsExist(registration)

		if !exists {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(ErrorResponse{"Car not found"})
			return
		}

		if car.Available {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(ErrorResponse{"Car is already available"})
			return
		}

		car.Available = true
		car.Mileage += float64(payload.Kilometers)
		s.ParkingLotService.DB.Save(&car)

		response := CarResponse{
			Message: "The Car with registration " + registration + " is returned!",
			Car:     car,
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
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
			json.NewEncoder(w).Encode(ErrorResponse{"Car not found"})
			return
		}

		w.WriteHeader(http.StatusOK)
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
			json.NewEncoder(w).Encode(ErrorResponse{"Car not found"})
			return
		}

		s.ParkingLotService.DB.Delete(&car)

		w.WriteHeader(http.StatusNoContent)
		json.NewEncoder(w).Encode("The Car with registration " + registration + " is deleted!")
	}
}

type ErrorResponse struct {
	Error string `json:"error"`
}
