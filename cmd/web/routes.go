package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rahulprgrmr/room-booking-web-app/pkg/config"
	"github.com/rahulprgrmr/room-booking-web-app/pkg/handlers"
)

func routes(app *config.AppConfig) http.Handler {
	// routing with pat module

	// mux := pat.New()

	// mux.Get("/", http.HandlerFunc(handlers.Repo.Home))
	// mux.Get("/about", http.HandlerFunc(handlers.Repo.About))

	// routing with chi module

	mux := chi.NewRouter()

	mux.Use(middleware.Logger, middleware.Recoverer, NoSurf, SessionLoad)

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)

	return mux
}