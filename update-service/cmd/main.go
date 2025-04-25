package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/khsnnn/diploma_fitness_server/update-service/internal/db"
	"github.com/khsnnn/diploma_fitness_server/update-service/internal/service"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	if err := godotenv.Load(); err != nil {
		logger.Warn("No .env file found")
	}

	dbConfig := struct {
		host, port, user, password, dbname string
	}{
		host:     os.Getenv("DB_HOST"),
		port:     os.Getenv("DB_PORT"),
		user:     os.Getenv("DB_USER"),
		password: os.Getenv("DB_PASSWORD"),
		dbname:   os.Getenv("DB_NAME"),
	}

	db, err := db.NewDB(dbConfig.host, dbConfig.port, dbConfig.user, dbConfig.password, dbConfig.dbname, logger)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	svc := service.NewService(db, logger)

	dataDir := "./data"
	if err := svc.UpdateClubs(dataDir); err != nil {
		logger.Fatal("Failed to update clubs", zap.Error(err))
	}
}
