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

var conn *websocket.Conn

func Connect() {
	token, err := createToken("restaurateur")
	if err != nil {
		fmt.Println("error on generating token on user side")
		return
	}

	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/ws", RawQuery: "token=" + token}

	var dialErr error
	conn, _, dialErr = websocket.DefaultDialer.Dial(u.String(), nil)
	if dialErr != nil {
		log.Fatal("dial:", dialErr)
	}

}
func SendMessage(wsMessage WsMessage) {
	err := conn.WriteJSON(wsMessage)
	if err != nil {
		fmt.Println("error generated during sending message inside restauratuers")
		fmt.Println(err)
	}
}
func ReadMessage(messageChannel chan<- string, uuid string) {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}

		if strings.Contains(string(message), "restaurateur") &&
			strings.Contains(string(message), uuid) {
			messageChannel <- string(message)
			break
		}
	}
	fmt.Println("exited from reading messages")

}
