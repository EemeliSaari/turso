package rabbit

import "github.com/streadway/amqp"

type Publisher struct {
	*RabbitMQ

	ContentType string
	Mandatory   bool
	Immidiate   bool
}

// NewPublisher...
func NewPublisher(config PublisherConfig, contentType string) (*Publisher, error) {
	mq, err := New(*config.Config)
	if err != nil {
		return nil, err
	}
	return &Publisher{
		RabbitMQ:    mq,
		ContentType: contentType,
		Mandatory:   config.Mandatory,
		Immidiate:   config.Immidiate,
	}, nil
}

// Publish...
func (p Publisher) Publish(body []byte) error {
	return p.Channel.Publish(
		"",
		p.Queue.Name,
		p.Mandatory,
		p.Immidiate,
		amqp.Publishing{
			ContentType: p.ContentType,
			Body:        body,
		},
	)
}
