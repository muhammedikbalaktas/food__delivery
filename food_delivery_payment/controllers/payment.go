package controllers

import "encoding/json"

type WsCheckPaymentMessage struct {
	ReceiverName string  `json:"reciever_name"`
	MessageID    string  `json:"message_id"`
	ActionName   string  `json:"action_name"`
	TotalPrice   float32 `json:"total_price"`
}

func CheckPayment(message string) {
	var paymentMessage WsCheckPaymentMessage
	err := json.Unmarshal([]byte(message), &paymentMessage)
	if err != nil {
		var response = Message{
			ReceiverName: "user",
			MessageID:    paymentMessage.MessageID,
			Response:     "error",
		}
		SendResponseMessage(response)
		return
	}
	if paymentMessage.TotalPrice > 150 {
		var response = Message{
			ReceiverName: "user",
			MessageID:    paymentMessage.MessageID,
			Response:     "error",
		}
		SendResponseMessage(response)
		return
	}
	var response = Message{
		ReceiverName: "user",
		MessageID:    paymentMessage.MessageID,
		Response:     "success",
	}
	SendResponseMessage(response)

}
