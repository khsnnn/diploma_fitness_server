package api

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/khsnnn/diploma_fitness_server/club-service/internal/db"
)

type Handler struct {
	db *db.DB
}

func NewHandler(db *db.DB) *Handler {
	return &Handler{db: db}
}

func (h *Handler) GetClubs(c *fiber.Ctx) error {
	filter := db.ClubFilter{
		Distance: 10.0,
	}

	if latStr := c.Query("lat"); latStr != "" {
		lat, err := strconv.ParseFloat(latStr, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid latitude"})
		}
		filter.Lat = &lat
	}

	if lonStr := c.Query("lon"); lonStr != "" {
		lon, err := strconv.ParseFloat(lonStr, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid longitude"})
		}
		filter.Lon = &lon
	}

	if distStr := c.Query("distance"); distStr != "" {
		distance, err := strconv.ParseFloat(distStr, 64)
		if err != nil || distance <= 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid distance"})
		}
		filter.Distance = distance
	}

	if ratingStr := c.Query("min_rating"); ratingStr != "" {
		rating, err := strconv.ParseFloat(ratingStr, 64)
		if err != nil || rating < 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid minimum rating"})
		}
		filter.MinRating = &rating
	}

	if categoriesStr := c.Query("categories"); categoriesStr != "" {
		categories := strings.Split(categoriesStr, ",")
		if len(categories) > 0 {
			filter.Categories = categories
		}
	}

	if subcategoriesStr := c.Query("subcategories"); subcategoriesStr != "" {
		subcategories := strings.Split(subcategoriesStr, ",")
		if len(subcategories) > 0 {
			filter.Subcategories = subcategories
		}
	}

	if typeStr := c.Query("type"); typeStr != "" {
		if typeStr != "commercial" && typeStr != "university" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid club type"})
		}
		filter.Type = &typeStr
	}

	clubs, err := h.db.GetClubs(filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch clubs"})
	}

	return c.JSON(clubs)
}

func (h *Handler) GetClub(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid club ID"})
	}

	club, err := h.db.GetClub(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Club not found"})
	}

	return c.JSON(club)
}