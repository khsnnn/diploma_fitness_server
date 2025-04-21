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

	for _, club := range commercialClubs {
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

	for _, club := range universityClubs {
		if err := s.db.InsertUniversityClub(club); err != nil {
			s.logger.Error("Failed to insert university club", zap.String("name", club.Name), zap.Error(err))
			continue
		}
	}

	s.logger.Info("Successfully updated clubs")
	return nil
}
