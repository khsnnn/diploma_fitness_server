package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/khsnnn/diploma_fitness_server/club-service/internal/api"
)

func main(){
	app := fiber.New()
	app.Get("/clubs", api.GetClubs)
	app.Listen(":8000")
}