package config

import (
	"github.com/caarlos0/env"
	"log"
	"time"
)

type Config struct {
	Port     string        `env:"PORT" envDefault:"8080"`
	Timeout  time.Duration `env:"TIMEOUT" envDefault:"5s"`
	MaxTasks int           `env:"MAX_TASKS" envDefault:"3"`
	MaxFiles int           `env:"MAX_FILES" envDefault:"3"`
}

func NewConfig() *Config {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		log.Fatalln(err, "env parse error")
	}
	return &cfg
}
