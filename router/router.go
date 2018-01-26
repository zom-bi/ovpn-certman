package router

import (
	"net/http"
	"os"
	"strings"

	"git.klink.asia/paul/certman/assets"
	"git.klink.asia/paul/certman/handlers"
	"git.klink.asia/paul/certman/views"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/csrf"
	"github.com/jinzhu/gorm"

	mw "git.klink.asia/paul/certman/middleware"
)

var (
	// TODO: make this configurable
	csrfCookieName = "csrf"
	csrfFieldName  = "csrf_token"
	csrfKey        = []byte("7Oj4DllZ9lTsxJnisTuWiiQBGQIzi6gX")
	cookieKey      = []byte("osx70sMD8HZG2ouUl8uKI4wcMugiJ2WH")
)

func HandleRoutes(db *gorm.DB) http.Handler {
	mux := chi.NewMux()

	//	mux.Use(middleware.RequestID)
	mux.Use(middleware.Logger)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.RedirectSlashes)
	mux.Use(mw.Recoverer)

	// we are serving the static files directly from the assets package
	fileServer(mux, "/static", assets.Assets)

	mux.Route("/", func(r chi.Router) {
		if os.Getenv("ENVIRONMENT") != "test" {
			r.Use(csrf.Protect(
				csrfKey,
				csrf.Secure(false),
				csrf.CookieName(csrfCookieName),
				csrf.FieldName(csrfFieldName),
			))
		}

		r.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
			view := views.New(req)
			view.Render(w, "debug")
		})

		r.Get("/login", func(w http.ResponseWriter, req *http.Request) {
			view := views.New(req)
			view.Render(w, "login")
		})

		r.Get("/certs", handlers.ListCertHandler(db))
		r.HandleFunc("/certs/new", handlers.GenCertHandler(db))

		r.HandleFunc("/500", func(w http.ResponseWriter, req *http.Request) {
			panic("500")
		})
	})

	// what should happen if no route matches
	mux.NotFound(handlers.NotFoundHandler)

	return mux
}

// fileServer sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func fileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	//fs := http.StripPrefix(path, http.FileServer(root))
	fs := http.FileServer(root)

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}
