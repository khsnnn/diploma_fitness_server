package db

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/khsnnn/diploma_fitness_server/update-service/internal/models"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type DB struct {
	conn   *sql.DB
	logger *zap.Logger
}

func NewDB(host, port, user, password, dbname string, logger *zap.Logger) (*DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	return &DB{conn: conn, logger: logger}, nil
}

func (db *DB) Close() error {
	return db.conn.Close()
}

func (db *DB) InsertCommercialClub(club models.Club) error {
	// Нормализация рейтинга
	rating := 0.0
	if strings.Contains(club.Rating, "/") {
		parts := strings.Split(club.Rating, "/")
		if len(parts) > 0 {
			ratingStr := strings.ReplaceAll(parts[0], ",", ".")
			var err error
			rating, err = strconv.ParseFloat(ratingStr, 64)
			if err != nil {
				db.logger.Warn("Invalid rating format", zap.String("rating", club.Rating))
				rating = 0.0
			}
		}
	}

	// Вставка клуба
	var clubID int
	err := db.conn.QueryRow(`
        INSERT INTO clubs (name, address, description, working_hours, rating, coordinates, type, status)
        VALUES ($1, $2, $3, $4, $5, POINT($6, $7), $8, $9)
        ON CONFLICT (name, address) DO NOTHING
        RETURNING id`,
		club.Name, club.Address, club.Description, club.WorkingHours, rating,
		club.Coordinates.Lng, club.Coordinates.Lat, "commercial", club.Status).Scan(&clubID)
	if err == sql.ErrNoRows {
		db.logger.Info("Club already exists", zap.String("name", club.Name))
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to insert club: %w", err)
	}
	db.logger.Info("Inserted commercial club", zap.Int("club_id", clubID), zap.String("name", club.Name))

	// Вставка категорий и подкатегорий
	return db.insertCategoriesAndSubcategories(clubID, club.Categories)
}

func (db *DB) InsertUniversityClub(club models.UniClub) error {
	// Вставка клуба
	var clubID int
	err := db.conn.QueryRow(`
        INSERT INTO clubs (name, address, description, working_hours, rating, coordinates, type, status)
        VALUES ($1, $2, $3, $4, $5, POINT($6, $7), $8, $9)
        ON CONFLICT (name, address) DO NOTHING
        RETURNING id`,
		club.Name, club.Address, club.Description, club.WorkingHours, 0.0,
		club.Coordinates.Lng, club.Coordinates.Lat, "university", club.Status).Scan(&clubID)
	if err == sql.ErrNoRows {
		db.logger.Info("Club already exists", zap.String("name", club.Name))
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to insert club: %w", err)
	}
	db.logger.Info("Inserted university club", zap.Int("club_id", clubID), zap.String("name", club.Name))

	// Вставка категорий и подкатегорий
	if err := db.insertCategoriesAndSubcategories(clubID, club.Categories); err != nil {
		return err
	}

	// Вставка расписания
	return db.insertSchedule(clubID, club.Schedule)
}

func (db *DB) insertCategoriesAndSubcategories(clubID int, categories map[string][]string) error {
	for catName, subcats := range categories {
		// Вставка категории
		var catID int
		err := db.conn.QueryRow(`
            INSERT INTO categories (name) VALUES ($1)
            ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
            RETURNING id`, catName).Scan(&catID)
		if err != nil {
			return fmt.Errorf("failed to insert category %s: %w", catName, err)
		}

		// Связывание клуба с категорией
		_, err = db.conn.Exec(`
            INSERT INTO club_categories (club_id, category_id)
            VALUES ($1, $2)
            ON CONFLICT DO NOTHING`, clubID, catID)
		if err != nil {
			return fmt.Errorf("failed to link club with category %s: %w", catName, err)
		}

		// Вставка подкатегорий
		for _, subcat := range subcats {
			var subcatID int
			err = db.conn.QueryRow(`
                INSERT INTO subcategories (category_id, name) VALUES ($1, $2)
                ON CONFLICT (category_id, name) DO UPDATE SET name = EXCLUDED.name
                RETURNING id`, catID, subcat).Scan(&subcatID)
			if err != nil {
				return fmt.Errorf("failed to insert subcategory %s: %w", subcat, err)
			}

			// Связывание клуба с подкатегорией
			_, err = db.conn.Exec(`
                INSERT INTO club_subcategories (club_id, subcategory_id)
                VALUES ($1, $2)
                ON CONFLICT DO NOTHING`, clubID, subcatID)
			if err != nil {
				return fmt.Errorf("failed to link club with subcategory %s: %w", subcat, err)
			}
		}
	}
	return nil
}

func (db *DB) insertSchedule(clubID int, schedule models.Schedule) error {
	schedules := []struct {
		day   string
		items []models.ScheduleItem
	}{
		{"Monday", schedule.Monday},
		{"Tuesday", schedule.Tuesday},
		{"Wednesday", schedule.Wednesday},
		{"Thursday", schedule.Thursday},
		{"Friday", schedule.Friday},
		{"Saturday", schedule.Saturday},
	}

	for _, s := range schedules {
		for _, item := range s.items {
			_, err := db.conn.Exec(`
                INSERT INTO schedules (club_id, day_of_week, time, activity, instructor)
                VALUES ($1, $2, $3, $4, $5)
                ON CONFLICT DO NOTHING`,
				clubID, s.day, item.Time, item.Activity, item.Instructor)
			if err != nil {
				return fmt.Errorf("failed to insert schedule item for %s: %w", s.day, err)
			}
		}
	}
	return nil
}
