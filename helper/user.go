package helper

import (
	"CRUD-Table-Backend/config"
	"CRUD-Table-Backend/model"
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CheckIfUsernameExistInDatabase(username string) error {
	collection := config.DB.Collection("User")

	var result model.User
	errUsername := collection.FindOne(context.TODO(), bson.D{primitive.E{Key: "username", Value: username}}).Decode(&result)

	return errUsername
}

func CheckIfEmailExistInDatabase(email string) error {
	collection := config.DB.Collection("User")

	var result model.User
	errEmail := collection.FindOne(context.TODO(), bson.D{primitive.E{Key: "email", Value: email}}).Decode(&result)

	return errEmail
}

func IsLoginInformationCorrect(username string, password string) (model.User, error) {
	collection := config.DB.Collection("User")

	var result model.User
	errLogin := collection.FindOne(context.TODO(), bson.D{primitive.E{Key: "username", Value: username}}).Decode(&result)

	if errLogin != nil {
		return model.User{}, errors.New("Invalid Username")
	}

	errLogin = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(password))

	if errLogin != nil {
		return model.User{}, errors.New("Invalid Password")
	}

	return result, nil
}
