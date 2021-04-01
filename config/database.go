package config

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DB => store connection ???
var DB *mongo.Database

// Client that access database
var Client *mongo.Client

// ConnectToDB => Establish connection to the database
func ConnectToDB() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	Client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	err = Client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected To MongoDB!")

	// point to specific database
	DB = Client.Database("crudtable")
}
