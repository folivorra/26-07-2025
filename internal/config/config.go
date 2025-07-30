package config

import (
	"github.com/caarlos0/env"
	"log"
	"time"
)

type Config struct {
	Port           string        `env:"PORT" envDefault:"8080"`
	Timeout        time.Duration `env:"TIMEOUT" envDefault:"5s"`
	MaxTasks       uint64        `env:"MAX_TASKS" envDefault:"3"`
	MaxFilesInTask uint64        `env:"MAX_FILES" envDefault:"3"`
	ArchDir        string        `env:"ARCH_DIR" envDefault:"archives"`
	DownloadDir    string        `env:"DOWNLOAD_DIR" envDefault:"downloads"`
	WorkersNum     int           `env:"WORKERS_NUM" envDefault:"3"`
}

func NewConfig() Config {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		log.Fatalln(err, "env parse error")
	}
	return cfg
}
