package service

import (
	"bot/models"
	"bot/storage/postgres"
	"log/slog"
)

type PasswordService struct {
	product *postgres.PasswordStorage
}

func NewPasswordService(pr *postgres.PasswordStorage) *PasswordService {
	return &PasswordService{product: pr}
}

// CreatePassword adds a new password associated with the user's ID
func (p *PasswordService) CreatePassword(userID string, pr models.Password) error {
	slog.Info("CreatePassword Service", "user_id", userID, "password", pr)
	err := p.product.CreatePassword(userID, &pr)
	if err != nil {
		slog.Error("Error while creating password", "err", err)
		return err
	}
	return nil
}

// GetAllPasswordsByUserID retrieves all passwords associated with the user's ID
func (p *PasswordService) GetAllPasswordsByUserID(userID string) ([]models.Password, error) {
	slog.Info("GetAllPasswordsByUserID Service", "user_id", userID)
	passwords, err := p.product.GetAllPasswordsByUserID(userID)
	if err != nil {
		slog.Error("Error while fetching passwords by user ID", "err", err)
		return nil, err
	}

	slog.Info("Successfully fetched passwords by user ID", "passwords", passwords)
	return passwords, nil
}

// GetByName retrieves passwords associated with the user's ID and matching the site name
func (p *PasswordService) GetByName(userID string, site string) ([]models.Password, error) {
	slog.Info("GetByName Service", "user_id", userID, "site", site)
	passwords, err := p.product.GetByName(userID, site)
	if err != nil {
		slog.Error("Error while fetching passwords by name", "err", err)
		return nil, err
	}

	slog.Info("Successfully fetched passwords by name", "passwords", passwords)
	return passwords, nil
}
