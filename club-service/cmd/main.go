package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/khsnnn/diploma_fitness_server/club-service/internal/api"
	"github.com/khsnnn/diploma_fitness_server/club-service/internal/db"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=fitness_user password=fitness_password dbname=fitness_db port=5434 sslmode=disable"
	}

	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	dbInstance := db.NewDB(gormDB)

	handler := api.NewHandler(dbInstance)

	app := fiber.New()

	// Настройка CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3001",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	log.Println("Registering routes...")
	app.Get("/clubs", handler.GetClubs)
	app.Get("/clubs/:id", handler.GetClub)
	log.Println("Routes registered")

	if err := app.Listen(":8080"); err != nil {
		panic("failed to start server: " + err.Error())
	}
}
