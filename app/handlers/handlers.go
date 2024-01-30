package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	jtoken "github.com/golang-jwt/jwt/v5"
	log "github.com/sirupsen/logrus"
)

// SessionStore app wide session store.
var SessionStore *session.Store

// Index renders the application index page.
//
//	@Summary     Index page
//	@Description The application index page
//	@Tags        pages
func Index(c *fiber.Ctx) error {
	fields := log.Fields{
		"service": "app",
		"context": "handlers.index",
	}

	sess, err := SessionStore.Get(c)
	if err != nil {
		log.Println(err)
	}

	log.WithFields(fields).Debug(sess)

	email := sess.Get("email")
	log.WithFields(fields).Debugf("email: %s", email)
	isAuthenticated := email != nil
	log.WithFields(fields).Debugf("isAuthenticated: %t", isAuthenticated)
	unauthorizedMessage := "You are not logged in"
	authorizedMessage := fmt.Sprintf("Welcome %v", email)

	return c.Render("index", fiber.Map{
		"title":               "Hello, World!",
		"authorizedMessage":   authorizedMessage,
		"unauthorizedMessage": unauthorizedMessage,
		"isAuthenticated":     isAuthenticated,
	}, "layouts/base")
}

// Dashboard renders the dashboard page.
func Dashboard(c *fiber.Ctx) error {
	log.Debugf("ctx: %v", c)
	user := c.Locals("user").(*jtoken.Token)
	log.Debugf("user: %v", user)
	claims := user.Claims.(jtoken.MapClaims)
	email := claims["email"].(string)

	log.Debugf("email: %s", email)
	log.Debugf("claims: %v", claims)

	return c.Render("dashboard", fiber.Map{
		"email": email,
	}, "layouts/base")
}
