package http

import (
	"fmt"
)

type Config struct {
	Host string `env:"HOST,notEmpty"`
	Port uint   `env:"PORT,notEmpty"`
}

func (c *Config) Addr() string {
	return fmt.Sprintf("%s:%v", c.Host, c.Port)
}
