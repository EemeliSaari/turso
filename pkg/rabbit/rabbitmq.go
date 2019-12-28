package rabbit

import (
	"github.com/streadway/amqp"
)

// RabbitMQ...
type RabbitMQ struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
	Queue      *amqp.Queue
}

// New...
func New(config Config) (*RabbitMQ, error) {

	uri := amqp.URI{
		Scheme:   "amqp",
		Host:     config.Host,
		Port:     config.Port,
		Username: config.Username,
		Password: config.Password,
		Vhost:    config.Vhost,
	}.String()

	conn, err := amqp.Dial(uri)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	q, err := ch.QueueDeclare(
		config.Queue,
		config.Durable,
		config.DeleteUnused,
		config.Exclusive,
		config.NoWait,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &RabbitMQ{
		Connection: conn,
		Channel:    ch,
		Queue:      &q,
	}, nil
}
