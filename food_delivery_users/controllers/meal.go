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

type ListMealsResponse struct {
	Restaurants []m.Meal `json:"meals"`
}

func ListMeals(c *gin.Context) {
	newUUID, err := uuid.NewRandom()
	if err != nil {
		fmt.Println("error generated for generating new uuid in restaurateur")
		fmt.Println(err)
	}
	var wsMessage = WsMessage{
		ReceiverName: "restaurant",
		MessageID:    newUUID.String(),
		ActionName:   "list_meals",
	}
	SendMessage(wsMessage)
	messageChannel := make(chan string)
	go ReadMessage(messageChannel, newUUID.String())

	message := <-messageChannel
	if strings.Contains(message, "error") {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "error happened when getting meals"})
		return
	} else if strings.Contains(message, "success") {
		fmt.Println("response is ", message)
		var response ListMealsResponse
		err := json.Unmarshal([]byte(message), &response)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "error on getting meals"})
			return
		}
		c.IndentedJSON(http.StatusOK, response)
		return
	}
}

type WsAddMealMessage struct {
	ReceiverName string `json:"reciever_name"`
	MessageID    string `json:"message_id"`
	ActionName   string `json:"action_name"`
	MealID       int    `json:"meal_id"`
}
type Param struct {
	MealID int `json:"meal_id"`
}
type WsAddMealResponseMessage struct {
	Result       string  `json:"result"`
	ReceiverName string  `json:"reciever_name"`
	MessageID    string  `json:"message_id"`
	MealPrice    float32 `json:"meal_price"`
	RestaurantID int     `json:"restaurant_id"`
}

func AddMealToBasket(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid auth"})
		return
	}
	userID, err := parseToken(token)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid token"})
		return
	}
	var param Param
	err = c.BindJSON(&param)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid params"})
		return
	}
	fmt.Println(userID)
	//check if food avaliable
	newUUID, err := uuid.NewRandom()
	if err != nil {
		fmt.Println("error generated for generating new uuid in restaurateur")
		fmt.Println(err)
	}
	var wsAddMealMessage = WsAddMealMessage{
		ReceiverName: "restaurant",
		MessageID:    newUUID.String(),
		ActionName:   "add_meal",
		MealID:       param.MealID,
	}
	SendAddMealMessage(wsAddMealMessage)
	messageChannel := make(chan string)
	go ReadMessage(messageChannel, newUUID.String())

	message := <-messageChannel
	if strings.Contains(message, "error") {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Error generated while getting food info"})
		return
	}
	db, err := createDb()
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Error generated opening db"})
		return
	}
	var response WsAddMealResponseMessage
	err = json.Unmarshal([]byte(message), &response)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Error generated marshalling on user"})
		return
	}
	query := "insert into basket (meal_id, user_id, price, res_id) values(?,?,?,?)"

	_, err = db.Exec(query, param.MealID, userID, response.MealPrice, response.RestaurantID)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Error adding to basket db"})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"success": "meal added succesfully"})

}
