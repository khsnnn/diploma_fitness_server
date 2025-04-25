package db

import (
	"github.com/khsnnn/diploma_fitness_server/club-service/internal/models"
	"github.com/khsnnn/diploma_fitness_server/club-service/internal/utils"
	"gorm.io/gorm"
)

type DB struct {
	db *gorm.DB
}

func NewDB(db *gorm.DB) *DB {
	return &DB{db: db}
}

type ClubFilter struct {
	Lat           *float64
	Lon           *float64
	Distance      float64
	MinRating     *float64
	Categories    []string
	Subcategories []string
	Type          *string
}

func (d *DB) GetClubs(filter ClubFilter) ([]models.APIClub, error) {
	query := d.db.Model(&models.DBClub{}).
		Preload("Categories.Subcategories").
		Preload("Schedules")

	if filter.Type != nil {
		query = query.Where("type = ?", *filter.Type)
	}

	if filter.MinRating != nil {
		query = query.Where("rating >= ?", *filter.MinRating)
	}

	if len(filter.Categories) > 0 {
		query = query.Joins("JOIN club_categories cc ON clubs.id = cc.club_id").
			Joins("JOIN categories c ON cc.category_id = c.id").
			Where("c.name IN ?", filter.Categories).
			Group("clubs.id")
	}

	if len(filter.Subcategories) > 0 {
		query = query.Joins("JOIN club_subcategories cs ON clubs.id = cs.club_id").
			Joins("JOIN subcategories s ON cs.subcategory_id = s.id").
			Where("s.name IN ?", filter.Subcategories).
			Group("clubs.id")
	}

	var dbClubs []models.DBClub
	if err := query.Find(&dbClubs).Error; err != nil {
		return nil, err
	}

	if filter.Lat != nil && filter.Lon != nil {
		filteredClubs := make([]models.DBClub, 0)
		for _, club := range dbClubs {
			distance := utils.HaversineDistance(*filter.Lat, *filter.Lon, club.Lat, club.Lon)
			if distance <= filter.Distance {
				filteredClubs = append(filteredClubs, club)
			}
		}
		dbClubs = filteredClubs
	}

	apiClubs := make([]models.APIClub, 0, len(dbClubs))
	for _, dbClub := range dbClubs {
		apiCategories := make([]models.APICategory, 0, len(dbClub.Categories))
		for _, cat := range dbClub.Categories {
			subcatNames := make([]string, 0, len(cat.Subcategories))
			for _, subcat := range cat.Subcategories {
				subcatNames = append(subcatNames, subcat.Name)
			}
			apiCategories = append(apiCategories, models.APICategory{
				Name:          cat.Name,
				Subcategories: subcatNames,
			})
		}

		apiClubs = append(apiClubs, models.APIClub{
			ID:           dbClub.ID,
			Name:         dbClub.Name,
			Address:      dbClub.Address,
			Description:  dbClub.Description,
			WorkingHours: dbClub.WorkingHours,
			Rating:       dbClub.Rating,
			Lat:          dbClub.Lat,
			Lon:          dbClub.Lon,
			Type:         dbClub.Type,
			Status:       dbClub.Status,
			Categories:   apiCategories,
			Schedules:    dbClub.Schedules,
		})
	}

	return apiClubs, nil
}
