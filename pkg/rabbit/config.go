package rabbit

// Config...
type Config struct {
	Host         string
	Port         int
	Username     string
	Password     string
	Vhost        string
	Queue        string
	Durable      bool
	DeleteUnused bool
	Exclusive    bool
	NoWait       bool
}

// PublisherConfig...
type PublisherConfig struct {
	*Config

	Mandatory bool
	Immidiate bool
}

// ConsumerConfig...
type ConsumerConfig struct {
	*Config

	AutoAck bool
	NoLocal bool
	NoWait  bool
	Tag     string
}
