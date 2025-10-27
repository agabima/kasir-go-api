package models

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Nama      string    `json:"nama"`
	Username  string    `gorm:"unique" json:"username"`
	Password  string    `json:"-"`
	Role      string    `json:"role"` // "admin" atau "kasir"
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
