package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/khsnnn/diploma_fitness_server/club-service/internal/models"
)

func GetClubs(c *fiber.Ctx) error {
	clubs := []models.Club{}
	return c.JSON(clubs)
}