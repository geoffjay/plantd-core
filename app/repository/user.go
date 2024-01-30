package repository

import (
	"errors"

	"github.com/geoffjay/plantd/app/models"
)

// FindUserByCredentials finds a user by their email and password.
func FindUserByCredentials(email, password string) (*models.User, error) {
	// Here you would query your database for the user with the given email
	if email == "test@example.com" && password == "test12345" {
		return &models.User{
			ID:       1,
			Email:    "test@example.com",
			Password: "test12345",
		}, nil
	}
	return nil, errors.New("user not found")
}
