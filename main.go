package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/securecookie"

	"git.klink.asia/paul/certman/services"

	"git.klink.asia/paul/certman/router"
	"git.klink.asia/paul/certman/views"

	// import sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	c := services.Config{
		DB: &services.DBConfig{
			Type: "sqlite3",
			DSN:  "db.sqlite3",
			Log:  true,
		},
		Sessions: &services.SessionsConfig{
			SessionName: "_session",
			CookieKey:   string(securecookie.GenerateRandomKey(32)),
			HttpOnly:    true,
			Lifetime:    24 * time.Hour,
		},
	}

	serviceProvider := services.NewProvider(&c)

	// load and parse template files
	views.LoadTemplates()

	mux := router.HandleRoutes(serviceProvider)

	err := http.ListenAndServe(":8000", mux)
	log.Fatalf(err.Error())
}
