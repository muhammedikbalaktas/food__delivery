package controllers

import (
	"encoding/json"
	"fmt"
	m "food-delivery/users/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type WsListRestaurantResponse struct {
	Restaurants []m.Restaurant `json:"restaurants"`
}
type WsMessage struct {
	ReceiverName string `json:"reciever_name"`
	MessageID    string `json:"message_id"`
	ActionName   string `json:"action_name"`
}

func ListRestaurants(c *gin.Context) {
	newUUID, err := uuid.NewRandom()
	if err != nil {
		fmt.Println("error generated for generating new uuid in restaurateur")
		fmt.Println(err)
	}
	var wsMessage = WsMessage{
		ReceiverName: "restaurant",
		MessageID:    newUUID.String(),
		ActionName:   "list_restaurants",
	}
	SendMessage(wsMessage)
	messageChannel := make(chan string)
	go ReadMessage(messageChannel, newUUID.String())

	message := <-messageChannel
	if strings.Contains(message, "error") {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "error happened when getting restaurants"})
		return
	} else if strings.Contains(message, "success") {
		var response WsListRestaurantResponse
		err := json.Unmarshal([]byte(message), &response)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "error on getting restaurants"})
			return
		}
		c.IndentedJSON(http.StatusOK, response)
		return
	}
}
