package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/khsnnn/diploma_fitness_server/club-service/internal/api"
	"github.com/khsnnn/diploma_fitness_server/club-service/internal/db"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=localhost user=postgres password=your_password dbname=fitness_db port=5432 sslmode=disable"
	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	dbInstance := db.NewDB(gormDB)

	handler := api.NewHandler(dbInstance)

	app := fiber.New()

	app.Get("/clubs", handler.GetClubs)

	if err := app.Listen(":8080"); err != nil {
		panic("failed to start server: " + err.Error())
	}
}
