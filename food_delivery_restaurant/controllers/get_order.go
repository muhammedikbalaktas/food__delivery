package controllers

import (
	"encoding/json"
)

type WsOrderMealMessage struct {
	ReceiverName string `json:"reciever_name"`
	MessageID    string `json:"message_id"`
	ActionName   string `json:"action_name"`
	Meals        []Meal `json:"meals"`
}

type Meal struct {
	MealID int
	Price  float32
	ResID  int
}

func GetOrder(message string) {
	var orderMessage WsOrderMealMessage

	err := json.Unmarshal([]byte(message), &orderMessage)
	if err != nil {
		sendErrorMessage(orderMessage.MessageID, "user")
		return
	}
	db, err := createDb()
	if err != nil {
		sendErrorMessage(orderMessage.MessageID, "user")
		return
	}
	defer db.Close()

	stmt, err := db.Prepare("insert into orders ( meal_id, res_id, price) values (?,?,?)")
	if err != nil {
		sendErrorMessage(orderMessage.MessageID, "user")
		return
	}
	defer stmt.Close()
	meals := orderMessage.Meals
	for _, meal := range meals {
		_, err := stmt.Exec(meal.MealID, meal.ResID, meal.Price)
		if err != nil {
			sendErrorMessage(orderMessage.MessageID, "user")
			return
		}
	}
	sendSuccessMessage(orderMessage.MessageID, "user")
}
