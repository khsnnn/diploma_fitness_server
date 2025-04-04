package api

import (
	"testing"
	"net/http"
	"net/http/httptest"

	"github.com/gofiber/fiber/v2"
)

func TestGetClubs(t *testing.T) {
	app := fiber.New()

	app.Get("/clubs", GetClubs)

	req := httptest.NewRequest("GET", "/clubs", nil)
	resp, err := app.Test(req)

	if err != nil{
		t.Fatalf("Failed to perform request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	contertType := resp.Header.Get("Content-Type")
	if contertType != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got '%s'", contertType)
	}
	
}