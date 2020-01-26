package couch

import "fmt"

// Config ...
type Config struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
}

func (c Config) asURL() string {
	return fmt.Sprintf("http://%s:%d/%s/", c.Host, c.Port, c.Database)
}
