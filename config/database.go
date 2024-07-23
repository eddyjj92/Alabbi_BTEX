package config

import (
	"Alabbi_BTEX/models"
	"context"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var Db *gorm.DB
var Ctx context.Context

func InitDatabase() {
	db, err := gorm.Open(sqlite.Open("btex_db.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	Db = db

	// Migrate the schema
	db.AutoMigrate(&models.Process{})
}
