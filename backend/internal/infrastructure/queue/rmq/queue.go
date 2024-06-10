package rmq

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RMQClient struct {
	cc *amqp.Connection
}

// NewClient creates a new RabbitMQ client. Will panic if there are an error while dealing
func NewClient(
	address string,
	user string,
	password string,
) *RMQClient {
	cc, err := amqp.Dial("amqp://" + user + ":" + password + "@localhost:5672/")
	if err != nil {
		log.Fatal("error connect to rabbitmq server", err)
	}

	return &RMQClient{
		cc: cc,
	}
}

func NewWithConnection(
	cc *amqp.Connection,
) *RMQClient {
	return &RMQClient{
		cc: cc,
	}
}
