package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	log "github.com/sirupsen/logrus"
	"github.com/skandyla/deploy-versions/config"

	"github.com/skandyla/deploy-versions/internal/repository"
	"github.com/skandyla/deploy-versions/internal/service"
	"github.com/skandyla/deploy-versions/internal/transport"
	"github.com/skandyla/deploy-versions/pkg/db"
	"github.com/skandyla/deploy-versions/pkg/hash"
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
		log.Println("clossing database connection")
	}()

	// init deps
	fmt.Println("tokenttl:", config.Auth.TokenTTL)
	fmt.Println("logLevel:", config.LogLevel)

	hasher := hash.NewSHA1Hasher("salt")

	versionsRepository := repository.NewVersionRepository(dbc)
	versionsService := service.NewVersions(versionsRepository)

	tokensRepo := repository.NewTokens(dbc)
	usersRepo := repository.NewUsers(dbc)
	//usersService := service.NewUsers(usersRepo, hasher, []byte("sample secret"), config.Auth.TokenTTL)
	usersService := service.NewUsers(usersRepo, tokensRepo, hasher, []byte("sample secret"), config.Auth.TokenTTL)

	handler := transport.NewHandler(versionsService, usersService)

	server := http.Server{
		Addr:           config.ListenAddress,
		Handler:        handler.InitRouter(),
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
	// Blocking main and waiting for shutdown of the daemon.

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Waiting for an osSignal or a non-HTTP related server error.
	select {
	case err := <-serverErrors:
		log.Printf("server error: %v", err)
		return

	case sig := <-quit:
		log.Info("shutdown started, signal: ", sig)
		//log.WithFields(log.Fields{"shutdown_status": "started"}).Info(sig)
		defer log.Info("shutdown complete, signal: ", sig)

		// Give outstanding requests a deadline for completion.
		ctx, cancel := context.WithTimeout(context.Background(), config.ShutdownTimeout)
		defer cancel()

		// Asking listener to shutdown and shed load.
		if err := server.Shutdown(ctx); err != nil {
			server.Close()
			log.Printf("could not stop server gracefully: %v", err)
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
