package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/abdeel07/backend-go-cars/model"
	"github.com/abdeel07/backend-go-cars/server"
	"github.com/gorilla/mux"
)

// CarResponse represents the response structure for car-related operations.
type CarResponse struct {
	Message string    `json:"message"`
	Car     model.Car `json:"car"`
}

// ListCars handles the GET HTTP request to list all cars.
func ListCars(s *server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var cars []model.Car
		s.ParkingLotService.DB.Find(&cars)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(cars)
	}
}

// AddCar handles the POST HTTP request to add a new car to the system.
func AddCar(s *server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Decode the request body into a model.Car instance.
		var car model.Car
		if err := json.NewDecoder(r.Body).Decode(&car); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{"Invalid request payload"})
			return
		}

		// Validate the request parameters.
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

		// Check if the car already exists in the system.
		exists, _ := s.ParkingLotService.IsExist(car.Registration)

		if exists {
			// Return a not found response if the car is not found.
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(ErrorResponse{"Car already exists"})
			return

		}

		// Set car availability and create the car in the database.
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

// RentCar handles the PUT HTTP request to rent a car.
func RentCar(s *server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Extract registration parameter from the request.
		params := mux.Vars(r)
		registration := params["registration"]

		// Check if the specified car exists in the system.
		exists, car := s.ParkingLotService.IsExist(registration)

		if !exists {
			// Return a not found response if the car is not found.
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(ErrorResponse{"Car not found"})
			return
		}

		// Check if the car is available for rent.
		if !car.Available {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(ErrorResponse{"Car is not available"})
			return
		}

		// Set car availability to false and save the changes in the database.
		car.Available = false
		s.ParkingLotService.DB.Save(&car)

		// Create a response indicating the successful rental of the car.
		response := CarResponse{
			Message: "The Car with registration " + registration + " is rented!",
			Car:     car,
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

// ReturnCar handles the PUT HTTP request to return a rented car.
func ReturnCar(s *server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Extract registration parameter from the request.
		params := mux.Vars(r)
		registration := params["registration"]

		// Define a structure to hold the payload for kilometers returned.
		type KilometersPayload struct {
			Kilometers float64 `json:"kilometers"`
		}

		// Decode the request body into the KilometersPayload structure.
		var payload KilometersPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{"Invalid kilometers payload"})
			return

		}

		// Validate that the returned kilometers value is non-negative.
		if payload.Kilometers < 0 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{"Kilometers parameter must be positive"})
			return
		}

		// Check if the specified car exists in the system.
		exists, car := s.ParkingLotService.IsExist(registration)

		if !exists {
			// Return a not found response if the car is not found.
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(ErrorResponse{"Car not found"})
			return
		}

		// Check if the car is currently available (should not be available for return).
		if car.Available {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(ErrorResponse{"Car is already available"})
			return
		}

		// Update the car availability, increase the mileage, and save changes in the database.
		car.Available = true
		car.Mileage += float64(payload.Kilometers)
		s.ParkingLotService.DB.Save(&car)

		// Create a response indicating the successful return of the car.
		response := CarResponse{
			Message: "The Car with registration " + registration + " is returned!",
			Car:     car,
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

// GetCar handles the GET HTTP request to get details of a specific car.
func GetCar(s *server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Extract registration parameter from the request.
		params := mux.Vars(r)
		registration := params["registration"]

		// Check if the specified car exists in the system.
		exists, car := s.ParkingLotService.IsExist(registration)

		if !exists {
			// Return a not found response if the car is not found.
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(ErrorResponse{"Car not found"})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(car)
	}
}

// DeleteCar handles the DELETE HTTP request to delete a specific car from the system.
func DeleteCar(s *server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Extract registration parameter from the request.
		params := mux.Vars(r)
		registration := params["registration"]

		// Check if the specified car exists in the system.
		exists, car := s.ParkingLotService.IsExist(registration)

		if !exists {
			// Return a not found response if the car is not found.
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(ErrorResponse{"Car not found"})
			return
		}

		// Delete the car from the database.
		s.ParkingLotService.DB.Delete(&car)

		w.WriteHeader(http.StatusNoContent)
		json.NewEncoder(w).Encode("The Car with registration " + registration + " is deleted!")
	}
}

type ErrorResponse struct {
	Error string `json:"error"`
}
