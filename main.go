package main

import (
	"fmt"
	"net/http"

	"github.com/abdeel07/backend-go-cars/routes"
	"github.com/abdeel07/backend-go-cars/server"
	"github.com/abdeel07/backend-go-cars/service"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	db, err := service.InitializeDB("root:@tcp(127.0.0.1:3306)/parking_lot?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println("Error initializing database:", err)
	}

	service.MigrateDB(db)

	parkingLotServer := server.NewServer(db)

	routes.SetupRoutes(router, parkingLotServer)

	http.ListenAndServe(":8080", router)
	fmt.Println("Server started on port 8080")
}
