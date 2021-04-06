package controller

import (
	"CRUD-Table-Backend/config"
	"CRUD-Table-Backend/helper"
	"CRUD-Table-Backend/model"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"

	jwt "github.com/dgrijalva/jwt-go"

	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	var user model.User
	body, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(body, &user)

	res := model.ResponseResult{}
	if err != nil {
		res.Error = "Error while processing input data, Try again"
		c.JSON(http.StatusBadRequest, res)
		return
	}

	collection := config.DB.Collection("User")

	if err != nil {
		res.Error = "Error while getting data from database, Try again"
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	// checking if the username and email have already existed in the database
	errUsername := helper.CheckIfUsernameExistInDatabase(user.Username)
	errEmail := helper.CheckIfEmailExistInDatabase(user.Email)

	if errUsername == nil {
		res.Error = "Username has been used. Please use different username"
		c.JSON(http.StatusBadRequest, res)
		return
	}

	if errEmail == nil {
		res.Error = "Email has been registered in the database"
		c.JSON(http.StatusBadRequest, res)
		return
	}

	// if both username and email have not existed in the database
	// then create the user.
	if errUsername.Error() == "mongo: no documents in result" && errUsername.Error() == errEmail.Error() {
		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 5)

		if err != nil {
			res.Error = "Error While Hashing Password"
			c.JSON(http.StatusInternalServerError, res)
			return
		}

		user.Password = string(hash)

		_, err = collection.InsertOne(context.TODO(), user)

		if err != nil {
			res.Error = "Error While Creating User"
			c.JSON(http.StatusInternalServerError, res)
			return
		}

		res.Result = "Registration Successful"
		c.JSON(http.StatusOK, res)
		return
	}
}

func Login(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	var user model.User
	var res model.ResponseResult

	body, _ := ioutil.ReadAll(c.Request.Body)

	err := json.Unmarshal(body, &user)

	if err != nil {
		res.Error = "Error while processing input data, Try again"
		c.JSON(http.StatusBadRequest, res)
		return
	}

	result, errLogin := helper.IsLoginInformationCorrect(user.Username, user.Password)

	if errLogin != nil {
		res.Error = errLogin.Error()
		c.JSON(http.StatusBadRequest, res)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": result.Username,
		"email":    result.Email,
	})

	tokenString, err := token.SignedString([]byte("secret"))

	if err != nil {
		res.Error = err.Error()
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	result.Token = tokenString
	result.Password = ""

	res.Error = ""
	res.Result = result

	c.JSON(http.StatusOK, res)
}
