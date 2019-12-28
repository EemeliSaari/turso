package mongo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoClient struct {
	*mongo.Client
	Collection      *mongo.Collection
	Context         context.Context
	KeepAlive       bool
	HealthCheckFreq int // Ping frequency in seconds

	healthy chan bool
	ticker  *time.Ticker
}

func New(config MongoConfig) (*MongoClient, error) {
	uri := fmt.Sprintf("mongodb://%s:%d", config.Host, config.Port)
	client, err := mongo.Connect(config.Context, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	collection := client.Database(config.Database).Collection(config.Collection)

	c := MongoClient{
		Client:     client,
		Collection: collection,
		Context:    config.Context,
		KeepAlive:  config.KeepAlive,
		healthy:    make(chan bool),
		ticker:     time.NewTicker(config.HealthCheckFreq * time.Second),
	}

	if config.KeepAlive {
		go func() {
			for {
				select {
				case <-c.ticker.C:
					err := c.Ping(nil)
					c.healthy <- err == nil
				}
			}
		}()
	} else {
		if err = c.Ping(nil); err != nil {
			return nil, err
		}
	}
	return &c, nil
}

func (c MongoClient) Ping(rp *readpref.ReadPref) error {
	return c.Client.Ping(c.Context, rp)
}

func (c MongoClient) InsertMany(data []interface{}) (*mongo.InsertManyResult, error) {
	if err := c.CheckHealth(); err != nil {
		return nil, err
	}
	return c.Collection.InsertMany(c.Context, data)
}

func (c MongoClient) InsertOne(data interface{}) (*mongo.InsertOneResult, error) {
	if err := c.CheckHealth(); err != nil {
		return nil, err
	}
	return c.Collection.InsertOne(c.Context, data)
}

func (c MongoClient) CheckHealth() error {
	status := <-c.healthy
	if !status {
		return errors.New("MongoDB connection unhealthy")
	} else {
		return nil
	}
}
