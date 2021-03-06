package mqClient

import (
	"encoding/json"
	"lib/dbClient"
	"os"

	"github.com/streadway/amqp"
)

type MQClient struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func InitMQClient() (mqClient MQClient, err error) {
	if mqClient.conn, err = amqp.Dial(os.Getenv("MESSAGE_QUEUE_ADDRESS")); err != nil {
		return
	}
	if mqClient.channel, err = mqClient.conn.Channel(); err != nil {
		return
	}
	_, err = mqClient.channel.QueueDeclare(
		"Importer", // name
		false,      // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	return
}

func (c *MQClient) SendBatch(batch []dbClient.Product) (err error) {
	var bytes []byte
	if bytes, err = json.Marshal(batch); err != nil {
		return
	}
	err = c.channel.Publish(
		"",         // exchange
		"Importer", // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        bytes,
		})
	return
}
