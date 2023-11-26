package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/abdeel07/backend-go-cars/model"
	"github.com/stretchr/testify/assert"
)

func TestAddCar(t *testing.T) {

	router := setupRouter()

	car := model.Car{
		CarModel:     "New Model",
		Registration: "New Registration",
		Mileage:      100,
	}

	carJSON, err := json.Marshal(car)
	assert.NoError(t, err)

	request, err := http.NewRequest("POST", "/cars", bytes.NewBuffer(carJSON))
	assert.NoError(t, err)
	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	fmt.Printf("\n------\n")
	fmt.Printf("Test Add Car - HTTP Status Code: %d (Must be 201)\n", response.Code)
	assert.Equal(t, http.StatusCreated, response.Code)

	var addedCar model.Car
	err = json.Unmarshal(response.Body.Bytes(), &addedCar)
	assert.NoError(t, err)

	fmt.Printf("Test Add Car - Car Availability: %t (Must be true)\n", addedCar.Available)
	assert.True(t, addedCar.Available == true)
}

func TestAddCar2(t *testing.T) {

	router := setupRouter()

	car := model.Car{
		CarModel:     "New Model",
		Registration: "New Registration",
		Mileage:      100,
	}

	carJSON, err := json.Marshal(car)
	assert.NoError(t, err)

	request, err := http.NewRequest("POST", "/cars", bytes.NewBuffer(carJSON))
	assert.NoError(t, err)
	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()

	fmt.Printf("\n")
	fmt.Printf("Test Add Car 2 - With deplicate Registration\n")

	router.ServeHTTP(response, request)

	fmt.Printf("Test Add Car 2 - HTTP Status Code: %d (Must be 409)\n", response.Code)
	assert.Equal(t, http.StatusConflict, response.Code)
}
