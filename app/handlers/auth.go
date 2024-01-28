package handlers

import (
	"github.com/geoffjay/plantd/app/models"

	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

func Register(c *fiber.Ctx) error {
	// Validate user input (username, email, password)
	// Hash the password
	// Store user data in the database
	// Return a success message or error response
	return c.Send([]byte("Register"))
}

func Login(c *fiber.Ctx) error {
	fields := log.Fields{
		"service": "app",
		"context": "handlers.login",
	}
	user := new(models.User)
	if err := c.BodyParser(user); err != nil {
		return err
	}

	log.WithFields(fields).Debugf("email: %s", user.Email)

	sess, err := SessionStore.Get(c)
	if err != nil {
		panic(err)
	}

	log.WithFields(fields).Debugf("logging in: %s", user.Email)
	sess.Set("email", user.Email)
	if err := sess.Save(); err != nil {
		panic(err)
	}

	return c.Redirect("/")
}

func Logout(c *fiber.Ctx) error {
	sess, err := SessionStore.Get(c)
	if err != nil {
		panic(err)
	}

	sess.Delete("name")

	if err := sess.Destroy(); err != nil {
		panic(err)
	}

	return c.Redirect("/")
}
