package main

import (
	"CRUD-Table-Backend/config"
	"CRUD-Table-Backend/controller"
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// initialize server before starting one.
func init() {
	config.ConnectToDB()
}

// Main function
func main() {
	// r := mux.NewRouter()

	// r.HandleFunc("/register", controller.Register).Methods("POST")
	// r.HandleFunc("/login", controller.Login).Methods("POST")

	r := gin.Default()

	r.POST("/register", controller.Register)
	r.POST("/login", controller.Login)

	// disconnect after the program closes.
	defer config.Client.Disconnect(context.TODO())

	log.Fatal(http.ListenAndServe(":8080", r))
}
