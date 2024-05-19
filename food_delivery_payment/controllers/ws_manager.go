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
	token, err := createToken("payment")
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
		if strings.Contains(string(message), "payment") &&
			strings.Contains(string(message), "check_payment") {
			CheckPayment(string(message))

		}

	}

}

type Message struct {
	ReceiverName string `json:"reciever_name"`
	MessageID    string `json:"message_id"`
	Response     string `json:"response"`
}

func SendResponseMessage(response Message) {

	err := Conn.WriteJSON(response)
	if err != nil {
		fmt.Println("error generated while pushing message in side of payment")
		fmt.Println(err)
	}
}
