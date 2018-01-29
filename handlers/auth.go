package handlers

import (
	"net/http"

	"git.klink.asia/paul/certman/services"

	"git.klink.asia/paul/certman/models"
)

func RegisterHandler(w http.ResponseWriter, req *http.Request) {
	// Get parameters
	email := req.Form.Get("email")
	password := req.Form.Get("password")

	user := models.User{}
	user.Email = email
	user.SetPassword(password)

	err := services.Database.Create(&user).Error
	if err != nil {
		panic(err.Error)
	}

	services.SessionStore.Flash(w, req,
		services.Flash{
			Type:    "success",
			Message: "The user was created. Check your inbox for the confirmation email.",
		},
	)

	http.Redirect(w, req, "/login", http.StatusFound)
	return
}

func LoginHandler(w http.ResponseWriter, req *http.Request) {
	// Get parameters
	email := req.Form.Get("email")
	password := req.Form.Get("password")

	user := models.User{}

	err := services.Database.Where(&models.User{Email: email}).Find(&user).Error
	if err != nil {
		// could not find user
		services.SessionStore.Flash(
			w, req, services.Flash{
				Type: "warning", Message: "Invalid Email or Password.",
			},
		)
		http.Redirect(w, req, "/login", http.StatusFound)
		return
	}

	if !user.EmailValid {
		services.SessionStore.Flash(
			w, req, services.Flash{
				Type: "warning", Message: "You need to confirm your email before logging in.",
			},
		)
		http.Redirect(w, req, "/login", http.StatusFound)
		return
	}

	if err := user.CheckPassword(password); err != nil {
		// wrong password
		services.SessionStore.Flash(
			w, req, services.Flash{
				Type: "warning", Message: "Invalid Email or Password.",
			},
		)
		http.Redirect(w, req, "/login", http.StatusFound)
		return
	}

	// user is logged in, set cookie
	services.SessionStore.SetUserEmail(w, req, email)

	http.Redirect(w, req, "/certs", http.StatusSeeOther)
}
