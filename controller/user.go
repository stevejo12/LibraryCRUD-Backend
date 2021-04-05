package controller

import (
	"CRUD-Table-Backend/config"
	"CRUD-Table-Backend/helper"
	"CRUD-Table-Backend/model"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"

	"golang.org/x/crypto/bcrypt"
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
	errUsername := helper.CheckIfUsernameExistInDatabase(user.Username)
	errEmail := helper.CheckIfEmailExistInDatabase(user.Email)

	if errUsername == nil {
		res.Error = "Username has been used. Please use different username"
		json.NewEncoder(w).Encode(res)
		return
	}

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

	result, errLogin := helper.IsLoginInformationCorrect(user.Username, user.Password)

	if errLogin != nil {
		res.Error = errLogin.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": result.Username,
		"email":    result.Email,
	})

	tokenString, err := token.SignedString([]byte("secret"))

	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	result.Token = tokenString
	result.Password = ""

	json.NewEncoder(w).Encode(result)
}
