package handlers_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRentCar(t *testing.T) {

	router := setupRouter()

	// Create a PUT request to rent a car with registration "Reg1".
	request, err := http.NewRequest("PUT", "/cars/Reg1/rentals", nil)
	assert.NoError(t, err)

	// Record the response.
	response := httptest.NewRecorder()

	// Serve the HTTP request.
	router.ServeHTTP(response, request)

	fmt.Printf("\n------\n")
	fmt.Printf("Test Rent Car - HTTP Status Code: %d (Must be 200)\n", response.Code)
	assert.Equal(t, http.StatusOK, response.Code)

	// Convert the response body into a CarResponse.
	var rentCarResponse CarResponse
	err = json.Unmarshal(response.Body.Bytes(), &rentCarResponse)
	if err != nil {
		fmt.Printf("Error unmarshalling response body: %s\n", err)
	}

	fmt.Printf("Test Rent Car - Car Availability: %t (Must be false)\n", rentCarResponse.Car.Available)
	assert.True(t, rentCarResponse.Car.Available == false)
}

func TestRentCar2(t *testing.T) {

	router := setupRouter()

	// Create a PUT request to rent a car with registration "Reg1" (already rented).
	request, err := http.NewRequest("PUT", "/cars/Reg1/rentals", nil)
	assert.NoError(t, err)

	// Record the response.
	response := httptest.NewRecorder()

	fmt.Printf("\n")
	fmt.Printf("Test Rent Car 2 - With Car already rented\n")

	// Serve the HTTP request.
	router.ServeHTTP(response, request)

	fmt.Printf("Test Rent Car 2 - HTTP Status Code: %d (Must be 409)\n", response.Code)
	assert.Equal(t, http.StatusConflict, response.Code)
}

func TestRentCar3(t *testing.T) {

	router := setupRouter()

	// Create a PUT request to rent a non-existing car with registration "RegXXX".
	request, err := http.NewRequest("PUT", "/cars/RegXXX/rentals", nil)
	assert.NoError(t, err)

	// Record the response.
	response := httptest.NewRecorder()

	fmt.Printf("\n")
	fmt.Printf("Test Rent Car 3 - With Car not exist\n")

	// Serve the HTTP request.
	router.ServeHTTP(response, request)

	fmt.Printf("Test Rent Car 3 - HTTP Status Code: %d (Must be 404)\n", response.Code)
	assert.Equal(t, http.StatusNotFound, response.Code)
}
