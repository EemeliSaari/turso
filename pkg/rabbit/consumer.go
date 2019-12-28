package rabbit

import "github.com/streadway/amqp"

// Consumer...
type Consumer struct {
	*RabbitMQ
	Tag       string
	AutoAck   bool
	NoLocal   bool
	NoWait    bool
	Exclusive bool
	handler   func(amqp.Delivery)
	running   chan bool
}

// NewConsumer...
func NewConsumer(config ConsumerConfig) (*Consumer, error) {
	mq, err := New(*config.Config)
	if err != nil {
		return nil, err
	}
	return &Consumer{
		RabbitMQ:  mq,
		AutoAck:   config.AutoAck,
		NoLocal:   config.NoLocal,
		NoWait:    config.NoWait,
		Exclusive: config.Exclusive,
		running:   make(chan bool),
	}, nil
}

// Consume...
func (c Consumer) Consume(handler func(d amqp.Delivery)) error {
	deliveries, err := c.Channel.Consume(
		c.Queue.Name,
		c.Tag,
		c.AutoAck,
		c.Exclusive,
		c.NoLocal,
		c.NoWait,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for {
			for delivery := range deliveries {
				handler(delivery)
			}
		}
	}()
	c.running <- true

	return nil
}
