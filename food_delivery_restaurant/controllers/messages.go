package controllers

import (
	"fmt"
	m "food-delivery/restaurant/models"
)

type WsResponseMessage struct {
	Result       string `json:"result"`
	ReceiverName string `json:"reciever_name"`
	MessageID    string `json:"message_id"`
}

type WsAddMealResponseMessage struct {
	Result       string  `json:"result"`
	ReceiverName string  `json:"reciever_name"`
	MessageID    string  `json:"message_id"`
	MealPrice    float32 `json:"meal_price"`
	RestaurantID int     `json:"restaurant_id"`
}
type WsListRestaurantResponse struct {
	Result       string         `json:"result"`
	ReceiverName string         `json:"reciever_name"`
	MessageID    string         `json:"message_id"`
	Restaurants  []m.Restaurant `json:"restaurants"`
}
type WsMessage struct {
	ReceiverName   string `json:"reciever_name"`
	MessageID      string `json:"message_id"`
	ActionName     string `json:"action_name"`
	RestaurantName string `json:"restaurant_name"`
	UserToken      string `json:"user_token"`
}
type WsListMealsResponse struct {
	Result       string   `json:"result"`
	ReceiverName string   `json:"reciever_name"`
	MessageID    string   `json:"message_id"`
	Meals        []m.Meal `json:"meals"`
}
type WsAddMealMessage struct {
	ReceiverName string `json:"reciever_name"`
	MessageID    string `json:"message_id"`
	ActionName   string `json:"action_name"`
	MealID       int    `json:"meal_id"`
}

func sendErrorMessage(messageID string, receiverName string) {
	var response = WsResponseMessage{
		Result:       "error",
		ReceiverName: receiverName,
		MessageID:    messageID,
	}
	err := Conn.WriteJSON(response)
	if err != nil {
		fmt.Println("error on sending response error message in restauant")
		fmt.Println(err)
	}

}
func sendSuccessMessage(messageID string, receiverName string) {
	var response = WsResponseMessage{
		Result:       "success",
		ReceiverName: receiverName,
		MessageID:    messageID,
	}

	err := Conn.WriteJSON(response)
	if err != nil {
		fmt.Println("error on sending success message in restauant")
		fmt.Println(err)
	}
}

func SendListRestaurantMessage(response WsListRestaurantResponse) {
	err := Conn.WriteJSON(response)
	if err != nil {
		fmt.Println("error on sending list restaurant message in restauant")
		fmt.Println(err)
	}
}
func SendListMealsMessage(response WsListMealsResponse) {

	err := Conn.WriteJSON(response)
	if err != nil {
		fmt.Println("error on sending list meals in restauant")
		fmt.Println(err)
	}

}
func SendAddMealMessage(response WsAddMealResponseMessage) {
	err := Conn.WriteJSON(response)
	if err != nil {
		fmt.Println("error on sending list meals in restauant")
		fmt.Println(err)
	}
}
