package config

import (
	"sync"

	"github.com/geoffjay/plantd/core"

	log "github.com/sirupsen/logrus"
)

// TODO:
// - add a new configuration section for the session middleware
// - add a new configuration section for the database
// - assign defaults to the configuration values

type logConfig struct {
	Debug     bool   `mapstructure:"debug"`
	Formatter string `mapstructure:"formatter"`
	Level     string `mapstructure:"level"`
}

type Config struct {
	core.Config
	Env            string        `mapstructure:"env"`
	ClientEndpoint string        `mapstructure:"client-endpoint"`
	Secret         string        `mapstructure:"secret"`
	Log            logConfig     `mapstructure:"log"`
	Cors           corsConfig    `mapstructure:"cors"`
	Session        sessionConfig `mapstructure:"session"`
}

var lock = &sync.Mutex{}
var instance *Config

func (c *Config) setDefaults() {
	if c.Env == "" {
		c.Env = "development"
	}

	if c.ClientEndpoint == "" {
		c.ClientEndpoint = "tcp://localhost:9797"
	}

	if c.Secret == "" {
		c.Secret = "secret"
	}

	if c.Log.Formatter == "" {
		c.Log.Formatter = "text"
	}

	if c.Log.Level == "" {
		c.Log.Level = "info"
	}

	c.Cors.setDefaults()
	c.Session.setDefaults()
}

// GetConfig returns the application configuration singleton.
func GetConfig() *Config {
	if instance == nil {
		lock.Lock()
		defer lock.Unlock()
		if instance == nil {
			if err := core.LoadConfig("app", &instance); err != nil {
				log.Fatalf("error reading config file: %s\n", err)
			}
			instance.setDefaults()
		}
	}

	log.Tracef("config: %+v", instance)

	return instance
}
