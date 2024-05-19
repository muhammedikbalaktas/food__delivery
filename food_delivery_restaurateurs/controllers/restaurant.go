package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type WsMessage struct {
	ReceiverName   string `json:"reciever_name"`
	MessageID      string `json:"message_id"`
	ActionName     string `json:"action_name"`
	RestaurantName string `json:"restaurant_name"`
	UserToken      string `json:"user_token"`
}
type InputParam struct {
	RestaurantName string `json:"restaurant_name"`
}

func CreateRestaurant(c *gin.Context) {
	token := c.GetHeader("Authorization")
	fmt.Println(token)
	_, err := parseToken(token)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid token for restaurateur"})
		return
	}
	var inputParam InputParam
	err = c.BindJSON(&inputParam)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid params for restaurateur"})
		return
	}

	newUUID, err := uuid.NewRandom()
	if err != nil {
		fmt.Println("error generated for generating new uuid in restaurateur")
		fmt.Println(err)
	}
	var wsMessage = WsMessage{
		ReceiverName:   "restaurant",
		MessageID:      newUUID.String(),
		ActionName:     "create_restaurant",
		RestaurantName: inputParam.RestaurantName,
		UserToken:      token,
	}

	SendMessage(wsMessage)
	messageChannel := make(chan string)
	go ReadMessage(messageChannel, newUUID.String())

	message := <-messageChannel
	if strings.Contains(message, "error") {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "error happened when creating restaurant"})
		return
	} else if strings.Contains(message, "success") {
		c.IndentedJSON(http.StatusOK, gin.H{"success": "restaurant created succesfully"})
		return
	}

}
