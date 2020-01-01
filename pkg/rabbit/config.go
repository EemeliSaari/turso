package rabbit

// Config ...
type Config struct {
	Host         string `json:"host"`
	Port         int    `json:"port"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	Vhost        string `json:"vhost"`
	Queue        string `json:"queue"`
	Durable      bool   `json:"durable"`
	DeleteUnused bool   `json:"delete_unused"`
	Exclusive    bool   `json:"exlusive"`
	NoWait       bool   `json:"no_wait"`
}

// PublisherConfig ...
type PublisherConfig struct {
	*Config

	Mandatory bool
	Immidiate bool
}

// ConsumerConfig ...
type ConsumerConfig struct {
	*Config

	AutoAck bool
	NoLocal bool
	NoWait  bool
	Tag     string
}
