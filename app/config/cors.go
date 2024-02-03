package config

import (
	"strings"

	"github.com/gofiber/fiber/v2/middleware/cors"
)

type corsConfig struct {
	AllowCredentials bool   `mapstructure:"allow-credentials"`
	AllowOrigins     string `mapstructure:"allow-origins"`
	AllowHeaders     string `mapstructure:"allow-headers"`
	AllowMethods     string `mapstructure:"allow-methods"`
}

func (c *corsConfig) setDefaults() {
	// FIXME: can't default bool this way
	//
	// if !c.AllowCredentials {
	// 	c.AllowCredentials = true
	// }

	if c.AllowOrigins == "" {
		c.AllowOrigins = "*"
	}

	if c.AllowHeaders == "" {
		headers := []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Content-Length",
			"Accept-Language",
			"Accept-Encoding",
			"Connection",
			"Access-Control-Allow-Origin",
		}
		c.AllowHeaders = strings.Join(headers, ",")
	}

	if c.AllowMethods == "" {
		c.AllowMethods = "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS"
	}
}

func (c *corsConfig) ToCorsConfig() cors.Config {
	return cors.Config{
		// AllowCredentials: c.AllowCredentials,
		AllowCredentials: true,
		AllowOrigins:     c.AllowOrigins,
		AllowHeaders:     c.AllowHeaders,
		AllowMethods:     c.AllowMethods,
	}
}
