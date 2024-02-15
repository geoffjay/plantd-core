package handlers

import (
	"net/http"

	"github.com/geoffjay/plantd/app/repository"
	"github.com/geoffjay/plantd/app/types"
	"github.com/geoffjay/plantd/app/views"
	"github.com/geoffjay/plantd/app/views/pages"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

// TODO: add registration page.
func Register(c *fiber.Ctx) error {
	// Validate user input (username, email, password)
	// Hash the password
	// Store user data in the database
	// Return a success message or error response
	return c.Send([]byte("Register"))
}

func LoginPage(c *fiber.Ctx) error {
	csrfToken, ok := c.Locals("csrf").(string)
	if !ok {
		log.Info("csrf token not found")
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	log.Infof("login page with csrf token: %s", csrfToken)

	c.Locals("title", "Login")

	return views.Render(c, pages.Login(), templ.WithStatus(http.StatusOK))
}

func Login(c *fiber.Ctx) error {
	fields := log.Fields{
		"service": "app",
		"context": "handlers.login",
	}

	// Extract the credentials from the request body
	loginRequest := new(types.LoginRequest)
	if err := c.BodyParser(loginRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	log.WithFields(fields).Debugf("email: %s", loginRequest.Email)
	_, err := repository.FindUserByCredentials(loginRequest.Email, loginRequest.Password)
	if err != nil {
		log.WithFields(fields).Error(err)
		csrfToken, ok := c.Locals("csrf").(string)
		if !ok {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		return c.Render("login", fiber.Map{
			"title": "Login",
			"csrf":  csrfToken,
			"error": "Invalid credentials",
		}, "layouts/base")
	}

	log.WithFields(fields).Debugf("logging in: %s", loginRequest.Email)

	session, err := SessionStore.Get(c)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if err := session.Reset(); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	session.Set("loggedIn", true)
	if err := session.Save(); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	c.Set("HX-Redirect", "/")

	return c.SendStatus(fiber.StatusOK)
}

func Logout(c *fiber.Ctx) error {
	session, err := SessionStore.Get(c)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	// Revoke users authentication
	if err := session.Destroy(); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.Redirect("/login")
}
