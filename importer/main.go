package main

import (
	"DC-homework-1/importer/mqClient"
	"encoding/json"
	"lib/dbClient"
	"log"
)

func main() {
	var db dbClient.Client
	var err error
	if db, err = dbClient.InitClient(); err != nil {
		log.Fatal(err)
	}
	var MQClient mqClient.MQClient
	if MQClient, err = mqClient.InitMQClient(); err != nil {
		log.Fatal(err)
	}
	messages, err := MQClient.GetMessages()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Importer service started")
	for message := range messages {
		var batch []dbClient.Product
		_ = json.Unmarshal(message.Body, &batch)
		if err = db.ImportProductBatch(batch); err != nil {
			log.Fatal(err)
		}
		if err = message.Ack(false); err != nil {
			log.Fatal(err)
		}
	}
}
