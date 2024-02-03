package config

import (
	"time"

	"github.com/gofiber/fiber/v2/middleware/session"
)

type sessionConfig struct {
	Expiration     int    `mapstructure:"expiration"`
	KeyLookup      string `mapstructure:"key-lookup"`
	CookieSecure   bool   `mapstructure:"cookie-secure"`
	CookieHTTPOnly bool   `mapstructure:"cookie-http-only"`
	CookieSameSite string `mapstructure:"cookie-same-site"`
}

func (c *sessionConfig) setDefaults() {
	// FIXME: can't default int this way if 0 is a valid value
	//
	// if c.Expiration == 0 {
	// 	c.Expiration = 30 * time.Minute
	// }

	if c.KeyLookup == "" {
		c.KeyLookup = "cookie:__Host-session"
	}

	// FIXME: can't default bool this way
	//
	// if !c.CookieSecure {
	// 	c.CookieSecure = true
	// }
	//
	// if !c.CookieHTTPOnly {
	// 	c.CookieHTTPOnly = true
	// }

	if c.CookieSameSite == "" {
		c.CookieSameSite = "Lax"
	}
}

func (c *sessionConfig) ToSessionConfig() session.Config {
	return session.Config{
		// Expiration:     c.Expiration,
		Expiration: 30 * time.Minute,
		KeyLookup:  c.KeyLookup,
		// CookieSecure:   c.CookieSecure,
		// CookieHTTPOnly: c.CookieHTTPOnly,
		CookieSecure:   true,
		CookieHTTPOnly: true,
		CookieSameSite: c.CookieSameSite,
	}
}
