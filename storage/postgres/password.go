package postgres

import (
	"bot/models"
	"database/sql"
	"errors"
	"fmt"
	"log"
)

type PasswordStorage struct {
	db *sql.DB
}

func NewPasswordStorage(db *sql.DB) *PasswordStorage {
	return &PasswordStorage{db: db}
}

func (u *PasswordStorage) CreatePassword(userID string, password *models.Password) error {
	query := `
		INSERT INTO passwords (user_id, site, password)
		VALUES ($1, $2, $3)
	`
	_, err := u.db.Exec(query, userID, password.Site, password.Password)
	if err != nil {
		return fmt.Errorf("failed to create password: %w", err)
	}
	return nil
}

func (u *PasswordStorage) GetAllPasswordsByUserID(userID string) ([]models.Password, error) {
	query := `
		SELECT site, password
		FROM passwords
		WHERE user_id = $1
	`
	rows, err := u.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var passwords []models.Password
	for rows.Next() {
		var password models.Password
		if err := rows.Scan(&password.Site, &password.Password); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		passwords = append(passwords, password)
	}

	if len(passwords) == 0 {
		return nil, errors.New("no passwords found for the given user ID")
	}

	return passwords, nil
}

// GetByName retrieves passwords for a specific user ID and site name
func (u *PasswordStorage) GetByName(userID string, site string) ([]models.Password, error) {
	query := `
		SELECT site, password
		FROM passwords
		WHERE user_id = $1 AND site ILIKE '%' || $2 || '%'
	`
	log.Printf("Executing query with user_id: %s, site: %s", userID, site)
	rows, err := u.db.Query(query, userID, site)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var passwords []models.Password
	for rows.Next() {
		var password models.Password
		if err := rows.Scan(&password.Site, &password.Password); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		passwords = append(passwords, password)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	if len(passwords) == 0 {
		return nil, errors.New("no passwords found for the given user ID and site name")
	}

	return passwords, nil
}
