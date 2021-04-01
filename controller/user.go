package controller

import (
	"CRUD-Table-Backend/config"
	"CRUD-Table-Backend/model"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user model.User
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &user)

	var res model.ResponseResult
	if err != nil {
		res.Error = "Error while processing input data, Try again"
		json.NewEncoder(w).Encode(res)
		return
	}

	collection := config.DB.Collection("User")

	if err != nil {
		res.Error = "Error while getting data from database, Try again"
		json.NewEncoder(w).Encode(res)
		return
	}

	// checking if the username and email have already existed in the database
	var result model.User
	errUsername := collection.FindOne(context.TODO(), bson.D{primitive.E{Key: "username", Value: user.Username}}).Decode(&result)

	if errUsername == nil {
		res.Error = "Username has been used. Please use different username"
		json.NewEncoder(w).Encode(res)
		return
	}

	errEmail := collection.FindOne(context.TODO(), bson.D{primitive.E{Key: "email", Value: user.Email}}).Decode(&result)

	if errEmail == nil {
		res.Error = "Email has been registered in the database"
		json.NewEncoder(w).Encode(res)
		return
	}

	// if both username and email have not existed in the database
	// then create the user.
	if errUsername.Error() == "mongo: no documents in result" && errUsername.Error() == errEmail.Error() {
		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 5)

		if err != nil {
			res.Error = "Error While Hashing Password"
			json.NewEncoder(w).Encode(res)
			return
		}

		user.Password = string(hash)

		_, err = collection.InsertOne(context.TODO(), user)

		if err != nil {
			res.Error = "Error While Creating User"
			json.NewEncoder(w).Encode(res)
			return
		}

		res.Result = "Registration Successful"
		json.NewEncoder(w).Encode(res)
		return
	}

	res.Error = err.Error()
	json.NewEncoder(w).Encode(res)
	return
}

func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user model.User
	var res model.ResponseResult

	body, _ := ioutil.ReadAll(r.Body)

	err := json.Unmarshal(body, &user)

	if err != nil {
		res.Error = "Error while processing input data, Try again"
		json.NewEncoder(w).Encode(res)
		return
	}

	// get collections from database
	collection := config.DB.Collection("User")

	var result model.User

	err = collection.FindOne(context.TODO(), bson.D{primitive.E{Key: "username", Value: user.Username}}).Decode(&result)

}
