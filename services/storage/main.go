package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/streadway/amqp"

	"github.com/EemeliSaari/turso/internal/utils"
	"github.com/EemeliSaari/turso/pkg/rabbit"
	"github.com/EemeliSaari/turso/pkg/rss"
)

type storedArticle struct {
	rss.Article
	Content   string `json:"-"`
	Loaded    bool   `json:"-"`
	Erroneous string `json:"-"`
}

func handler(delivery amqp.Delivery) {
	var article storedArticle

	err := json.Unmarshal(delivery.Body, &article)
	if err != nil {
		log.Print("Failed to deserialize message when parsing.")
	}
	log.Printf("Stored article: %s", article.Title)
}

func main2() {
	var article storedArticle

	file, err := ioutil.ReadFile("data/0a0208a1ba159404ca99cd5f58271fb5.json")

	utils.FailOnError(err, "")
	err = json.Unmarshal([]byte(file), &article)

	fmt.Println(article)
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
