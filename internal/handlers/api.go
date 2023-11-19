package handlers

import (
	"github.com/go-chi/chi"
	chimiddle "github.com/go-chi/chi/middleware"
	"github.com/tiberiu1204/olx_web_scraper/internal/middleware"
)

func Handler(r *chi.Mux) {
	r.Use(chimiddle.StripSlashes)

	r.Route("/query", func(router chi.Router) {
		router.Use(middleware.Authorization)
		router.Get("/", GetQueryResults)
	})
}
