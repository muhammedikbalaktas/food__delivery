package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type WsCheckPaymentMessage struct {
	ReceiverName string  `json:"reciever_name"`
	MessageID    string  `json:"message_id"`
	ActionName   string  `json:"action_name"`
	TotalPrice   float32 `json:"total_price"`
}
type WsOrderMealMessage struct {
	ReceiverName string `json:"reciever_name"`
	MessageID    string `json:"message_id"`
	ActionName   string `json:"action_name"`
	Meals        []Meal `json:"meals"`
}

func OrderMeal(c *gin.Context) {
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
	fmt.Println(userID)
	newUUID, err := uuid.NewRandom()
	if err != nil {

		fmt.Println("error generated for generating new uuid in restaurateur")
		fmt.Println(err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "error on marshalling"})
		return
	}
	meals, err := calculatePrice(userID)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var totalPrice float32
	for _, meal := range meals {
		totalPrice += meal.Price
	}
	var paymentMessage = WsCheckPaymentMessage{
		ReceiverName: "payment",
		MessageID:    newUUID.String(),
		ActionName:   "check_payment",
		TotalPrice:   totalPrice,
	}
	SendCheckPaymentMessage(paymentMessage)
	messageChannel := make(chan string)
	go ReadMessage(messageChannel, newUUID.String())

	message := <-messageChannel
	if strings.Contains(message, "error") {
		fmt.Println("error generated for payment")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "error on payment. Invalid amount of money"})
		return
	}
	//send order message should be here
	newUUID, err = uuid.NewRandom()
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "error on executing statements"})
		return
	}
	var orderMessage = WsOrderMealMessage{
		ReceiverName: "restaurant",
		MessageID:    newUUID.String(),
		ActionName:   "order_meal",
		Meals:        meals,
	}
	SendOrderMealMessage(orderMessage)

	go ReadMessage(messageChannel, newUUID.String())

	message = <-messageChannel
	if strings.Contains(message, "error") {
		fmt.Println("error generated for adding meal in side of restaurant")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "error on restaurant side "})
		return
	}
	fmt.Println("stucked")
	fmt.Println("message is ", message)
	uid := uuid.New()
	db, err := createDb()
	if err != nil {
		fmt.Println(err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "error on creating database"})
		return
	}
	defer db.Close()
	query := "insert into orders (user_id,total_amount,uid) values(?,?,?)"
	_, err = db.Exec(query, userID, totalPrice, uid.String())
	if err != nil {
		fmt.Println(err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "error on creating database"})
		return
	}
	query = "select id from orders where uid=?"
	rows, err := db.Query(query, uid)
	if err != nil {
		fmt.Println(err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "error on scanning rows"})
		return
	}
	var orderID int
	if rows.Next() {
		err := rows.Scan(&orderID)
		if err != nil {
			fmt.Println(err)
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "error on scanning rows"})
			return
		}
	}
	stmt, err := db.Prepare("insert into order_details (order_id, meal_id, price) values (?,?,?)")
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "error on creating statements"})
		return
	}
	defer stmt.Close()

	for _, meal := range meals {
		_, err := stmt.Exec(orderID, meal.MealID, meal.Price)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "error on executing statements"})
			return
		}
	}

	query = "delete from basket where user_id=?"
	_, err = db.Exec(query, userID)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "error on removing items from basket"})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"succes": "order has been made succesfully"})
}

type Meal struct {
	MealID int
	Price  float32
	ResID  int
}

func calculatePrice(userID int) ([]Meal, error) {
	db, err := createDb()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := "select b.meal_id, b.price, b.res_id from users user " +
		"inner join basket b on user.id=b.user_id where user.id=?"

	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	var meals = make([]Meal, 0)
	rowCount := 0
	for rows.Next() {
		var meal Meal
		err := rows.Scan(&meal.MealID, &meal.Price, &meal.ResID)
		if err != nil {
			return nil, err
		}
		meals = append(meals, meal)
		rowCount++
	}
	if rowCount == 0 {
		return nil, errors.New("the basket is empty")
	}
	return meals, nil

}
