package main

import (
	"os"
	"reflect"
	"testing"

	"github.com/geoffjay/plantd/core"
)

func buildConfig(t *testing.T, path string) *brokerConfig {
	var config *brokerConfig

	_ = os.Setenv("PLANTD_BROKER_CONFIG", path)
	err := core.LoadConfig("broker", &config)
	if err != nil {
		t.Fatalf("Cannot create configuration: %v", err)
	}

	return config
}

func TestLoadConfig(t *testing.T) {
	emptyConfig := &brokerConfig{
		Env:               "",
		Endpoint:          "",
		HeartbeatLiveness: 0,
		HeartbeatInterval: 0,
		Log: logConfig{
			Debug:     false,
			Formatter: "",
			Level:     "",
		},
	}
	validConfig := &brokerConfig{
		Env:               "testing",
		Endpoint:          "tcp://*:7200",
		ClientEndpoint:    "tcp://localhost:7200",
		HeartbeatLiveness: 3,
		HeartbeatInterval: 2500000,
		Log: logConfig{
			Debug:     true,
			Formatter: "text",
			Level:     "debug",
		},
	}
	cases := []struct {
		fixture string
		want    *brokerConfig
		name    string
	}{
		{
			fixture: "fixtures/config/empty.yaml",
			want:    emptyConfig,
			name:    "EmptyFixture",
		},
		{
			fixture: "fixtures/config/valid.yaml",
			want:    validConfig,
			name:    "ValidYAMLFixture",
		},
		{
			fixture: "fixtures/config/valid.json",
			want:    validConfig,
			name:    "ValidJSONFixture",
		},
		{
			fixture: "fixtures/config/valid.toml",
			want:    validConfig,
			name:    "ValidTOMLFixture",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			config := buildConfig(t, tc.fixture)

			// split up data to make it easier to see where an error occurs
			got := map[string]interface{}{
				"env":      config.Env,
				"endpoint": config.Endpoint,
				"log":      config.Log,
			}

			wanted := map[string]interface{}{
				"env":      tc.want.Env,
				"endpoint": tc.want.Endpoint,
				"log":      tc.want.Log,
			}

			for key, value := range wanted {
				if !reflect.DeepEqual(value, got[key]) {
					t.Errorf("Expected: %v, got: %v", wanted[key], got[key])
				}
			}
		})
	}
}
