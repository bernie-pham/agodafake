package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/bernie-pham/agodafake/pkg/config"
	"github.com/bernie-pham/agodafake/pkg/model"
)

var functions = template.FuncMap{}
var app *config.AppConfig

func NewTemplates(a *config.AppConfig) {
	app = a
}

func RenderTemplateTODO(w http.ResponseWriter, tmpl string, data model.TodoPageData) {
	parsedTemplate, _ := template.ParseFiles("./templates/" + tmpl)
	if err := parsedTemplate.Execute(w, data); err != nil {
		fmt.Println("Error parsing template: ", err)
		return
	}
}
func RenderTemplate(w http.ResponseWriter, tmpl string, td *model.TemplateData) {
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
	_ = t.Execute(buf, td)
	_, err := buf.WriteTo(w)
	if err != nil {
		// log.Fatal(err)
		fmt.Println("Error writing template to browser", err)
	}
}

func LoadTemplate() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob("./templates/*.page.html")
	if err != nil {
		log.Fatal(err)
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			log.Fatal(err)
		}

		matches, err := filepath.Glob("./templates/*.layout.html")
		if err != nil {
			log.Fatal(err)
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.html")
			if err != nil {
				log.Fatal(err)
			}
		}

		myCache[name] = ts
	}
	return myCache, nil
}
