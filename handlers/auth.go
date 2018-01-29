package handlers

import (
	"net/http"

	"git.klink.asia/paul/certman/services"

	"git.klink.asia/paul/certman/models"
)

func RegisterHandler(p *services.Provider) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// Get parameters
		email := req.Form.Get("email")
		password := req.Form.Get("password")

		user := models.User{}
		user.Email = email
		user.SetPassword(password)

		err := p.DB.Create(&user).Error
		if err != nil {
			panic(err.Error)
		}

		p.Sessions.Flash(w, req,
			services.Flash{
				Type:    "success",
				Message: "The user was created. Check your inbox for the confirmation email.",
			},
		)

		http.Redirect(w, req, "/login", http.StatusFound)
		return
	}
}

func LoginHandler(p *services.Provider) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// Get parameters
		email := req.Form.Get("email")
		password := req.Form.Get("password")

		user := models.User{}

		err := p.DB.Where(&models.User{Email: email}).Find(&user).Error
		if err != nil {
			// could not find user
			p.Sessions.Flash(
				w, req, services.Flash{
					Type: "warning", Message: "Invalid Email or Password.",
				},
			)
			http.Redirect(w, req, "/login", http.StatusFound)
			return
		}

		if !user.EmailValid {
			p.Sessions.Flash(
				w, req, services.Flash{
					Type: "warning", Message: "You need to confirm your email before logging in.",
				},
			)
			http.Redirect(w, req, "/login", http.StatusFound)
			return
		}

		if err := user.CheckPassword(password); err != nil {
			// wrong password
			p.Sessions.Flash(
				w, req, services.Flash{
					Type: "warning", Message: "Invalid Email or Password.",
				},
			)
			http.Redirect(w, req, "/login", http.StatusFound)
			return
		}

		// user is logged in, set cookie
		p.Sessions.SetUserEmail(w, req, email)

		http.Redirect(w, req, "/certs", http.StatusSeeOther)
	}
}
