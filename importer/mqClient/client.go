package mqClient

import (
	"github.com/streadway/amqp"
	"os"
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
	return
}

func (c *MQClient) GetMessages() (messages <-chan amqp.Delivery, err error) {
	var q amqp.Queue
	if q, err = c.channel.QueueDeclare(
		"Importer", // name
		false,      // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	); err != nil {
		return
	}
	messages, err = c.channel.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	return
}
