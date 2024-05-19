package controllers

import (
	"fmt"
	m "food-delivery/users/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func CreateUser(c *gin.Context) {
	var user m.User

	err := c.BindJSON(&user)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "cannot create user. Invalid json"})
		return
	}
	validator := validator.New()
	err = validator.Struct(user)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "cannot create user. Invalid parameters"})
		return
	}
	db, err := createDb()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "error opening database"})
		return
	}
	defer db.Close()
	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "error hashing password"})
		return
	}
	query := "insert into users(first_name, last_name, email, password) values(?,?,?,?)"
	_, err = db.Exec(query, user.FirstName, user.LastName, user.Email, hashedPassword)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "error adding new user"})
		return
	}
	c.IndentedJSON(http.StatusOK, map[string]string{"succes": "user created succesfully"})

}

type params struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func GetUser(c *gin.Context) {
	var parameters params
	err := c.BindJSON(&parameters)
	if err != nil {
		fmt.Println(err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "cannot create user. Invalid json"})
		return
	}

	validator := validator.New()
	err = validator.Struct(parameters)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "cannot create user. Invalid parameters"})
		return
	}

	db, err := createDb()
	if err != nil {
		fmt.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "error opening database"})
		return
	}
	defer db.Close()
	query := "select id, password from users where email=?"

	rows, err := db.Query(query, parameters.Email)
	if err != nil {
		fmt.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "error getting user"})
		return
	}
	var (
		password string
		userID   int
	)
	if rows.Next() {

		err := rows.Scan(&userID, &password)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "error parsing row 1"})
			return
		}

	} else {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "user does not exist"})
		return
	}
	err = comparePasswords([]byte(password), []byte(parameters.Password))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid password for user"})
		return
	}
	token, err := generateToken(userID)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "could not get token"})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"token": token})

}
