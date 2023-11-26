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

	kilometersPayload := map[string]float64{
		"kilometers": 250.7,
	}

	payloadBytes, err := json.Marshal(kilometersPayload)
	assert.NoError(t, err)

	request, err := http.NewRequest("PUT", "/cars/Reg1/returns", bytes.NewBuffer(payloadBytes))
	assert.NoError(t, err)

	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	fmt.Printf("Test Return Car - HTTP Status Code: %d (Must be 200)\n", response.Code)
	assert.Equal(t, http.StatusOK, response.Code)

	var returnCarResponse CarResponse
	err = json.Unmarshal(response.Body.Bytes(), &returnCarResponse)

	fmt.Printf("Test Return Car - Car Availability: %t (Must be true)\n", returnCarResponse.Car.Available)
	assert.True(t, returnCarResponse.Car.Available == true)

	fmt.Printf("Test Return Car - Car Mileage: %f (Must be 750.7)\n", returnCarResponse.Car.Mileage)
	assert.True(t, returnCarResponse.Car.Mileage == 750.7)
}

func TestReturnCar2(t *testing.T) {

	router := setupRouter()

	kilometersPayload := map[string]float64{
		"kilometers": -500,
	}

	payloadBytes, err := json.Marshal(kilometersPayload)
	assert.NoError(t, err)

	request, err := http.NewRequest("PUT", "/cars/Reg1/returns", bytes.NewBuffer(payloadBytes))
	assert.NoError(t, err)

	response := httptest.NewRecorder()

	fmt.Printf("Test Return Car 2 - With negative kilometers\n")
	router.ServeHTTP(response, request)

	fmt.Printf("Test Return Car 2 - HTTP Status Code: %d (Must be 400)\n", response.Code)
	assert.Equal(t, http.StatusBadRequest, response.Code)
}

func TestReturnCar3(t *testing.T) {

	router := setupRouter()

	kilometersPayload := map[string]float64{
		"kilometers": 250.7,
	}

	payloadBytes, err := json.Marshal(kilometersPayload)
	assert.NoError(t, err)

	request, err := http.NewRequest("PUT", "/cars/RegXXX/returns", bytes.NewBuffer(payloadBytes))
	assert.NoError(t, err)

	response := httptest.NewRecorder()

	fmt.Printf("Test Return Car 3 - With Car not exist\n")
	router.ServeHTTP(response, request)

	fmt.Printf("Test Return Car 3 - HTTP Status Code: %d (Must be 404)\n", response.Code)
	assert.Equal(t, http.StatusNotFound, response.Code)
}
