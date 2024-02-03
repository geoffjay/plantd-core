package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

// SessionStore app wide session store.
var SessionStore *session.Store

// Index renders the application index page.
//
//	@Summary     Index page
//	@Description The application index page
//	@Tags        pages
func Index(c *fiber.Ctx) error {
	session, err := SessionStore.Get(c)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	loggedIn, _ := session.Get("loggedIn").(bool)
	if !loggedIn {
		return c.Redirect("/login")
	}

	return c.Render("index", fiber.Map{
		"Title": "App",
	}, "layouts/base")
}

// Dashboard renders the dashboard page.
func Dashboard(c *fiber.Ctx) error {
	session, err := SessionStore.Get(c)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	loggedIn, _ := session.Get("loggedIn").(bool)
	if !loggedIn {
		// User is not authenticated, redirect to the login page
		return c.Redirect("/login")
	}

	csrfToken, ok := c.Locals("csrf").(string)
	if !ok {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.Render("dashboard", fiber.Map{
		"Title": "Dashboard",
		"csrf":  csrfToken,
	}, "layouts/base")
}
