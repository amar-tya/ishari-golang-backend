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
type MySQLShalawat struct {
	ID                         int64     `gorm:"column:id"`
	MuhudID                    int64     `gorm:"column:muhud_id"`
	TextShalawat               string    `gorm:"column:text_shalawat"`
	Transliteration            string    `gorm:"column:transliteration"`
	TranslationID              string    `gorm:"column:translation_id"`
	NumberOfDiwan              *int      `gorm:"column:numberOfDiwan"`              // Bisa dapet Null
	NumberOfMaulidSyarafulAnam *int      `gorm:"column:numberOfMaulidSyarafulAnam"` // Bisa dapet Null
	CreatedAt                  time.Time `gorm:"column:created_at"`
	UpdatedAt                  time.Time `gorm:"column:updated_at"`
}

func (MySQLShalawat) TableName() string {
	return "shalawats"
}

// Struktur tabel Muhud MySQL
type MySQLMuhud struct {
	ID                int64  `gorm:"column:id"`
	TransliterationID string `gorm:"column:transliteration_id"`
	TranslationID     string `gorm:"column:translation_id"`
}

func (MySQLMuhud) TableName() string {
	return "muhuds"
}

func main() {
	// 1. KONEKSI KE MYSQL
	mysqlDSN := "root:Inipassword:)^@tcp(localhost:3306)/ishari?charset=utf8mb4&parseTime=True&loc=Local"
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

	// 3. AMBIL DATA DARI MYSQL MUHUD (Filter 'id BETWEEN 7 AND 14')
	var mysqlMuhuds []MySQLMuhud
	if err := mysqlDB.Where("id BETWEEN ? AND ?", 7, 14).Find(&mysqlMuhuds).Error; err != nil {
		log.Fatalf("Gagal load data muhud dari MySQL: %v", err)
	}
	fmt.Printf("Ditemukan %d muhud untuk id 7 - 14.\n", len(mysqlMuhuds))

	if len(mysqlMuhuds) == 0 {
		fmt.Println("Tidak ada data muhud yang perlu diproses.")
		return
	}

	successCount := 0
	translationSuccessCount := 0

	// 4. LOOP MUHUD & AMBIL SHALAWATS-NYA
	for _, muhud := range mysqlMuhuds {
		// Ambil semua shalawats untuk muhud ini
		var shalawats []MySQLShalawat
		if err := mysqlDB.Where("muhud_id = ?", muhud.ID).Order("id asc").Find(&shalawats).Error; err != nil {
			log.Printf("Gagal load shalawats untuk muhud id %d: %v\n", muhud.ID, err)
			continue
		}

		// Loop dan Insert ke Postgres
		for i, ms := range shalawats {
			// Tentukan Category berdasarkan Logic TERBALIK khusus untuk Shalawats (Muhuds 7-14)
			var category string
			if ms.NumberOfDiwan == nil {
				category = "Diwan" // Jika numberOfDiwan Kosong, maka ini Diwan
			} else if ms.NumberOfMaulidSyarafulAnam == nil {
				category = "Diba" // Jika numberOfMaulidSyarafulAnam Kosong, maka ini Diba
			} else {
				log.Printf("Peringatan: Shalawat ID %d tidak memiliki kategori yang jelas (keduanya terisi / bukan null). Dilewati.\n", ms.ID)
				continue
			}

			// Cari ID Chapter di PostgreSQL berdasarkan transliteration_id (title) DAN category tersebut
			var pgChapter entity.Chapter
			if err := pgDB.Where("title = ? AND category = ?", muhud.TransliterationID, category).First(&pgChapter).Error; err != nil {
				log.Printf("Peringatan: Chapter '%s' (Kategori: %s) tidak ditemukan di PostgreSQL. Shalawat MySQL id %d dilewati.\n", muhud.TransliterationID, category, ms.ID)
				continue
			}

			verseNumber := uint(i + 1) // 1-indexed

			// Map transliteration
			transliteration := ms.Transliteration

			pgVerse := entity.Verse{
				ChapterID:       pgChapter.ID,
				VerseNumber:     verseNumber,
				ArabicText:      ms.TextShalawat,
				Transliteration: &transliteration,
				CreatedAt:       ms.CreatedAt,
				UpdatedAt:       ms.UpdatedAt,
			}

			// Insert Verse
			if err := pgDB.Create(&pgVerse).Error; err != nil {
				log.Printf("Gagal insert shalawat verse (chapter_id pg: %d, verse_number: %d). Error: %v\n", pgChapter.ID, pgVerse.VerseNumber, err)
				continue
			}
			successCount++

			// Insert Translation jika ada text di mysql translation_id (teks indonesia)
			if ms.TranslationID != "" {
				pgTranslation := entity.Translation{
					VerseID:         pgVerse.ID,
					LanguageCode:    "id",
					TranslationText: ms.TranslationID,
					CreatedAt:       ms.CreatedAt,
					UpdatedAt:       ms.UpdatedAt,
				}
				if err := pgDB.Create(&pgTranslation).Error; err != nil {
					log.Printf("Gagal insert translation untuk verse_id %d. Error: %v\n", pgVerse.ID, err)
					continue
				}
				translationSuccessCount++
			}
		}
	}

	fmt.Printf("Proses Migrasi Selesai! Berhasil migrasi %d verse dan %d translations.\n", successCount, translationSuccessCount)
}
