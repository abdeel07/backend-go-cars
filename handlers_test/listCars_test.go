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

	request, err := http.NewRequest("GET", "/cars", nil)
	assert.NoError(t, err)

	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	fmt.Printf("Test List Cars - HTTP Status Code: %d (Must be 200)\n", response.Code)
	assert.Equal(t, http.StatusOK, response.Code)

	var cars []model.Car
	err = json.Unmarshal(response.Body.Bytes(), &cars)
	assert.NoError(t, err)

	fmt.Printf("Test List Cars - Number of Cars in the Response: %d (Must be 3)\n", len(cars))
	assert.True(t, len(cars) == 3)
}
