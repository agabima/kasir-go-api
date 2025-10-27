package models

import "time"

type Produk struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Nama        string    `json:"nama"`
	Harga       float64   `json:"harga"`
	Stok        int       `json:"stok"`
	KategoriID  uint      `json:"kategori_id"`
	Kategori    Kategori  `gorm:"foreignKey:KategoriID" json:"kategori"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
