package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	log "github.com/sirupsen/logrus"
	"github.com/skandyla/deploy-versions/config"
	"github.com/skandyla/deploy-versions/internal"
	"github.com/skandyla/deploy-versions/pkg/db"
)

func main() {
	config, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	initLogger(config.LogLevel, config.JsonLogOutput)

	dbc, err := db.NewConnection(config.PostgresDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := dbc.Close(); err != nil {
			log.Printf("error closing database: %v", err)
		}
	}()

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

	server := http.Server{
		Addr:           config.ListenAddress,
		Handler:        r,
		ReadTimeout:    config.ReadTimeout,
		WriteTimeout:   config.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	serverErrors := make(chan error, 1)
	go func() {
		log.Printf("server started, listening on %s", server.Addr)
		serverErrors <- server.ListenAndServe()
	}()

	//------------------------------
	//shutdown

	// Blocking main and waiting for shutdown of the daemon.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Waiting for an osSignal or a non-HTTP related server error.
	select {
	case err := <-serverErrors:
		log.Printf("server error: %w", err)
		return

	case sig := <-shutdown:
		log.Info("shutdown started, signal: ", sig)
		//log.WithFields(log.Fields{"shutdown_status": "started"}).Info(sig)
		defer log.Info("shutdown complete, signal: ", sig)

		// Give outstanding requests a deadline for completion.
		ctx, cancel := context.WithTimeout(context.Background(), config.ShutdownTimeout)
		defer cancel()

		// Asking listener to shutdown and shed load.
		if err := server.Shutdown(ctx); err != nil {
			server.Close()
			log.Printf("could not stop server gracefully: %w", err)
			return
		}
	}
}

//------------------------------
func initLogger(logLevel string, json bool) {
	if json {
		log.SetFormatter(&log.JSONFormatter{})
	}
	log.SetOutput(os.Stderr)

	switch strings.ToLower(logLevel) {
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	default:
		log.SetLevel(log.DebugLevel)
	}
}
