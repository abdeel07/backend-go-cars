package routes

import (
	"github.com/abdeel07/backend-go-cars/handlers"
	"github.com/abdeel07/backend-go-cars/server"
	"github.com/gorilla/mux"
)

func SetupRoutes(router *mux.Router, s *server.Server) {
	router.HandleFunc("/cars", handlers.ListCars(s)).Methods("GET")
	router.HandleFunc("/cars", handlers.AddCar(s)).Methods("POST")
	router.HandleFunc("/cars/{registration}", handlers.GetCar(s)).Methods("GET")
	router.HandleFunc("/cars/{registration}", handlers.DeleteCar(s)).Methods("DELETE")
	router.HandleFunc("/cars/{registration}/rentals", handlers.RentCar(s)).Methods("PUT")
	router.HandleFunc("/cars/{registration}/returns", handlers.ReturnCar(s)).Methods("PUT")
}
