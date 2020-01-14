package main

import (
	"encoding/json"
	"log"

	"github.com/EemeliSaari/turso/internal/utils"

	"github.com/streadway/amqp"

	"github.com/EemeliSaari/turso/pkg/document"
	"github.com/EemeliSaari/turso/pkg/rabbit"
	"github.com/EemeliSaari/turso/pkg/rss"
)

type resultHandler struct {
	*rabbit.Publisher
}

func (h resultHandler) handler(delivery amqp.Delivery) {
	var article rss.Article

	err := json.Unmarshal(delivery.Body, &article)
	if err != nil {
		log.Print("Failed to deserialize message when parsing.")
	}

	doc, err := document.New([]byte(article.Content))
	if err != nil {
		log.Print("Failed to parse document:", article.Title)
	}

	bin, err := json.Marshal(doc)
	if err != nil {
		log.Print("Failed to serialize parsed document.")
	}

	h.Publish(bin)
}

func main() {
	base := rabbit.Config{
		Host:         "localhost",
		Port:         5672,
		Username:     "dev",
		Password:     "dev",
		Vhost:        "/",
		Queue:        "data",
		Durable:      false,
		DeleteUnused: false,
		Exclusive:    false,
		NoWait:       false,
	}
	consumerConf := rabbit.ConsumerConfig{
		Config:  &base,
		AutoAck: true,
		NoLocal: false,
		NoWait:  false,
		Tag:     "",
	}
	publisherConf := rabbit.PublisherConfig{
		Config:    &base,
		Mandatory: false,
		Immidiate: false,
	}

	c, err := rabbit.NewConsumer(consumerConf)
	utils.FailOnError(err, "failed to start consumer")
	p, err := rabbit.NewPublisher(publisherConf, "text/plain")
	utils.FailOnError(err, "failed to start publisher")

	h := resultHandler{Publisher: p}

	c.Consume(h.handler)
}
