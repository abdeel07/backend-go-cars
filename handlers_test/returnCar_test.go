package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReturnCar(t *testing.T) {

	router := setupRouter()

	// Create a kilometers payload for returning the car.
	kilometersPayload := map[string]float64{
		"kilometers": 250.7,
	}

	// Convert the payload to JSON.
	payloadBytes, err := json.Marshal(kilometersPayload)
	assert.NoError(t, err)

	// Create a PUT request to return the car with registration "Reg1".
	request, err := http.NewRequest("PUT", "/cars/Reg1/returns", bytes.NewBuffer(payloadBytes))
	assert.NoError(t, err)

	// Record the response.
	response := httptest.NewRecorder()

	// Serve the HTTP request.
	router.ServeHTTP(response, request)

	fmt.Printf("\n------\n")
	fmt.Printf("Test Return Car - HTTP Status Code: %d (Must be 200)\n", response.Code)
	assert.Equal(t, http.StatusOK, response.Code)

	// Convert the response body into a CarResponse.
	var returnCarResponse CarResponse
	err = json.Unmarshal(response.Body.Bytes(), &returnCarResponse)

	fmt.Printf("Test Return Car - Car Availability: %t (Must be true)\n", returnCarResponse.Car.Available)
	assert.True(t, returnCarResponse.Car.Available == true)

	fmt.Printf("Test Return Car - Car Mileage: %f (Must be 750.7)\n", returnCarResponse.Car.Mileage)
	assert.True(t, returnCarResponse.Car.Mileage == 750.7)
}

func TestReturnCar2(t *testing.T) {

	router := setupRouter()

	// Create a kilometers payload for returning the car.
	kilometersPayload := map[string]float64{
		"kilometers": -500,
	}

	// Convert the payload to JSON.
	payloadBytes, err := json.Marshal(kilometersPayload)
	assert.NoError(t, err)

	// Create a PUT request to return the car with registration "Reg1".
	request, err := http.NewRequest("PUT", "/cars/Reg1/returns", bytes.NewBuffer(payloadBytes))
	assert.NoError(t, err)

	// Record the response.
	response := httptest.NewRecorder()

	fmt.Printf("\n")
	fmt.Printf("Test Return Car 2 - With negative kilometers\n")

	// Serve the HTTP request.
	router.ServeHTTP(response, request)

	fmt.Printf("Test Return Car 2 - HTTP Status Code: %d (Must be 400)\n", response.Code)
	assert.Equal(t, http.StatusBadRequest, response.Code)
}

func TestReturnCar3(t *testing.T) {

	router := setupRouter()

	// Create a kilometers payload for returning the car.
	kilometersPayload := map[string]float64{
		"kilometers": 250.7,
	}

	// Convert the payload to JSON.
	payloadBytes, err := json.Marshal(kilometersPayload)
	assert.NoError(t, err)

	// Create a PUT request to return a non-existing car with registration "RegXXX".
	request, err := http.NewRequest("PUT", "/cars/RegXXX/returns", bytes.NewBuffer(payloadBytes))
	assert.NoError(t, err)

	// Record the response.
	response := httptest.NewRecorder()

	fmt.Printf("\n")
	fmt.Printf("Test Return Car 3 - With Car not exist\n")

	// Serve the HTTP request.
	router.ServeHTTP(response, request)

	fmt.Printf("Test Return Car 3 - HTTP Status Code: %d (Must be 404)\n", response.Code)
	assert.Equal(t, http.StatusNotFound, response.Code)
}
