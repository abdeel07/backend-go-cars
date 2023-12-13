package handlers_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/abdeel07/backend-go-cars/model"
	"github.com/stretchr/testify/assert"
)

func TestListCars(t *testing.T) {

	router := setupRouter()

	// Create a GET request for the "/cars" endpoint.
	request, err := http.NewRequest("GET", "/cars", nil)
	assert.NoError(t, err)

	// Record the response.
	response := httptest.NewRecorder()

	// Serve the HTTP request.
	router.ServeHTTP(response, request)

	fmt.Printf("\n------\n")
	fmt.Printf("Test List Cars - HTTP Status Code: %d (Must be 200)\n", response.Code)
	// Assert the HTTP status code.
	assert.Equal(t, http.StatusOK, response.Code)

	// Convert the response body into a slice of model.Car.
	var cars []model.Car
	err = json.Unmarshal(response.Body.Bytes(), &cars)
	assert.NoError(t, err)

	fmt.Printf("Test List Cars - Number of Cars in the Response: %d (Must be >= 3)\n", len(cars))
	assert.True(t, len(cars) >= 3)
}
