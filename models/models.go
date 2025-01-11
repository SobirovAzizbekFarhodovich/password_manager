package models

type Password struct {
	ID string `json:"id" db:"id"`
	Site string `json:"site" db:"site"`
	Password string `json:"password" db:"password"`
}