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

func (p *PasswordService) CreatePassword(pr models.Password) error {
	slog.Info("CreatePassword Service", "password", pr)
	err := p.product.CreatePassword(&pr)
	if err != nil {
		slog.Error("Error while creating password", "err", err)
		return err
	}
	return nil
}

func (p *PasswordService) GetAllPasswordsByPhone(phone string) ([]models.Password, error) {
	slog.Info("GetAllPasswordsByPhone Service", "phone", phone)
	passwords, err := p.product.GetAllPasswordsByPhone(phone)
	if err != nil {
		slog.Error("Error while fetching passwords by phone", "err", err)
		return nil, err
	}

	slog.Info("Successfully fetched passwords by phone", "passwords", passwords)
	return passwords, nil
}

func (p *PasswordService) GetByName(phone string, site string) ([]models.Password, error) {
	slog.Info("GetByName Service", "phone", phone, "site", site)
	passwords, err := p.product.GetByName(phone, site)
	if err != nil {
		slog.Error("Error while fetching passwords by name", "err", err)
		return nil, err
	}

	slog.Info("Successfully fetched passwords by name", "passwords", passwords)
	return passwords, nil
}
