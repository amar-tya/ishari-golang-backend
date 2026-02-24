package main

import (
	"fmt"
	"log"
	"time"

	"ishari-backend/internal/core/entity"
	"ishari-backend/pkg/config"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Struktur tabel MySQL
type MySQLVerse struct {
	ID              int64     `gorm:"column:id"`
	MuhudID         int64     `gorm:"column:muhud_id"`
	Position        int       `gorm:"column:position"`
	Text            string    `gorm:"column:text"`
	TranslateID     string    `gorm:"column:translate_id"`
	Transliteration string    `gorm:"column:transliteration"`
	IsDiwan         string    `gorm:"column:is_diwan"`
	IsDiba          string    `gorm:"column:is_diba"`
	CreatedAt       time.Time `gorm:"column:created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at"`
}

func (MySQLVerse) TableName() string {
	return "arabic_text"
}

// Struktur tabel Muhud MySQL
type MySQLMuhud struct {
	ID       int64  `gorm:"column:id"`
	Position int    `gorm:"column:position"`
	Name     string `gorm:"column:name"`
}

func (MySQLMuhud) TableName() string {
	return "muhud"
}

func main() {
	// 1. KONEKSI KE MYSQL
	mysqlDSN := "root:Inipassword:)^@tcp(localhost:3306)/ishari_240225?charset=utf8mb4&parseTime=True&loc=Local"
	mysqlDB, err := gorm.Open(mysql.Open(mysqlDSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("Gagal connect ke MySQL: %v", err)
	}
	fmt.Println("Berhasil connect ke MySQL!")

	// 2. KONEKSI KE POSTGRESQL MENGGUNAKAN CONFIG APP
	cfg := config.Load()
	pgDSN := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Jakarta",
		cfg.Database.Host, cfg.Database.User, cfg.Database.Password, cfg.Database.DBName, cfg.Database.Port, cfg.Database.SSLMode)

	pgDB, err := gorm.Open(postgres.Open(pgDSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("Gagal connect ke PostgreSQL: %v", err)
	}
	fmt.Println("Berhasil connect ke PostgreSQL!")

	// 3. AMBIL DATA DARI MYSQL
	var mysqlMuhuds []MySQLMuhud
	if err := mysqlDB.Find(&mysqlMuhuds).Error; err != nil {
		log.Fatalf("Gagal load data muhud dari MySQL: %v", err)
	}

	// Buat map untuk mempercepat pencarian muhud berdasarkan ID
	muhudMap := make(map[int64]MySQLMuhud)
	for _, m := range mysqlMuhuds {
		muhudMap[m.ID] = m
	}

	var mysqlVerses []MySQLVerse
	if err := mysqlDB.Find(&mysqlVerses).Error; err != nil {
		log.Fatalf("Gagal load data verse dari MySQL: %v", err)
	}
	fmt.Printf("Ditemukan %d verse di MySQL. Mulai proses migrasi...\n", len(mysqlVerses))

	// 4. LOOP & INSERT KE POSTGRESQL
	successCount := 0
	for _, mv := range mysqlVerses {
		// Dapatkan Muhud / Chapter Number
		muhud, exists := muhudMap[mv.MuhudID]
		if !exists {
			log.Printf("Peringatan: Muhud ID %d tidak ditemukan di tabel muhud MySQL, verse id %d dilewati.\n", mv.MuhudID, mv.ID)
			continue
		}

		// Tentukan Category (Diwan / Diba)
		category := "Diwan" // Default
		if mv.IsDiba == "Y" {
			category = "Diba"
		}

		// Cari ID Chapter di PostgreSQL berdasarkan chapter_number dan category
		var pgChapter entity.Chapter
		if err := pgDB.Where("chapter_number = ? AND category = ?", muhud.Position, category).First(&pgChapter).Error; err != nil {
			log.Printf("Peringatan: Chapter dgn urutan %d dan category %s tidak ditemukan di PostgreSQL. Verse id MySQL %d dilewati.\n", muhud.Position, category, mv.ID)
			continue
		}

		// Map transliteration
		transliteration := mv.Transliteration

		pgVerse := entity.Verse{
			ChapterID:       pgChapter.ID,
			VerseNumber:     uint(mv.Position),
			ArabicText:      mv.Text,
			Transliteration: &transliteration,
			CreatedAt:       mv.CreatedAt,
			UpdatedAt:       mv.UpdatedAt,
		}

		// Insert row ke table `verses` PostgreSQL
		if err := pgDB.Create(&pgVerse).Error; err != nil {
			log.Printf("Gagal insert verse (chapter_id pg: %d, verse_number: %d). Error: %v\n", pgChapter.ID, pgVerse.VerseNumber, err)
			continue
		}
		successCount++
	}

	fmt.Printf("Proses Migrasi Selesai! Berhasil migrasi %d dari %d verse.\n", successCount, len(mysqlVerses))
}
