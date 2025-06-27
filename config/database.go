package config

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/Zauro25/Capstone-PerpusKominfosan/models"
)

var DB *gorm.DB

func InitDB() {
	var err error
	
	// Database connection string
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", "password"),
		getEnv("DB_NAME", "perpustakaan_db"),
		getEnv("DB_PORT", "5432"),
	)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate tables
	err = DB.AutoMigrate(
		&models.AdminPerpustakaan{},
		&models.AdminDPK{},
		&models.Executive{},
		&models.Perpustakaan{},
		&models.Koleksi{},
		&models.SDM{},
		&models.Pengunjung{},
		&models.Anggota{},
		&models.Verifikasi{},
		&models.Laporan{},
		&models.LogRevisi{},
		&models.Notifikasi{},
		&models.AuditLog{},
	)

	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Database connected and migrated successfully")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}