package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/khsnnn/diploma_fitness_server/club-service/internal/api"
	"github.com/khsnnn/diploma_fitness_server/club-service/internal/repository"
)

func main() {
	if err := godotenv.Load(); err != nil { 
		log.Fatalf("Failed to load .env file: %v", err)
	}

	db, err := repository.NewDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	api.DB = db

	app := fiber.New()
	app.Get("/api/clubs", api.GetClubs)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080" 
	}

	log.Printf("Server starting on port %s", port)
	if err := app.Listen(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
