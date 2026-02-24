package main

import (
	"fmt"
	"ishari-backend/pkg/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg := config.Load()
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Jakarta",
		cfg.Database.Host, cfg.Database.User, cfg.Database.Password, cfg.Database.DBName, cfg.Database.Port, cfg.Database.SSLMode)

	pgDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	type Result struct {
		Indexname string `gorm:"column:indexname"`
		Indexdef  string `gorm:"column:indexdef"`
	}
	var res Result
	pgDB.Raw(`
		SELECT indexname, indexdef
		FROM pg_indexes
		WHERE indexname = 'unique_verse_number'
	`).Scan(&res)

	fmt.Printf("Index name: %s\nDefinition: %s\n", res.Indexname, res.Indexdef)
}
