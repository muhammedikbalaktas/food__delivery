package controllers

// func GetFoods(c *gin.Context) {
// 	// Create a channel to receive messages from the WebSocket connection
// 	messageChannel := make(chan string)

// 	// Start the Goroutine to read messages from the WebSocket connection
// 	go ReadMessage(messageChannel)

// 	// Wait for a message to be received on the messageChannel
// 	message := <-messageChannel

// 	// Respond to the HTTP request with the received message
// 	c.IndentedJSON(http.StatusOK, map[string]string{"success": message})
// }
