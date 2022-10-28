package main

import (
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/bernie-pham/agodafake/pkg/config"
	"github.com/bernie-pham/agodafake/pkg/handlers"
	"github.com/bernie-pham/agodafake/pkg/render"
)

var session *scs.SessionManager
var app config.AppConfig

func main() {
	app.InProduction = false
	tc, err := render.LoadTemplate()

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	if err != nil {
		log.Fatal("Cannot create template cache")
	}

	app.TemplateCache = tc
	app.UseCache = true

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)
	render.NewTemplates(&app)

	// fs := http.FileServer(http.Dir("./images/"))
	// http.Handle("/images/", http.StripPrefix("/images/", fs))

	server := &http.Server{
		Addr:    ":8080",
		Handler: routes(&app),
	}

	err = server.ListenAndServe()
	log.Fatal(err)
}
