package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/skandyla/deploy-versions/config"
	"github.com/skandyla/deploy-versions/internal"
	"github.com/skandyla/deploy-versions/pkg/db"
)

func main() {
	config, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	dbc, err := db.NewConnection(config.PostgresDSN)
	if err != nil {
		log.Fatal(err)
	}

	storage := internal.NewVersionStorage(dbc)
	h := internal.NewVersionHandler(storage)

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)

	r.Route("/info", func(r chi.Router) {
		r.Get("/", h.Info)
	})

	r.Route("/versions", func(r chi.Router) {
		r.Get("/", h.GetAllVersions)
	})

	r.Route("/version", func(r chi.Router) {
		//r.Get("/", h.GetVersion)
		r.Post("/", h.PostVersion)
		r.Route("/{buildID}", func(r chi.Router) {
			r.Get("/", h.GetVersionByID)
			r.Put("/", h.PutVersionByID) //update entity
			r.Delete("/", h.DeleteVersionByID)
		})
	})

	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal(err)
	}
}
