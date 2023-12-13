package handlers_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/abdeel07/backend-go-cars/model"
	"github.com/abdeel07/backend-go-cars/routes"
	"github.com/abdeel07/backend-go-cars/server"
	"github.com/abdeel07/backend-go-cars/service"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func TestMain(m *testing.M) {

	// Initialize the test database.
	db, err := service.InitializeDB("root:@tcp(127.0.0.1:3306)/parking_lot_test?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println("Error initializing database:", err)
		os.Exit(1)
	}

	service.MigrateDB(db)

	seedTestData(db)

	exitCode := m.Run()

	clearTestData(db)

	sqlDB, _ := db.DB()
	sqlDB.Close()

	os.Exit(exitCode)
}

// seedTestData inserts test data into the database.
func seedTestData(db *gorm.DB) {
	cars := []model.Car{
		{CarModel: "Model1", Registration: "Reg1", Mileage: 500, Available: true},
		{CarModel: "Model2", Registration: "Reg2", Mileage: 160, Available: true},
		{CarModel: "Model3", Registration: "Reg3", Mileage: 1000, Available: false},
	}

	for _, car := range cars {
		db.Create(&car)
	}
}

// clearTestData deletes all test data from the database.
func clearTestData(db *gorm.DB) {
	db.Exec("DELETE FROM cars")
}

// CarResponse represents the response structure for car-related operations.
type CarResponse struct {
	Message string    `json:"message"`
	Car     model.Car `json:"car"`
}

func setupRouter() *mux.Router {
	router := mux.NewRouter()

	db, err := service.InitializeDB("root:@tcp(127.0.0.1:3306)/parking_lot_test?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println("Error initializing database:", err)
	}

	service.MigrateDB(db)

	parkingLotServer := server.NewServer(db)

	routes.SetupRoutes(router, parkingLotServer)

	return router
}
