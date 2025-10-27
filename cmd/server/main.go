package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"kasir-go-api/config"
	"kasir-go-api/models"
	"kasir-go-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"github.com/gin-contrib/cors"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Koneksi ke Database
	config.ConnectDatabase()

	// Auto migrate tabel Kategori
	config.DB.AutoMigrate(
		&models.Kategori{}, 
		&models.Produk{},
		&models.Transaksi{},
		&models.DetailTransaksi{},
		&models.User{},)

	router := gin.Default()

	router.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:5500"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
    }))

	// =====================
	// CRUD KATEGORI
	// =====================

	// GET /kategori -> semua data
	router.GET("/kategori", func(c *gin.Context) {
		var kategori []models.Kategori
		result := config.DB.Find(&kategori)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}
		c.JSON(http.StatusOK, kategori)
	})

	// GET /kategori/:id -> ambil 1 data
	router.GET("/kategori/:id", func(c *gin.Context) {
		id := c.Param("id")
		var kategori models.Kategori

		if err := config.DB.First(&kategori, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Kategori tidak ditemukan"})
			return
		}
		c.JSON(http.StatusOK, kategori)
	})

	// POST /kategori -> tambah baru
	router.POST("/kategori", func(c *gin.Context) {
		var input models.Kategori
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		config.DB.Create(&input)
		c.JSON(http.StatusOK, input)
	})

	// PUT /kategori/:id -> update
	router.PUT("/kategori/:id", func(c *gin.Context) {
		id := c.Param("id")
		var kategori models.Kategori

		if err := config.DB.First(&kategori, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Kategori tidak ditemukan"})
			return
		}

		var input models.Kategori
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		kategori.Nama = input.Nama
		config.DB.Save(&kategori)

		c.JSON(http.StatusOK, kategori)
	})

	// DELETE /kategori/:id -> hapus
	router.DELETE("/kategori/:id", func(c *gin.Context) {
		id := c.Param("id")
		var kategori models.Kategori

		if err := config.DB.First(&kategori, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Kategori tidak ditemukan"})
			return
		}

		config.DB.Delete(&kategori)
		c.JSON(http.StatusOK, gin.H{"message": "Kategori berhasil dihapus"})
	})

	// =======================================================
	// ================ CRUD PRODUK ==========================
	// =======================================================
	router.GET("/produk", func(c *gin.Context) {
		var produk []models.Produk
		if err := config.DB.Preload("Kategori").Find(&produk).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, produk)
	})

	router.GET("/produk/:id", func(c *gin.Context) {
		id := c.Param("id")
		var produk models.Produk
		if err := config.DB.Preload("Kategori").First(&produk, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Produk tidak ditemukan"})
			return
		}
		c.JSON(http.StatusOK, produk)
	})

	router.POST("/produk", func(c *gin.Context) {
		var input models.Produk
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Pastikan kategori ada
		var kategori models.Kategori
		if err := config.DB.First(&kategori, input.KategoriID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Kategori tidak ditemukan"})
			return
		}

		config.DB.Create(&input)
		c.JSON(http.StatusOK, input)
	})

	router.PUT("/produk/:id", func(c *gin.Context) {
		id := c.Param("id")
		var produk models.Produk
		if err := config.DB.First(&produk, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Produk tidak ditemukan"})
			return
		}

		var input models.Produk
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		produk.Nama = input.Nama
		produk.Harga = input.Harga
		produk.Stok = input.Stok
		produk.KategoriID = input.KategoriID
		config.DB.Save(&produk)

		c.JSON(http.StatusOK, produk)
	})

	router.DELETE("/produk/:id", func(c *gin.Context) {
		id := c.Param("id")
		var produk models.Produk
		if err := config.DB.First(&produk, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Produk tidak ditemukan"})
			return
		}
		config.DB.Delete(&produk)
		c.JSON(http.StatusOK, gin.H{"message": "Produk berhasil dihapus"})
	})

	// =======================================================
	// ================ AUTH (LOGIN / REGISTER) ===============
	// =======================================================
	router.POST("/register", func(c *gin.Context) {
		var input struct {
			Nama     string `json:"nama"`
			Username string `json:"username"`
			Password string `json:"password"`
			Role     string `json:"role"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		hash, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

		user := models.User{
			Nama:     input.Nama,
			Username: input.Username,
			Password: string(hash),
			Role:     input.Role,
		}
		config.DB.Create(&user)
		c.JSON(http.StatusOK, user)
	})

	router.POST("/login", func(c *gin.Context) {
		var input struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var user models.User
		if err := config.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User tidak ditemukan"})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Password salah"})
			return
		}

		token, _ := utils.GenerateToken(user.ID, user.Role)
		c.JSON(http.StatusOK, gin.H{
			"token": token,
			"user":  user,
		})
	})


	// =======================================================
	// ================ TRANSAKSI ============================
	// =======================================================

	router.GET("/transaksi", func(c *gin.Context) {
		var transaksi []models.Transaksi
		config.DB.Preload("DetailTransaksi.Produk").Find(&transaksi)
		c.JSON(http.StatusOK, transaksi)
	})

	router.POST("/transaksi", func(c *gin.Context) {
		var input struct {
			Kasir   string `json:"kasir"`
			Details []struct {
				ProdukID uint `json:"produk_id"`
				Qty      int  `json:"qty"`
			} `json:"details"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var total float64
		var detailList []models.DetailTransaksi

		for _, d := range input.Details {
			var produk models.Produk
			if err := config.DB.First(&produk, d.ProdukID).Error; err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Produk ID %d tidak ditemukan", d.ProdukID)})
				return
			}

			subtotal := float64(d.Qty) * produk.Harga
			total += subtotal

			detail := models.DetailTransaksi{
				ProdukID: d.ProdukID,
				Qty:      d.Qty,
				Subtotal: subtotal,
			}
			detailList = append(detailList, detail)

			// Kurangi stok produk
			produk.Stok -= d.Qty
			config.DB.Save(&produk)
		}

		transaksi := models.Transaksi{
			Kasir:           input.Kasir,
			Total:           total,
			DetailTransaksi: detailList,
		}

		config.DB.Create(&transaksi)
		// Kirim ke Kafka
		go utils.SendKafkaMessage("transaksi_log", transaksi)

		c.JSON(http.StatusOK, transaksi)
	})

	router.GET("/transaksi/:id", func(c *gin.Context) {
		id := c.Param("id")
		var transaksi models.Transaksi
		if err := config.DB.Preload("DetailTransaksi.Produk").First(&transaksi, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Transaksi tidak ditemukan"})
			return
		}
		c.JSON(http.StatusOK, transaksi)
	})

	router.GET("/transaksi/:id/struk", func(c *gin.Context) {
		id := c.Param("id")
		var transaksi models.Transaksi
		if err := config.DB.Preload("DetailTransaksi.Produk").First(&transaksi, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Transaksi tidak ditemukan"})
			return
		}

		struk := fmt.Sprintf("===== STRUK TRANSAKSI =====\nKasir: %s\nTanggal: %s\n\n",
			transaksi.Kasir, transaksi.CreatedAt.Format("02-01-2006 15:04"))

		for _, d := range transaksi.DetailTransaksi {
			struk += fmt.Sprintf("%s x%d = Rp%.0f\n", d.Produk.Nama, d.Qty, d.Subtotal)
		}
		struk += fmt.Sprintf("\nTOTAL: Rp%.0f\n=============================", transaksi.Total)

		c.String(http.StatusOK, struk)
	})


	// TEST /ping
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	fmt.Println("🚀 Server running on port:", port)
	router.Run(":" + port)
}
