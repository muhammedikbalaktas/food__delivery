package controllers

import (
	"encoding/json"
	"fmt"
	m "food-delivery/restaurant/models"
)

func ListRestaurants(wsJsonMessage string) {

	var wsMessage WsMessage
	err := json.Unmarshal([]byte(wsJsonMessage), &wsMessage)
	messageID := wsMessage.MessageID
	if err != nil {
		fmt.Println("error on marshalling json get from user in restaurant")
		fmt.Println(err)
		sendErrorMessage(messageID, "user")
		return
	}
	var restaurants = make([]m.Restaurant, 0)
	db, err := createDb()
	if err != nil {
		fmt.Println("error on opening database in side of ListRestaurants in restaurant")
		sendErrorMessage(messageID, "user")
		fmt.Println(err)
		return
	}
	defer db.Close()
	query := "select id, name from restaurant"
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("error on getting rows in side of ListRestaurants in restaurant")

		fmt.Println(err)
		sendErrorMessage(messageID, "user")
		return
	}
	for rows.Next() {
		var restaurant m.Restaurant
		err := rows.Scan(&restaurant.ID, &restaurant.Name)
		if err != nil {
			fmt.Println("error on scanning rows in side of ListRestaurants in restaurant")

			fmt.Println(err)
			sendErrorMessage(messageID, "user")
			return
		}
		restaurants = append(restaurants, restaurant)
	}
	var wsListResponse = WsListRestaurantResponse{
		Result:       "success",
		ReceiverName: "user",
		MessageID:    wsMessage.MessageID,
		Restaurants:  restaurants,
	}
	SendListRestaurantMessage(wsListResponse)
}

func CreateRestaurant(wsJsonMessage string) {
	//parse message first
	var wsMessage WsMessage
	err := json.Unmarshal([]byte(wsJsonMessage), &wsMessage)
	messageID := wsMessage.MessageID
	if err != nil {
		fmt.Println("error on marshalling json get from user in restaurant")
		fmt.Println(err)
		sendErrorMessage(messageID, "restaurateur")
		return
	}

	db, err := createDb()
	if err != nil {
		fmt.Println("error on opening database on restaurant")
		sendErrorMessage(messageID, "restaurateur")
		return
	}
	defer db.Close()

	userID, err := parseToken(wsMessage.UserToken)
	if err != nil {
		fmt.Println("could not validate user in side of restaurant")
		sendErrorMessage(messageID, "restaurateur")
		fmt.Println(err)
		return
	}

	query := "insert into restaurant (name, owner_id) values(?,?)"
	_, err = db.Exec(query, wsMessage.RestaurantName, userID)
	if err != nil {
		sendErrorMessage(messageID, "restaurateur")
	}
	sendSuccessMessage(messageID, "restaurateur")
}
