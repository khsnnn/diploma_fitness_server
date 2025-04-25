package service

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	"github.com/khsnnn/diploma_fitness_server/update-service/internal/db"
	"github.com/khsnnn/diploma_fitness_server/update-service/internal/models"
	"go.uber.org/zap"
)

type Service struct {
	db     *db.DB
	logger *zap.Logger
}

func NewService(db *db.DB, logger *zap.Logger) *Service {
	return &Service{db: db, logger: logger}
}

func (s *Service) UpdateClubs(dataDir string) error {
	// Обработка коммерческих клубов
	commercialPath := filepath.Join(dataDir, "commercial_clubs.json")
	commercialData, err := ioutil.ReadFile(commercialPath)
	if err != nil {
		s.logger.Error("Failed to read commercial clubs file", zap.Error(err))
		return err
	}

	var commercialClubs []models.Club
	if err := json.Unmarshal(commercialData, &commercialClubs); err != nil {
		s.logger.Error("Failed to unmarshal commercial clubs", zap.Error(err))
		return err
	}
	s.logger.Info("Read commercial clubs", zap.Int("count", len(commercialClubs)))
	for _, club := range commercialClubs {
		s.logger.Info("Processing commercial club",
			zap.String("name", club.Name),
			zap.Any("categories", club.Categories))
		if club.Categories == nil {
			s.logger.Warn("Categories is nil for commercial club", zap.String("name", club.Name))
			club.Categories = make(map[string][]string) // Инициализируем пустой словарь
		} else if len(club.Categories) == 0 {
			s.logger.Info("Categories is empty for commercial club", zap.String("name", club.Name))
		} else {
			s.logger.Info("Categories found for commercial club",
				zap.String("name", club.Name),
				zap.Int("category_count", len(club.Categories)))
		}
		if err := s.db.InsertCommercialClub(club); err != nil {
			s.logger.Error("Failed to insert commercial club", zap.String("name", club.Name), zap.Error(err))
			continue
		}
	}

	// Обработка университетских клубов
	universityPath := filepath.Join(dataDir, "university_clubs.json")
	universityData, err := ioutil.ReadFile(universityPath)
	if err != nil {
		s.logger.Error("Failed to read university clubs file", zap.Error(err))
		return err
	}

	var universityClubs []models.UniClub
	if err := json.Unmarshal(universityData, &universityClubs); err != nil {
		s.logger.Error("Failed to unmarshal university clubs", zap.Error(err))
		return err
	}
	s.logger.Info("Read university clubs", zap.Int("count", len(universityClubs)))
	for _, club := range universityClubs {
		s.logger.Info("Processing university club",
			zap.String("name", club.Name),
			zap.Any("categories", club.Categories))
		if club.Categories == nil {
			s.logger.Warn("Categories is nil for university club", zap.String("name", club.Name))
			club.Categories = make(map[string][]string) // Инициализируем пустой словарь
		} else if len(club.Categories) == 0 {
			s.logger.Info("Categories is empty for university club", zap.String("name", club.Name))
		} else {
			s.logger.Info("Categories found for university club",
				zap.String("name", club.Name),
				zap.Int("category_count", len(club.Categories)))
		}
		if err := s.db.InsertUniversityClub(club); err != nil {
			s.logger.Error("Failed to insert university club", zap.String("name", club.Name), zap.Error(err))
			continue
		}
	}

	s.logger.Info("Successfully updated clubs")
	return nil
}