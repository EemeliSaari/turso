package main

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"

	"github.com/EemeliSaari/turso/pkg/rabbit"
	"github.com/EemeliSaari/turso/pkg/rss"
)

func handler(delivery amqp.Delivery) {
	var article rss.Article

	err := json.Unmarshal(delivery.Body, &article)
	if err != nil {
		log.Print("Failed to deserialize message when parsing.")
	}
	log.Printf("Parsed article: %s", article.Title)
}

func main() {
	conf := rabbit.ConsumerConfig{
		Config: &rabbit.Config{
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
		},
		AutoAck: true,
		NoLocal: false,
		NoWait:  false,
		Tag:     "",
	}

	c, err := rabbit.NewConsumer(conf)
	if err != nil {
		panic(err)
	}
	c.Consume(handler)
}
