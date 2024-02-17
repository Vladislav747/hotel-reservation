package middleware

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
)

func JWTAuthentication(c *fiber.Ctx) error {
	fmt.Println("-- JWT auth")
	return nil
}
