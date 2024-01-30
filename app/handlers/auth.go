package handlers

import (
	"fmt"
	"time"

	conf "github.com/geoffjay/plantd/app/config"
	"github.com/geoffjay/plantd/app/repository"
	"github.com/geoffjay/plantd/app/types"

	"github.com/gofiber/fiber/v2"
	jtoken "github.com/golang-jwt/jwt/v5"
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

	config := conf.GetConfig()

	// Extract the credentials from the request body
	loginRequest := new(types.LoginRequest)
	if err := c.BodyParser(loginRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// sess, err := SessionStore.Get(c)
	// if err != nil {
	// 	panic(err)
	// }

	log.WithFields(fields).Debugf("email: %s", loginRequest.Email)
	user, err := repository.FindUserByCredentials(loginRequest.Email, loginRequest.Password)
	if err != nil {
		log.WithFields(fields).Error(err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	day := time.Hour * 24
	claims := jtoken.MapClaims{
		"ID":    user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(day * 1).Unix(),
	}
	token := jtoken.NewWithClaims(jtoken.SigningMethodHS256, claims)
	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(config.Secret))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// log.WithFields(fields).Debugf("logging in: %s", user.Email)
	// sess.Set("email", user.Email)
	// if err := sess.Save(); err != nil {
	// 	panic(err)
	// }

	c.Set("Set-Cookie", fmt.Sprintf("token=%s", t))
	c.Set("HX-Redirect", "/dashboard")

	return c.JSON(types.LoginResponse{
		Token: t,
	})
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
