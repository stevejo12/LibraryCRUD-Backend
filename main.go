package main

import (
	"CRUD-Table-Backend/config"
	"CRUD-Table-Backend/controller"
	"context"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// initialize server before starting one.
func init() {
	config.ConnectToDB()
}

// Main function
func main() {
	r := mux.NewRouter()

	r.HandleFunc("/register", controller.Register).Methods("POST")
	r.HandleFunc("/login", controller.Login).Methods("POST")

	// disconnect after the program closes.
	defer config.Client.Disconnect(context.TODO())

	log.Fatal(http.ListenAndServe(":8080", r))
}
