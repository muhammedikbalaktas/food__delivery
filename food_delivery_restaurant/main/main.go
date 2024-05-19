package main

import (
	c "food-delivery/restaurant/controllers"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	c.Connect()
	router.Run(":9092")

}
