package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/bernie-pham/agodafake/internal/config"
	"github.com/bernie-pham/agodafake/internal/model"
	"github.com/justinas/nosurf"
)

var functions = template.FuncMap{}
var app *config.AppConfig
var pathToTemplates string = "./templates"

func NewRenderer(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *model.TemplateData, req *http.Request) {
	td.Flash = app.Session.PopString(req.Context(), "flash")
	td.Warning = app.Session.PopString(req.Context(), "warning")
	td.Error = app.Session.PopString(req.Context(), "error")
	if app.Session.Exists(req.Context(), "user") {
		td.IsAuthenticated = 1
	}
	td.CSRFToken = nosurf.Token(req)
}

func Template(w http.ResponseWriter, req *http.Request, tmpl string, td *model.TemplateData) {
	var templates map[string]*template.Template
	if app.UseCache {
		templates = app.TemplateCache
	} else {
		templates, _ = LoadTemplate()
	}

	t, ok := templates[tmpl]
	if !ok {
		log.Fatal("Could not get template from template")
	}

	buf := new(bytes.Buffer)

	AddDefaultData(td, req)

	_ = t.Execute(buf, td)
	_, err := buf.WriteTo(w)
	if err != nil {
		// log.Fatal(err)
		fmt.Println("Error writing template to browser", err)
	}
}

func LoadTemplate() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.html", pathToTemplates))
	if err != nil {
		log.Fatal(err)
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			log.Fatal(err)
		}

		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.html", pathToTemplates))
		if err != nil {
			log.Fatal(err)
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.html", pathToTemplates))
			if err != nil {
				log.Fatal(err)
			}
		}

		myCache[name] = ts
	}
	return myCache, nil
}
