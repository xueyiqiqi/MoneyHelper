package repository

import (
	"fmt"
	"life-financial-assistant-backend/internal/config"
	"life-financial-assistant-backend/internal/model"
	"log"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(cfg config.DatabaseConfig) {
	var err error
	var db *gorm.DB

	if cfg.UseSQLite {
		// SQLite 模式
		db, err = gorm.Open(sqlite.Open(cfg.Path), &gorm.Config{})
	} else {
		// MySQL 模式
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	}

	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	DB = db
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
