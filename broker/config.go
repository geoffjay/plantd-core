package main

import (
	"github.com/geoffjay/plantd/core"
)

type busConfig struct {
	Name     string `mapstructure:"name"`
	Frontend string `mapstructure:"frontend"`
	Backend  string `mapstructure:"backend"`
	Capture  string `mapstructure:"capture"`
}

type logConfig struct {
	Debug     bool   `mapstructure:"debug"`
	Formatter string `mapstructure:"formatter"`
	Level     string `mapstructure:"level"`
}

type brokerConfig struct {
	core.Config
	Env               string      `mapstructure:"env"`
	Endpoint          string      `mapstructure:"endpoint"`
	ClientEndpoint    string      `mapstructure:"client-endpoint"`
	HeartbeatLiveness int         `mapstructure:"heartbeat-liveness"`
	HeartbeatInterval int         `mapstructure:"heartbeat-interval"`
	Buses             []busConfig `mapstructure:"buses"`
	Log               logConfig   `mapstructure:"log"`
}
