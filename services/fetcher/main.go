package main

import (
	"encoding/json"
	"log"

	"github.com/EemeliSaari/turso/internal/utils"

	"github.com/EemeliSaari/turso/pkg/constants"
	"github.com/EemeliSaari/turso/pkg/rabbit"
	"github.com/EemeliSaari/turso/pkg/rss"
)

type resultHandler struct {
	*rabbit.Publisher
}

func main() {
	forever := make(chan bool)
	conf := rabbit.PublisherConfig{
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
		Mandatory: false,
		Immidiate: false,
	}

	p, err := rabbit.NewPublisher(conf, "text/plain")
	if err != nil {
		panic(err)
	}

	h := resultHandler{Publisher: p}
	feeds := utils.FileLines("feeds.txt")
	lst := rss.NewListener(feeds, constants.DefaultInterval)
	lst.AddCallback(h.handler)

	err = lst.Start()
	if err != nil {
		log.Print("Failed to start the RSS listener.")
		panic(err)
	}
	<-forever
}

func (h resultHandler) handler(articles []*rss.Article) {
	log.Printf("Sending %d articles", len(articles))
	for _, a := range articles {
		bin, err := json.Marshal(a)
		utils.FailOnError(err, "Failed to serialize on fetcher.")
		h.Publish(bin)
	}
}
