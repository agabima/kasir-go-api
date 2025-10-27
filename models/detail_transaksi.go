package models

type DetailTransaksi struct {
	ID          uint     `gorm:"primaryKey" json:"id"`
	TransaksiID uint     `json:"transaksi_id"`
	ProdukID    uint     `json:"produk_id"`
	Qty         int      `json:"qty"`
	Subtotal    float64  `json:"subtotal"`
	Produk      Produk   `gorm:"foreignKey:ProdukID" json:"produk"`
}
