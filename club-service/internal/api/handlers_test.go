package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/khsnnn/diploma_fitness_server/club-service/internal/models"
)

func TestGetClubs(t *testing.T) {
	app := fiber.New()

	app.Get("/clubs", GetClubs)

	req := httptest.NewRequest(http.MethodGet, "/clubs", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
	}

	var clubs []models.Club
	if err := json.NewDecoder(resp.Body).Decode(&clubs); err != nil {
		t.Fatalf("Failed to decode responce: %v", err)
	}

	if len(clubs) != 0 {
		t.Errorf("Expected empty array, got %d clubs", len(clubs))
	}
}
