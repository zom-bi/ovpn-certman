package views

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/zom-bi/ovpn-certman/services"

	"github.com/gorilla/csrf"
)

type View struct {
	Vars         map[string]interface{}
	Request      *http.Request
	SessionStore *services.Sessions
}

func New(req *http.Request) *View {
	return &View{
		Request: req,
		Vars: map[string]interface{}{
			"CSRF_TOKEN": csrf.Token(req),
			"csrfField":  csrf.TemplateField(req),
			"Meta": map[string]interface{}{
				"Path": req.URL.Path,
				"Env":  "develop",
			},
			"flashes":  []services.Flash{},
			"username": "",
		},
	}
}

func NewWithSession(req *http.Request, sessionStore *services.Sessions) *View {
	return &View{
		Request:      req,
		SessionStore: sessionStore,
		Vars: map[string]interface{}{
			"CSRF_TOKEN": csrf.Token(req),
			"csrfField":  csrf.TemplateField(req),
			"Meta": map[string]interface{}{
				"Path": req.URL.Path,
				"Env":  "develop",
			},
			"flashes":  []services.Flash{},
			"username": sessionStore.GetUsername(req),
		},
	}
}

func (view View) Render(w http.ResponseWriter, name string) {
	var err error

	t, err := GetTemplate(name)
	if err != nil {
		log.Printf("the template '%s' does not exist.", name)
		view.RenderError(w, 404)
		return
	}

	if view.SessionStore != nil {
		// add flashes to template
		view.Vars["flashes"] = view.SessionStore.Flashes(w, view.Request)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	t.Execute(w, view.Vars)

}

func (view View) RenderError(w http.ResponseWriter, status int) {
	var name string

	switch status {
	case http.StatusUnauthorized:
		name = "401"
	case http.StatusForbidden:
		name = "403"
	case http.StatusNotFound:
		name = "404"
	default:
		name = "500"
	}

	t, err := GetTemplate(name)
	if err != nil {
		log.Printf("the error template '%s' does not exist.", name)
		fmt.Fprintf(w, "Error page for status '%d' could not be rendered.", status)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	t.Execute(w, view.Vars)

}

// GetTemplate returns a parsed template. The template ,ap needs to be
// Initialized by calling `LoadTemplates()` first.
func GetTemplate(name string) (*template.Template, error) {
	if tmpl, ok := templates[name]; ok {
		return tmpl, nil
	}

	return nil, errors.New("Template not found")
}
