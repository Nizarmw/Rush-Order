package config

import (
	"RushOrder/models"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPass, dbHost, dbPort, dbName)

	var db *gorm.DB
	for i := 0; i < 10; i++ {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		log.Printf("Retrying database connection (%d/10)...", i+1)
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		log.Fatalf("Failed to connect to database after retries: %v", err)
	}

	DB = db
	fmt.Println("Database connected!")

	fmt.Println("Migrating Pemesan...")
	if err := db.AutoMigrate(&models.Pemesan{}); err != nil {
		log.Fatal("❌ Gagal migrasi Pemesan:", err)
	}

	fmt.Println("Migrating Produk...")
	if err := db.AutoMigrate(&models.Produk{}); err != nil {
		log.Fatal("❌ Gagal migrasi Produk:", err)
	}

	fmt.Println("Migrating Order...")
	if err := db.AutoMigrate(&models.Order{}); err != nil {
		log.Fatal("❌ Gagal migrasi Order:", err)
	}

	fmt.Println("Migrating OrderItem...")
	if err := db.AutoMigrate(&models.OrderItem{}); err != nil {
		log.Fatal("❌ Gagal migrasi OrderItem:", err)
	}

	fmt.Println("Migrating Pegawai...")
	if err := db.AutoMigrate(&models.Pegawai{}); err != nil {
		log.Fatal("❌ Gagal migrasi Pegawai:", err)
	}

	fmt.Println("Migrating Payment...")
	if err := db.AutoMigrate(&models.Payment{}); err != nil {
		log.Fatal("❌ Gagal migrasi Payment:", err)
	}

	fmt.Println("Database migrated!")

	// Seed data
	seedData()
}

func seedData() {
	fmt.Println("Starting database seeding...")

	// Seed Pegawai data
	seedPegawai()

	// Seed Produk data
	seedProduk()

	fmt.Println("Database seeding completed!")
}

func seedPegawai() {
	fmt.Println("Seeding Pegawai data...")

	// Check if pegawai already exists
	var count int64
	DB.Model(&models.Pegawai{}).Count(&count)
	if count > 0 {
		fmt.Println("Pegawai data already exists, skipping...")
		return
	}

	// Create admin pegawai
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return
	}

	pegawai := models.Pegawai{
		Username: "admin",
		Password: string(hashedPassword),
	}

	if err := DB.Create(&pegawai); err.Error != nil {
		log.Printf("Error seeding pegawai: %v", err.Error)
	} else {
		fmt.Printf("✅ Created pegawai: %s\n", pegawai.Username)
	}
}

func seedProduk() {
	fmt.Println("Seeding Produk data...")

	// Check if produk already exists
	var count int64
	DB.Model(&models.Produk{}).Count(&count)
	if count > 0 {
		fmt.Println("Produk data already exists, skipping...")
		return
	}

	// Create sample products for each category
	products := []models.Produk{
		{
			IDProduk:    "MKN001",
			NamaProduk:  "Nasi Goreng Special",
			Deskripsi:   "Nasi goreng dengan telur, ayam, dan sayuran segar",
			HargaProduk: 25000,
			ImageURL:    "https://emslvefeidpmppzjxwfl.supabase.co/storage/v1/object/public/rushorder/Nasgor.jpeg",
			Kategori:    models.KategoriMakanan,
		},
		{
			IDProduk:    "MNM001",
			NamaProduk:  "Es Teh Manis",
			Deskripsi:   "Teh manis segar dengan es batu",
			HargaProduk: 8000,
			ImageURL:    "https://emslvefeidpmppzjxwfl.supabase.co/storage/v1/object/public/rushorder/Tehes.jpeg",
			Kategori:    models.KategoriMinuman,
		},
		{
			IDProduk:    "SNK001",
			NamaProduk:  "Keripik Kentang",
			Deskripsi:   "Keripik kentang renyah dengan berbagai rasa",
			HargaProduk: 12000,
			ImageURL:    "https://emslvefeidpmppzjxwfl.supabase.co/storage/v1/object/public/rushorder/Fries.jpeg",
			Kategori:    models.KategoriSnack,
		},
	}

	for _, product := range products {
		if err := DB.Create(&product); err.Error != nil {
			log.Printf("Error seeding produk %s: %v", product.IDProduk, err.Error)
		} else {
			fmt.Printf("✅ Created produk: %s - %s (%s)\n", product.IDProduk, product.NamaProduk, product.Kategori)
		}
	}
}
