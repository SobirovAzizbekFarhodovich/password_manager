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

func (u *PasswordStorage) CreatePassword(Password *models.Password) error {
	query := `
		INSERT INTO passwords (phone,site, password)
		VALUES ($1, $2, $3)
	`
	_, err := u.db.Exec(query, Password.Phone, Password.Site, Password.Password)
	return err
}

func (u *PasswordStorage) GetAllPasswordsByPhone(phone string) ([]models.Password, error) {
	query := `
		SELECT site, password
		FROM passwords
		WHERE phone = $1
	`
	rows, err := u.db.Query(query, phone)
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

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	if len(passwords) == 0 {
		return nil, errors.New("no passwords found for the given phone number")
	}

	return passwords, nil
}

func (u *PasswordStorage) GetByName(phone string, site string) ([]models.Password, error) {
	query := `
		SELECT site, password
		FROM passwords
		WHERE phone = $1 AND site ILIKE '%' || $2 || '%'
	`
	log.Printf("Executing query with phone: %s, site: %s", phone, site)
	rows, err := u.db.Query(query, phone, site)
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

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	if len(passwords) == 0 {
		return nil, errors.New("no passwords found for the given phone number and site name")
	}

	return passwords, nil
}