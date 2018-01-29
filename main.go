package main

import (
	"log"
	"net/http"

	"git.klink.asia/paul/certman/services"

	"git.klink.asia/paul/certman/router"
	"git.klink.asia/paul/certman/views"

	// import sqlite3 driver once
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	// Connect to the database
	db := services.InitDB()

	services.InitSession()

	//user := models.User{}
	//user.Username = "test"
	//user.SetPassword("test")
	//fmt.Println(user.HashedPassword)
	//fmt.Println(db.Create(&user).Error)

	// load and parse template files
	views.LoadTemplates()

	mux := router.HandleRoutes(db)

	err := http.ListenAndServe(":8000", mux)
	log.Fatalf(err.Error())
}
