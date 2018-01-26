package services

import (
	"log"

	"git.klink.asia/paul/certman/models"
	"git.klink.asia/paul/certman/settings"
	"github.com/jinzhu/gorm"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	dsn := settings.Get("DATABASE_URL", "db.sqlite3")

	// Establish connection
	db, err := gorm.Open("sqlite3", dsn)
	if err != nil {
		log.Fatalf("Could not open database: %s", err.Error())
	}

	// Migrate models
	db.AutoMigrate(models.User{}, models.ClientConf{})

	return db
}
