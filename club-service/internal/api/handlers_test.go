package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/khsnnn/diploma_fitness_server/club-service/internal/models"
)

func TestGetClubsWithDistanceAndClusters(t *testing.T) {
    app := fiber.New()
    app.Get("/api/clubs", GetClubs)

    req := httptest.NewRequest(http.MethodGet, "/api/clubs?distance=5&clusters=Танцевальные%20направления&lat=57.152&lng=65.534", nil)
    resp, err := app.Test(req)
    if err != nil {
        t.Fatalf("Failed to perform request: %v", err)
    }

    if resp.StatusCode != http.StatusOK {
        t.Errorf("Expected status 200 OK, got %d", resp.StatusCode)
    }

    var clubs []models.Club
    if err := json.NewDecoder(resp.Body).Decode(&clubs); err != nil {
        t.Fatalf("Failed to decode response: %v", err)
    }
}
