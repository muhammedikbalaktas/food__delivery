package controllers

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/websocket"
)

const password = "SystemTokenPassword"

var MessageChan = make(chan string)

func createToken(role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"role": role,
	})

	tokenString, err := token.SignedString([]byte(password))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

var Conn *websocket.Conn

func Connect() {
	token, err := createToken("restaurant")
	if err != nil {
		fmt.Println("error on generating token on user side")
		return
	}

	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/ws", RawQuery: "token=" + token}

	var dialErr error
	Conn, _, dialErr = websocket.DefaultDialer.Dial(u.String(), nil)
	if dialErr != nil {
		log.Fatal("dial:", dialErr)
	}
	ReadMessage()

}

func ReadMessage() {
	for {
		_, message, err := Conn.ReadMessage()

		if err != nil {
			log.Println("read:", err)
			return
		}
		if strings.Contains(string(message), "restaurant") &&
			strings.Contains(string(message), "create_restaurant") {

			CreateRestaurant(string(message))
		} else if strings.Contains(string(message), "restaurant") &&
			strings.Contains(string(message), "list_restaurants") {

			ListRestaurants(string(message))
		} else if strings.Contains(string(message), "restaurant") &&
			strings.Contains(string(message), "list_meals") {

			ListMeals(string(message))
		} else if strings.Contains(string(message), "restaurant") &&
			strings.Contains(string(message), "add_meal") {

			CheckAvailability(string(message))
		} else if strings.Contains(string(message), "restaurant") &&
			strings.Contains(string(message), "order_meal") {

			GetOrder(string(message))
		}

	}

}
