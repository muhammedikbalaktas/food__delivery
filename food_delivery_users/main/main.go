package main

import (
	c "food-delivery/users/controllers"

	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()
	router.POST("/create_user", c.CreateUser)
	router.POST("/get_user", c.GetUser)
	router.GET("/list_restaurants", c.ListRestaurants)
	router.GET("/list_meals", c.ListMeals)
	router.POST("/add_meal", c.AddMealToBasket)
	router.POST("/order_meal", c.OrderMeal)
	c.Connect()
	router.Run(":9090")
}
