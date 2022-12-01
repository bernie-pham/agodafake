package main

import (
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/bernie-pham/agodafake/db/driver"
	"github.com/bernie-pham/agodafake/internal/config"
	"github.com/bernie-pham/agodafake/internal/handlers"
	"github.com/bernie-pham/agodafake/internal/helpers"
	"github.com/bernie-pham/agodafake/internal/model"
	"github.com/bernie-pham/agodafake/internal/render"
)

var session *scs.SessionManager
var app config.AppConfig
var infoLog *log.Logger
var errLog *log.Logger
var dsn = "postgres://root:secret@localhost/booking_room?sslmode=disable"

func main() {
	// Tell the program which struct we will store in the session
	db, err := run()
	server := &http.Server{
		Addr:    ":8080",
		Handler: routes(&app),
	}
	defer db.SQL.Close()

	defer close(app.MailChan)
	listenForMail()

	err = server.ListenAndServe()
	log.Fatal(err)
}

func run() (*driver.DB, error) {

	gob.Register(model.Reservation{})
	gob.Register(model.Room{})
	gob.Register([]model.Room{})
	gob.Register(model.Restriction{})
	gob.Register(model.User{})
	app.InProduction = false

	mailChan := make(chan model.MailData)
	app.MailChan = mailChan

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog
	errLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.Errorlog = errLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session
	log.Println("Connecting to database...")
	db, err := driver.ConnectSQL(dsn)
	if err != nil {
		log.Fatal(err)
	}

	tc, err := render.LoadTemplate()
	if err != nil {
		log.Fatal("Cannot create template cache")
		return nil, err
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)
	return db, nil
}

// Homework
// - Finishing forgot password via email (generate changing password url)
// - Using redis for store session. prevening logout
