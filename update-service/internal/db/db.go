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

	// Проверка и конвертация координат
	var lat, lon float64
	var latValid, lonValid bool
	if club.Lat != "" {
		var err error
		lat, err = strconv.ParseFloat(club.Lat, 64)
		if err != nil {
			db.logger.Warn("Invalid latitude format", zap.String("lat", club.Lat), zap.String("name", club.Name))
		} else {
			latValid = true
		}
	}
	if club.Lon != "" {
		var err error
		lon, err = strconv.ParseFloat(club.Lon, 64)
		if err != nil {
			db.logger.Warn("Invalid longitude format", zap.String("lon", club.Lon), zap.String("name", club.Name))
		} else {
			lonValid = true
		}
	}

	var latPtr, lonPtr *float64
	if latValid && lonValid {
		latPtr = &lat
		lonPtr = &lon
	} else {
		db.logger.Warn("Missing or invalid coordinates, using NULL", zap.String("name", club.Name))
	}

	// Начало транзакции
	tx, err := db.conn.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Проверка существования клуба
	var clubID int
	err = tx.QueryRow(`
        SELECT id FROM clubs WHERE name = $1 AND address = $2`,
		club.Name, club.Address).Scan(&clubID)
	if err == sql.ErrNoRows {
		// Клуб не существует, вставляем новый
		err = tx.QueryRow(`
            INSERT INTO clubs (name, address, description, working_hours, rating, lat, lon, type, status)
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
            RETURNING id`,
			club.Name, club.Address, club.Description, club.WorkingHours, rating,
			latPtr, lonPtr, "commercial", club.Status).Scan(&clubID)
		if err != nil {
			return fmt.Errorf("failed to insert club: %w", err)
		}
		db.logger.Info("Inserted commercial club", zap.Int("club_id", clubID), zap.String("name", club.Name))
	} else if err != nil {
		return fmt.Errorf("failed to check club existence: %w", err)
	} else {
		// Клуб существует, обновляем данные
		_, err = tx.Exec(`
            UPDATE clubs
            SET description = $1, working_hours = $2, rating = $3, lat = $4, lon = $5, status = $6
            WHERE id = $7`,
			club.Description, club.WorkingHours, rating, latPtr, lonPtr, club.Status, clubID)
		if err != nil {
			return fmt.Errorf("failed to update club: %w", err)
		}
		db.logger.Info("Updated commercial club", zap.Int("club_id", clubID), zap.String("name", club.Name))
	}

	// Вставка категорий и подкатегорий
	if err := db.insertCategoriesAndSubcategories(tx, clubID, club.Categories); err != nil {
		return err
	}

	// Коммит транзакции
	return tx.Commit()
}

func (db *DB) InsertUniversityClub(club models.UniClub) error {
	// Проверка и конвертация координат
	var lat, lon float64
	var latValid, lonValid bool
	if club.Lat != "" {
		var err error
		lat, err = strconv.ParseFloat(club.Lat, 64)
		if err != nil {
			db.logger.Warn("Invalid latitude format", zap.String("lat", club.Lat), zap.String("name", club.Name))
		} else {
			latValid = true
		}
	}
	if club.Lon != "" {
		var err error
		lon, err = strconv.ParseFloat(club.Lon, 64)
		if err != nil {
			db.logger.Warn("Invalid longitude format", zap.String("lon", club.Lon), zap.String("name", club.Name))
		} else {
			lonValid = true
		}
	}

	var latPtr, lonPtr *float64
	if latValid && lonValid {
		latPtr = &lat
		lonPtr = &lon
	} else {
		db.logger.Warn("Missing or invalid coordinates, using NULL", zap.String("name", club.Name))
	}

	// Начало транзакции
	tx, err := db.conn.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Проверка существования клуба
	var clubID int
	err = tx.QueryRow(`
        SELECT id FROM clubs WHERE name = $1 AND address = $2`,
		club.Name, club.Address).Scan(&clubID)
	if err == sql.ErrNoRows {
		// Клуб не существует, вставляем новый
		err = tx.QueryRow(`
            INSERT INTO clubs (name, address, description, working_hours, rating, lat, lon, type, status)
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
            RETURNING id`,
			club.Name, club.Address, club.Description, club.WorkingHours, 0.0,
			latPtr, lonPtr, "university", club.Status).Scan(&clubID)
		if err != nil {
			return fmt.Errorf("failed to insert club: %w", err)
		}
		db.logger.Info("Inserted university club", zap.Int("club_id", clubID), zap.String("name", club.Name))
	} else if err != nil {
		return fmt.Errorf("failed to check club existence: %w", err)
	} else {
		// Клуб существует, обновляем данные
		_, err = tx.Exec(`
            UPDATE clubs
            SET description = $1, working_hours = $2, lat = $3, lon = $4, status = $5
            WHERE id = $6`,
			club.Description, club.WorkingHours, latPtr, lonPtr, club.Status, clubID)
		if err != nil {
			return fmt.Errorf("failed to update club: %w", err)
		}
		db.logger.Info("Updated university club", zap.Int("club_id", clubID), zap.String("name", club.Name))
	}

	// Вставка категорий и подкатегорий
	if err := db.insertCategoriesAndSubcategories(tx, clubID, club.Categories); err != nil {
		return err
	}

	// Вставка расписания
	if err := db.insertSchedule(tx, clubID, club.Schedule); err != nil {
		return err
	}

	// Коммит транзакции
	return tx.Commit()
}

