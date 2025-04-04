package api

import (
	"github.com/gofiber/fiber/v2"
)

func GetClubs(c *fiber.Ctx) error {
	return c.JSON([]interface{}{})
}