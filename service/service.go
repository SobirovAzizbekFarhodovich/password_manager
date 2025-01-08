package service

import (
	"database/sql"
	"bot/storage/postgres"
)

type Service struct {
	PrService *PasswordService
}

func InitServices(db *sql.DB) *Service {
	product := NewPasswordService(postgres.NewPasswordStorage(db))

	return &Service{product}
}
