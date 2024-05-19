package main

import (
	c "food-delivery/restaurateurs/controllers"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.POST("/create_restaurateur", c.CreateRestaurateur)
	router.POST("/get_restaurateur", c.GetRestaurateur)
	router.POST("/create_restaurant", c.CreateRestaurant)
	c.Connect()
	router.Run(":9091")

}
