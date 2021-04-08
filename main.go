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
	r.Use(Cors())

	r.POST("/register", controller.Register)
	r.POST("/login", controller.Login)

	// disconnect after the program closes.
	defer config.Client.Disconnect(context.TODO())

	log.Fatal(http.ListenAndServe(":8080", r))
}

// Cors => allow access to non origin
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, token")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
