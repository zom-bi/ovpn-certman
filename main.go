package main

import (
	"log"
	"net/http"

	"github.com/jinzhu/gorm"

	"git.klink.asia/paul/certman/models"
	"git.klink.asia/paul/certman/router"
	"git.klink.asia/paul/certman/views"

	// import sqlite3 driver once
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	// Connect to the database
	db, err := gorm.Open("sqlite3", "db.sqlite3")
	if err != nil {
		log.Fatalf("Could not open database: %s", err.Error())
	}
	defer db.Close()

	// Migrate
	db.AutoMigrate(models.User{}, models.ClientConf{})

	// load and parse template files
	views.LoadTemplates()

	mux := router.HandleRoutes(db)

	err = http.ListenAndServe(":8000", mux)
	log.Fatalf(err.Error())
}
