package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/rahulprgrmr/room-booking-web-app/pkg/config"
	"github.com/rahulprgrmr/room-booking-web-app/pkg/handlers"
	"github.com/rahulprgrmr/room-booking-web-app/pkg/render"
)

const portNumber = ":8080"
var app config.AppConfig
var session *scs.SessionManager

func main() {

	// change this to true when in production
	app.InProduction = false

	session = scs.New()
	session.Lifetime = time.Hour * 24
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	templateCache, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}
	app.TemplateCache = templateCache
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

	// http.HandleFunc("/", handlers.Repo.Home)
	// http.HandleFunc("/about", handlers.Repo.About)

	fmt.Println(fmt.Sprintf("Starting application on port %s", portNumber))
	// http.ListenAndServe(portNumber, nil)

	serve := &http.Server {
		Addr: portNumber,
		Handler: routes(&app),
	}

	err = serve.ListenAndServe()
	log.Fatal(err)
}