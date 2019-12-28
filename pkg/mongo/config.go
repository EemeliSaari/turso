package mongo

import (
	"context"
	"time"
)

type MongoConfig struct {
	Host            string
	Port            int
	Database        string
	Collection      string
	Context         context.Context
	KeepAlive       bool
	HealthCheckFreq time.Duration
}
