package main

import (
	_ "github.com/geoffjay/plantd/app/docs"
	"github.com/geoffjay/plantd/app/handlers"
	"github.com/geoffjay/plantd/core/util"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func initializeRouter(app *fiber.App) {
	staticContents := util.Getenv("PLANTD_APP_PUBLIC_PATH", "./app/public")

	app.Static("/public", staticContents)

	app.Get("/", handlers.Index)
	app.Post("/register", handlers.Register)
	app.Post("/login", handlers.Login)
	app.Get("/logout", handlers.Logout)

	defaultHandler := func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	}

	api := app.Group("/api")
	api.Get("/docs/*", swagger.HandlerDefault)

	v1 := api.Group("/v1", func(c *fiber.Ctx) error {
		c.Set("Version", "v1")
		return c.Next()
	})
	v1.Get("/status", defaultHandler)
	v1.Get("/errors", defaultHandler)
	v1.Get("/workers", defaultHandler)
	v1.Get("/workers/:id", defaultHandler)
	v1.Get("/info", defaultHandler)
}
