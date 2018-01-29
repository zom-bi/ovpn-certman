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
		Email: &services.EmailConfig{
			SMTPServer:   "example.com",
			SMTPPort:     25,
			SMTPUsername: "test",
			SMTPPassword: "test",
			From:         "Mailtest <test@example.com>",
		},
	}

	serviceProvider := services.NewProvider(&c)

	// Start the mail daemon, which re-uses connections to send mails to the
	// SMTP server
	go serviceProvider.Email.Daemon()

	//user := models.User{}
	//user.Username = "test"
	//user.SetPassword("test")
	//fmt.Println(user.HashedPassword)
	//fmt.Println(db.Create(&user).Error)

	// load and parse template files
	views.LoadTemplates()

	mux := router.HandleRoutes(serviceProvider)

	err := http.ListenAndServe(":8000", mux)
	log.Fatalf(err.Error())
}
