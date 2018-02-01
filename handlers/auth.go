package handlers

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"git.klink.asia/paul/certman/views"

	"github.com/go-chi/chi"

	"github.com/gorilla/securecookie"

	"git.klink.asia/paul/certman/services"

	"git.klink.asia/paul/certman/models"
)

func RegisterHandler(p *services.Provider) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// Get parameters
		email := req.Form.Get("email")

		user := models.User{}
		user.Email = email

		// don't set a password, user will get password reset request via mail
		user.HashedPassword = []byte{}

		err := p.DB.CreateUser(&user)
		if err != nil {
			panic(err.Error)
		}

		if err := createPasswordReset(p, &user); err != nil {
			p.Sessions.Flash(w, req,
				services.Flash{
					Type:    "danger",
					Message: "The registration email could not be generated.",
				},
			)
			http.Redirect(w, req, "/register", http.StatusFound)
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

		user, err := p.DB.GetUserByEmail(email)
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

func ConfirmEmailHandler(p *services.Provider) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		v := views.NewWithSession(req, p.Sessions)

		switch req.Method {
		case "GET":
			token := chi.URLParam(req, "token")
			pwr, err := p.DB.GetPasswordResetByToken(token)
			_ = pwr
			if err != nil {
				v.RenderError(w, 404)
				return
			}
			v.Render(w, "email-set-password")
		case "POST":
			password := req.Form.Get("password")
			token := req.Form.Get("token")
			pwr, err := p.DB.GetPasswordResetByToken(token)
			if err != nil {
				v.RenderError(w, 404)
				return
			}

			user, err := p.DB.GetUserByID(pwr.UserID)
			if err != nil {
				v.RenderError(w, 500)
				return
			}

			user.SetPassword(password)

			//err := p.DB.UpdateUser(user.ID, &user)
			if err != nil {
				v.RenderError(w, 500)
				return
			}

			err = p.DB.DeletePasswordResetsByUserID(pwr.UserID)

		default:
			v.RenderError(w, 405)
		}

		// try to get post params

		fmt.Fprintln(w, "Okay.")
	}
}

func createPasswordReset(p *services.Provider, user *models.User) error {
	// create the reset request
	pwr := models.PasswordReset{
		UserID:     user.ID,
		Token:      string(securecookie.GenerateRandomKey(32)),
		ValidUntil: time.Now().Add(6 * time.Hour),
	}

	if err := p.DB.CreatePasswordReset(&pwr); err != nil {
		return err
	}

	var subject string
	var text *bytes.Buffer

	if user.EmailValid {
		subject = "Password reset"
		text.WriteString("Somebody (hopefully you) has requested a password reset.\nClick below to reset your password:\n")
	} else {
		// If the user email has not been confirmed yet, send out
		// "mail confirmation"-mail instead
		subject = "Email confirmation"
		text.WriteString("Hello, thanks you for signing up!\nClick below to verify this email address\n")
	}

	return p.Email.Send(user.Email, subject, text.String(), "")
}
