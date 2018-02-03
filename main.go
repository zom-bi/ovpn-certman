package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"git.klink.asia/paul/certman/services"

	"git.klink.asia/paul/certman/router"
	"git.klink.asia/paul/certman/views"
)

func main() {
	log.Println("Initializing certman")
	if err := checkCAFilesExist(); err != nil {
		log.Fatalf("Could not read CA files: %s", err)
	}

	c := services.Config{
		CollectionPath: "./clients.json",
		Sessions: &services.SessionsConfig{
			SessionName: "_session",
			CookieKey:   os.Getenv("APP_KEY"),
			HttpOnly:    true,
			Lifetime:    24 * time.Hour,
		},
	}

	log.Println(".. services")
	serviceProvider := services.NewProvider(&c)

	// load and parse template files
	log.Println(".. templates")
	views.LoadTemplates()

	mux := router.HandleRoutes(serviceProvider)

	log.Println(".. server")
	err := http.ListenAndServe(":8000", mux)
	log.Fatalf(err.Error())
}

func checkCAFilesExist() error {
	for _, filename := range []string{"ca.crt", "ca.key"} {
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			return errors.New(filename + " not readable")
		}
	}
	return nil
}
