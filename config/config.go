package config

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	PostgresDSN     string        `envconfig:"POSTGRES_DSN" required:"true"`
	ListenAddress   string        `default:":8080" split_words:"true"`
	ClientTimeout   time.Duration `default:"5s" split_words:"true"`
	ReadTimeout     time.Duration `default:"5s" split_words:"true"`
	WriteTimeout    time.Duration `default:"500s" split_words:"true"`
	ShutdownTimeout time.Duration `default:"5s" split_words:"true"`
	LogLevel        string        `default:"info" split_wors:"true"`
	JsonLogOutput   bool          `required:"false" default:"true" split_words:"true"`
}

func New() (*Config, error) {
	cfg := &Config{}

	if err := envconfig.Process("app", cfg); err != nil {
		return nil, fmt.Errorf("parse error: %w", err)
	}

	return cfg, nil
}
