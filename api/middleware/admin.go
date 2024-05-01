package middleware

import (
	"github.com/gofiber/fiber/v2"
	"hotel-reservation/errors"
	"hotel-reservation/types"
)

func AdminAuth(c *fiber.Ctx) error {
	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return errors.ErrUnauthorized()
	}
	if !user.IsAdmin {
		return errors.ErrUnauthorized()
	}
	return c.Next()
}
