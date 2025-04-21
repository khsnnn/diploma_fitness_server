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
	// Инициализация логгера
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	// Загрузка переменных окружения
	if err := godotenv.Load(); err != nil {
		logger.Warn("No .env file found")
	}

	// Получение конфигурации базы данных
	dbConfig := struct {
		host, port, user, password, dbname string
	}{
		host:     os.Getenv("POSTGRES_HOST"),
		port:     os.Getenv("POSTGRES_PORT"),
		user:     os.Getenv("POSTGRES_USER"),
		password: os.Getenv("POSTGRES_PASSWORD"),
		dbname:   os.Getenv("POSTGRES_DB"),
	}

	// Подключение к базе данных
	db, err := db.NewDB(dbConfig.host, dbConfig.port, dbConfig.user, dbConfig.password, dbConfig.dbname, logger)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	// Инициализация сервиса
	svc := service.NewService(db, logger)

	// Обновление данных
	dataDir := "./data"
	if err := svc.UpdateClubs(dataDir); err != nil {
		logger.Fatal("Failed to update clubs", zap.Error(err))
	}
}
