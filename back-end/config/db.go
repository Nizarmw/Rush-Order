package config

import (
	"RushOrder/models"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
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
}
