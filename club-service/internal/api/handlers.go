package api

import (
    "github.com/gofiber/fiber/v2"
    "github.com/khsnnn/diploma_fitness_server/club-service/internal/models"
    "math"
    "strconv"
    "strings"
)

// haversineDistance в км между двумя точками
func haversineDistance(lat1, lng1, lat2, lng2 float64) float64 {
    const R = 6371 // Радиус Земли в км
    dLat := (lat2 - lat1) * math.Pi / 180
    dLng := (lng2 - lng1) * math.Pi / 180
    a := math.Sin(dLat/2)*math.Sin(dLat/2) +
        math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
            math.Sin(dLng/2)*math.Sin(dLng/2)
    c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
    return R * c
}

func GetClubs(c *fiber.Ctx) error {
    allClubs := []models.Club{
        {ID: 1, Name: "Клуб 1", Rating: 2.5, Coordinates: models.Coordinates{Lat: 57.152, Lng: 65.534}, Clusters: []string{"Фитнес и общая физическая активность"}},
        {ID: 2, Name: "Клуб 2", Rating: 3.5, Coordinates: models.Coordinates{Lat: 57.153, Lng: 65.535}, Clusters: []string{"Танцевальные направления"}},
        {ID: 3, Name: "Клуб 3", Rating: 4.0, Coordinates: models.Coordinates{Lat: 57.160, Lng: 65.540}, Clusters: []string{"Йога и духовные практики", "Танцевальные направления"}},
    }

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

    var filteredClubs []models.Club
    for _, club := range allClubs {
        // Рейтинг
        if club.Rating < minRating {
            continue
        }

        // Расстояние (если координаты пользователя указаны)
        if latStr != "" && lngStr != "" {
            distance := haversineDistance(userLat, userLng, club.Coordinates.Lat, club.Coordinates.Lng)
            if distance > maxDistance {
                continue
            }
        }

        // Кластеры
        if len(wantedClusters) > 0 {
            matches := false
            for _, wanted := range wantedClusters {
                for _, clubCluster := range club.Clusters {
                    if strings.TrimSpace(wanted) == clubCluster {
                        matches = true
                        break
                    }
                }
                if matches {
                    break
                }
            }
            if !matches {
                continue
            }
        }

        filteredClubs = append(filteredClubs, club)
    }

    return c.JSON(filteredClubs)
}