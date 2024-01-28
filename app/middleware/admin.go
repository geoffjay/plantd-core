package middleware

// import (
// 	"github.com/geoffjay/plantd/app/models"
//
// 	"github.com/gofiber/fiber/v2"
// )
//
// func AdminMiddleware(c *fiber.Ctx) error {
// 	userRole := getUserRoleFromContext(c)
//
// 	if userRole != models.AdminRole {
// 		return c.Status(fiber.StatusForbidden).SendString("Permission Denied")
// 	}
//
// 	return c.Next()
// }
