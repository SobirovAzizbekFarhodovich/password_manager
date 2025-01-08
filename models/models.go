package models

type Password struct {
	ID int64 `json:"id" db:"id"`
	Phone string `json:"phone" db:"phone"`
	Site string `json:"site" db:"site"`
	Password string `json:"password" db:"password"`
}