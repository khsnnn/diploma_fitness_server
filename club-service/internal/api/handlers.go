package api

import (
	"math"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/khsnnn/diploma_fitness_server/club-service/internal/models"
	"gorm.io/gorm"
)

var DB *gorm.DB

func haversineDistance(lat1, lng1, lat2, lng2 float64) float64 {
	const R = 6371
	dLat := (lat2 - lat1) * math.Pi / 180
	dLng := (lng2 - lng1) * math.Pi / 180
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLng/2)*math.Sin(dLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}

func GetClubs(c *fiber.Ctx) error {
	ratingStr := c.Query("rating", "0.0")
	minRating, err := strconv.ParseFloat(ratingStr, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid rating value"})
	}

	distanceStr := c.Query("distance", "1000")
	maxDistance, err := strconv.ParseFloat(distanceStr, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid distance value"})
	}

	latStr := c.Query("lat", "")
	lngStr := c.Query("lng", "")
	var userLat, userLng float64
	if latStr != "" && lngStr != "" {
		userLat, err = strconv.ParseFloat(latStr, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid lat value"})
		}
		userLng, err = strconv.ParseFloat(lngStr, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid lng value"})
		}
	}

	clustersStr := c.Query("clusters", "")
	wantedClusters := strings.Split(clustersStr, ",")
	if clustersStr == "" {
		wantedClusters = []string{}
	}

	var clubs []models.Club
	query := DB.Preload("Clusters").Where("rating >= ?", minRating)

	if len(wantedClusters) > 0 {
		query = query.Joins("JOIN club_clusters cc ON cc.club_id = clubs.id").
			Joins("JOIN clusters c ON c.id = cc.cluster_id").
			Where("c.name IN ?", wantedClusters).
			Group("clubs.id")
	}

	if err := query.Find(&clubs).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}

	if latStr != "" && lngStr != "" {
		var filteredClubs []models.Club
		for _, club := range clubs {
			distance := haversineDistance(userLat, userLng, club.Coordinates.Lat, club.Coordinates.Lng)
			if distance <= maxDistance {
				filteredClubs = append(filteredClubs, club)
			}
		}
		clubs = filteredClubs
	}

	return c.JSON(clubs)
}
