package main

import (
	"DC-homework-1/notification/mqClient"
	"DC-homework-1/notification/sms"
	"encoding/json"
	"log"
)

type Message struct {
	Receiver string `json:"receiver"`
	Text     string `json:"text"`
}

func main() {
	smsClient := sms.InitClient()
	var MQClient mqClient.MQClient
	var err error
	if MQClient, err = mqClient.InitMQClient(); err != nil {
		log.Fatal(err)
	}
	messages, err := MQClient.GetMessages()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Notification service started")
	for message := range messages {
		var notification Message
		_ = json.Unmarshal(message.Body, &notification)
		if err = smsClient.Send(notification.Receiver, notification.Text); err != nil {
			log.Fatal(err)
		}
		if err = message.Ack(false); err != nil {
			log.Fatal(err)
		}
	}
}
