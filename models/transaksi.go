package models

import "time"

type Transaksi struct {
	ID              uint              `gorm:"primaryKey" json:"id"`
	Kasir           string            `json:"kasir"`
	Total           float64           `json:"total"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
	DetailTransaksi []DetailTransaksi `json:"detail_transaksi"`
}
