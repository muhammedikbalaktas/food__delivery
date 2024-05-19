package main

import (
	c "food-delivery/controllers"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/ws", c.HandleConnections)
	router.Run(":8080")

	// token, err := CreateToken()
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println(token)
}
