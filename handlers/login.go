package handlers

import (
	"net/http"

	"git.klink.asia/paul/certman/models"
	"github.com/jinzhu/gorm"
)

func LoginHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// Get parameters
		username := req.Form.Get("username")
		password := req.Form.Get("password")

		user := models.User{}

		err := db.Where(&models.User{Username: username}).Find(&user).Error
		if err != nil {
			// could not find user
			http.Redirect(w, req, "/login", http.StatusFound)
		}

		if err := user.CheckPassword(password); err != nil {
			// wrong password
			http.Redirect(w, req, "/login", http.StatusFound)
		}

		// user is logged in
		// set cookie
		http.Redirect(w, req, "/certs", http.StatusFound)
	}
}