func (db *DB) insertCategoriesAndSubcategories(tx *sql.Tx, clubID int, categories map[string][]string) error {
	if categories == nil || len(categories) == 0 {
		db.logger.Warn("No categories to insert for club", zap.Int("club_id", clubID))
		return nil
	}

	// Удаляем старые связи с категориями и подкатегориями
	_, err := tx.Exec(`
        DELETE FROM club_categories WHERE club_id = $1`, clubID)
	if err != nil {
		db.logger.Error("Failed to delete old club categories", zap.Int("club_id", clubID), zap.Error(err))
		return fmt.Errorf("failed to delete old club categories: %w", err)
	}
	_, err = tx.Exec(`
        DELETE FROM club_subcategories WHERE club_id = $1`, clubID)
	if err != nil {
		db.logger.Error("Failed to delete old club subcategories", zap.Int("club_id", clubID), zap.Error(err))
		return fmt.Errorf("failed to delete old club subcategories: %w", err)
	}

	for catName, subcats := range categories {
		db.logger.Info("Inserting category", zap.String("category", catName), zap.Int("club_id", clubID))
		var catID int
		err := tx.QueryRow(`
            INSERT INTO categories (name) VALUES ($1)
            ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
            RETURNING id`, catName).Scan(&catID)
		if err != nil {
			db.logger.Error("Failed to insert or update category",
				zap.String("category", catName),
				zap.Int("club_id", clubID),
				zap.Error(err))
			return fmt.Errorf("failed to insert category %s: %w", catName, err)
		}

		_, err = tx.Exec(`
            INSERT INTO club_categories (club_id, category_id)
            VALUES ($1, $2)
            ON CONFLICT DO NOTHING`, clubID, catID)
		if err != nil {
			db.logger.Error("Failed to link club with category",
				zap.String("category", catName),
				zap.Int("club_id", clubID),
				zap.Error(err))
			return fmt.Errorf("failed to link club with category %s: %w", catName, err)
		}

		for _, subcat := range subcats {
			db.logger.Info("Inserting subcategory",
				zap.String("subcategory", subcat),
				zap.String("category", catName),
				zap.Int("club_id", clubID))
			var subcatID int
			err = tx.QueryRow(`
                INSERT INTO subcategories (category_id, name) VALUES ($1, $2)
                ON CONFLICT (category_id, name) DO UPDATE SET name = EXCLUDED.name
                RETURNING id`, catID, subcat).Scan(&subcatID)
			if err != nil {
				db.logger.Error("Failed to insert or update subcategory",
					zap.String("subcategory", subcat),
					zap.String("category", catName),
					zap.Int("club_id", clubID),
					zap.Error(err))
				return fmt.Errorf("failed to insert subcategory %s: %w", subcat, err)
			}

			_, err = tx.Exec(`
                INSERT INTO club_subcategories (club_id, subcategory_id)
                VALUES ($1, $2)
                ON CONFLICT DO NOTHING`, clubID, subcatID)
			if err != nil {
				db.logger.Error("Failed to link club with subcategory",
					zap.String("subcategory", subcat),
					zap.String("category", catName),
					zap.Int("club_id", clubID),
					zap.Error(err))
				return fmt.Errorf("failed to link club with subcategory %s: %w", subcat, err)
			}
		}
	}
	return nil
}

func (db *DB) insertSchedule(tx *sql.Tx, clubID int, schedule models.Schedule) error {
	// Удаляем старое расписание
	_, err := tx.Exec(`
        DELETE FROM schedules WHERE club_id = $1`, clubID)
	if err != nil {
		return fmt.Errorf("failed to delete old schedule: %w", err)
	}

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
			_, err := tx.Exec(`
                INSERT INTO schedules (club_id, day_of_week, time, activity, instructor)
                VALUES ($1, $2, $3, $4, $5)`,
				clubID, s.day, item.Time, item.Activity, item.Instructor)
			if err != nil {
				return fmt.Errorf("failed to insert schedule item for %s: %w", s.day, err)
			}
		}
	}
	return nil
}
