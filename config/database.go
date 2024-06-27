package config

import (
	"Alabbi_BTEX/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var Db *gorm.DB

func InitDatabase() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	Db = db

	// Migrate the schema
	db.AutoMigrate(&models.Process{})
}
