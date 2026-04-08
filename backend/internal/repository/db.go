package repository

import (
	"fmt"
	"life-financial-assistant-backend/internal/model"
	"log"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(dbPath string) {
	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	fmt.Println("Database connection established")

	// Auto migration
	err = DB.AutoMigrate(
		&model.User{},
		&model.Space{},
		&model.SpaceUserLink{},
		&model.Bill{},
		&model.AnalysisReport{},
	)
	if err != nil {
		log.Fatalf("failed to auto migrate: %v", err)
	}

	fmt.Println("Database migration completed")
}
