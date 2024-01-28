package main

import (
	"sync"

	"github.com/geoffjay/plantd/core"

	log "github.com/sirupsen/logrus"
)

type logConfig struct {
	Debug     bool   `mapstructure:"debug"`
	Formatter string `mapstructure:"formatter"`
	Level     string `mapstructure:"level"`
}

type appConfig struct {
	core.Config
	Env            string    `mapstructure:"env"`
	ClientEndpoint string    `mapstructure:"client-endpoint"`
	Log            logConfig `mapstructure:"log"`
}

var lock = &sync.Mutex{}
var instance *appConfig

func getConfig() *appConfig {
	if instance == nil {
		lock.Lock()
		defer lock.Unlock()
		if instance == nil {
			if err := core.LoadConfig("app", &instance); err != nil {
				log.Fatalf("error reading config file: %s\n", err)
			}
		}
	}

	return instance
}
