package middleware

// import (
// 	"github.com/gofiber/fiber/v2"
// 	"github.com/golang-jwt/jwt"
// )
//
// func AuthMiddleware(c *fiber.Ctx) error {
// 	// Check for authentication credentials
// 	// If credentials are valid, proceed to the next middleware or route
// 	// If credentials are invalid, return an unauthorized response
//
// 	tokenString := c.Get("Authorization")
//
// 	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 		// Verify the token signing method and return the secret key
// 	})
//
// 	if err != nil || !token.Valid {
// 		return c.Status(fiber.StatusUnauthorized).SendString("Invalid Token")
// 	}
//
// 	// Extract user information from the token and store it in the context
// 	c.Locals("user", getUserFromToken(token))
//
// 	return c.Next()
// }
