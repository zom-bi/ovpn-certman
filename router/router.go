package router

import (
	"net/http"
	"os"
	"strings"

	"git.klink.asia/paul/certman/services"

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

	//mux.Use(middleware.RequestID)
	mux.Use(middleware.Logger)          // log requests
	mux.Use(middleware.RealIP)          // use proxy headers
	mux.Use(middleware.RedirectSlashes) // redirect trailing slashes
	mux.Use(mw.Recoverer)               // recover on panic
	mux.Use(services.SessionStore.Use)  // use session storage

	// we are serving the static files directly from the assets package
	// this either means we use the embedded files, or live-load
	// from the file system (if `--tags="dev"` is used).
	fileServer(mux, "/static", assets.Assets)

	mux.Route("/", func(r chi.Router) {
		if os.Getenv("ENVIRONMENT") != "test" {
			r.Use(csrf.Protect(
				csrfKey,
				csrf.Secure(false),
				csrf.CookieName(csrfCookieName),
				csrf.FieldName(csrfFieldName),
				csrf.ErrorHandler(http.HandlerFunc(handlers.CSRFErrorHandler)),
			))
		}

		r.HandleFunc("/", v("debug"))

		r.Route("/register", func(r chi.Router) {
			r.Get("/", v("register"))
			r.Post("/", handlers.RegisterHandler)
		})

		r.Route("/login", func(r chi.Router) {
			r.Get("/", v("login"))
			r.Post("/", handlers.LoginHandler)
		})

		//r.Post("/confirm-email/{token}", handlers.ConfirmEmailHandler(db))

		r.Route("/forgot-password", func(r chi.Router) {
			r.Get("/", v("forgot-password"))
			r.Post("/", handlers.LoginHandler)
		})

		r.Route("/certs", func(r chi.Router) {
			r.Use(mw.RequireLogin)
			r.Get("/", handlers.ListCertHandler)
			r.Post("/new", handlers.CreateCertHandler)
			r.HandleFunc("/download/{ID}", handlers.DownloadCertHandler)
		})

		r.HandleFunc("/500", func(w http.ResponseWriter, req *http.Request) {
			panic("500")
		})
	})

	// what should happen if no route matches
	mux.NotFound(handlers.NotFoundHandler)

	return mux
}

// v is a helper function for quickly displaying a view
func v(template string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		view := views.New(req)
		view.Render(w, template)
	}
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
