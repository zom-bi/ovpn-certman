package views

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"git.klink.asia/paul/certman/assets"
)

// map of all parsed templates, by template name
var templates map[string]*template.Template

// LoadTemplates initializes the templates map, parsing all defined templates.
func LoadTemplates() {
	templates = map[string]*template.Template{
		"401": newTemplate("layouts/application.gohtml", "errors/401.gohtml"),
		"404": newTemplate("layouts/application.gohtml", "errors/404.gohtml"),
		"500": newTemplate("layouts/application.gohtml", "errors/500.gohtml"),

		"debug":     newTemplate("layouts/application.gohtml", "shared/header.gohtml", "shared/footer.gohtml", "views/debug.gohtml"),
		"login":     newTemplate("layouts/application.gohtml", "shared/header.gohtml", "shared/footer.gohtml", "views/login.gohtml"),
		"cert_list": newTemplate("layouts/application.gohtml", "shared/header.gohtml", "shared/footer.gohtml", "views/cert_list.gohtml"),
	}
	return
}

// newTemplate returns a new template from the assets
func newTemplate(filenames ...string) *template.Template {
	f := []string{}
	prefix := "/templates"

	for _, filename := range filenames {
		f = append(f, filepath.Join(prefix, filename))
	}

	baseTemplate := template.New("base").Funcs(funcs)
	tmpl, err := parseAssets(baseTemplate, assets.Assets, f...)
	if err != nil {
		log.Fatalf("could not parse template: %s", err.Error())
	}

	return tmpl

}

// parseAssets is a helper function to generate a template from multiple
// assets. If the argument template is nil, it is created from the first
// parameter that is passed (first file).
func parseAssets(t *template.Template, fs http.FileSystem, assets ...string) (*template.Template, error) {
	if len(assets) == 0 {
		// Not really a problem, but be consistent.
		return nil, fmt.Errorf("no templates supplied in call to parseAssets")
	}

	for _, filename := range assets {
		f, err := fs.Open(filename)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		buf := new(bytes.Buffer)
		buf.ReadFrom(f)
		s := buf.String()

		name := filepath.Base(filename)
		// First template becomes return value if not already defined,
		// and we use that one for subsequent New calls to associate
		// all the templates together.
		var tmpl *template.Template
		if t == nil {
			t = template.New(name)
		}
		if name == t.Name() {
			tmpl = t
		} else {
			tmpl = t.New(name)
		}
		_, err = tmpl.Parse(s)
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}
