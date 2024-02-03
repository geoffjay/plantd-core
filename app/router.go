package main

import (
	"time"

	_ "github.com/geoffjay/plantd/app/docs"
	"github.com/geoffjay/plantd/app/handlers"
	"github.com/geoffjay/plantd/core/util"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/swagger"
	log "github.com/sirupsen/logrus"
)

func csrfErrorHandler(c *fiber.Ctx, err error) error {
	// Log the error so we can track who is trying to perform CSRF attacks
	// customize this to your needs
	log.WithFields(log.Fields{
		"service": "app",
		"context": "router.csrfErrorHandler",
		"error":   err,
		"ip":      c.IP(),
		"request": c.OriginalURL(),
	}).Error("CSRF Error")

	log.Debugf("ctx: %v", c)

	// check accepted content types
	switch c.Accepts("html", "json") {
	case "json":
		// Return a 403 Forbidden response for JSON requests
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "403 Forbidden",
		})
	case "html":
		// Return a 403 Forbidden response for HTML requests
		return c.Status(fiber.StatusForbidden).Render("error", fiber.Map{
			"Title":     "Error",
			"Error":     "403 Forbidden",
			"ErrorCode": "403",
		}, "layouts/base")
	default:
		// Return a 403 Forbidden response for all other requests
		return c.Status(fiber.StatusForbidden).SendString("403 Forbidden")
	}
}

func initializeRouter(app *fiber.App) {
	staticContents := util.Getenv("PLANTD_APP_PUBLIC_PATH", "./app/public")

	csrfConfig := csrf.Config{
		Session:        handlers.SessionStore,
		KeyLookup:      "form:csrf",
		CookieName:     "__Host-csrf",
		CookieSameSite: "Lax",
		CookieSecure:   true,
		CookieHTTPOnly: true,
		ContextKey:     "csrf",
		ErrorHandler:   csrfErrorHandler,
		Expiration:     30 * time.Minute,
	}
	csrfMiddleware := csrf.New(csrfConfig)

	app.Static("/public", staticContents)

	app.Get("/", handlers.Index)

	app.Get("/login", csrfMiddleware, handlers.LoginPage)
	app.Post("/login", csrfMiddleware, handlers.Login)
	app.Get("/logout", handlers.Logout)
	app.Post("/register", handlers.Register)

	app.Get("/dashboard", csrfMiddleware, handlers.Dashboard)

	// TODO: this is just here until the API is implemented.
	defaultHandler := func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	}

	api := app.Group("/api")
	api.Get("/docs/*", swagger.HandlerDefault)

	v1 := api.Group("/v1", func(c *fiber.Ctx) error {
		c.Set("Version", "v1")
		return c.Next()
	})

	broker := v1.Group("/broker")
	broker.Get("/status", defaultHandler)
	broker.Get("/errors", defaultHandler)
	broker.Get("/workers", defaultHandler)
	broker.Get("/workers/:id", defaultHandler)
	broker.Get("/info", defaultHandler)
}
