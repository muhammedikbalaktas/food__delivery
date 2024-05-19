package controllers

import (
	"encoding/json"
	m "food-delivery/restaurant/models"

	"fmt"
)

func ListMeals(wsJsonMessage string) {
	var wsMessage WsMessage
	err := json.Unmarshal([]byte(wsJsonMessage), &wsMessage)
	messageID := wsMessage.MessageID
	if err != nil {
		fmt.Println("error on marshalling json get from user in restaurant")
		fmt.Println(err)
		sendErrorMessage(messageID, "user")
		return
	}
	db, err := createDb()

	if err != nil {
		sendErrorMessage(messageID, "user")
		fmt.Println(err)
		return
	}
	defer db.Close()
	var meals = make([]m.Meal, 0)

	query := "select res.name, m.id, m.name, m.price, m.amount from restaurant res " +
		"inner join meals m on res.id=m.res_id"
	rows, err := db.Query(query)
	if err != nil {
		sendErrorMessage(messageID, "user")
		return
	}
	for rows.Next() {
		var meal m.Meal
		err := rows.Scan(&meal.ResName, &meal.ID, &meal.Name, &meal.Price, &meal.Amount)
		if err != nil {
			sendErrorMessage(messageID, "user")
			return
		}
		meals = append(meals, meal)
	}

	var mealsResponse = WsListMealsResponse{
		Result:       "success",
		ReceiverName: "user",
		MessageID:    wsMessage.MessageID,
		Meals:        meals,
	}
	SendListMealsMessage(mealsResponse)
}

func CheckAvailability(wsJsonMessage string) {
	// type WsAddMealMessage struct {
	// 	ReceiverName string `json:"reciever_name"`
	// 	MessageID    string `json:"message_id"`
	// 	ActionName   string `json:"action_name"`
	// 	MealID       int    `json:"meal_id"`
	// }
	var wsAddMealMessage WsAddMealMessage
	err := json.Unmarshal([]byte(wsJsonMessage), &wsAddMealMessage)
	messageID := wsAddMealMessage.MessageID
	if err != nil {
		fmt.Println(err)
		sendErrorMessage(messageID, "user")
		return
	}
	db, err := createDb()
	if err != nil {
		fmt.Println(err)
		sendErrorMessage(messageID, "user")
		return
	}
	query := "select amount, price, res_id from meals where id=?"
	var amount int
	var price float32
	rows, err := db.Query(query, wsAddMealMessage.MealID)
	if err != nil {
		fmt.Println(err)
		sendErrorMessage(messageID, "user")
		return
	}
	var resID int
	if rows.Next() {
		err := rows.Scan(&amount, &price, &resID)
		if err != nil {
			fmt.Println(err)
			sendErrorMessage(messageID, "user")
			return
		}
	}
	if amount <= 0 {
		fmt.Println("meal is not available")
		sendErrorMessage(messageID, "user")
		return
	}
	SendAddMealMessage(WsAddMealResponseMessage{
		Result:       "success",
		ReceiverName: "user",
		MessageID:    wsAddMealMessage.MessageID,
		MealPrice:    price,
		RestaurantID: resID,
	})
	fmt.Println("message has been sent")
}
