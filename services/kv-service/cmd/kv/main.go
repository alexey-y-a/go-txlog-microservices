package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alexey-y-a/go-txlog-microservices/libs/logger"
	"github.com/alexey-y-a/go-txlog-microservices/services/kv-service/internal/server"
)

func main() {
	logger.Init()
	log := logger.L().With().Str("service", "kv-service").Logger()

	srv, logFile, err := server.NewServer()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create kv-service")
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		err = srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Error().Err(err).Msg("kv-service stopped with error")
		}
	}()

	sig := <-sigChan
	log.Info().Str("signal", sig.String()).Msg("shutting down kv-service")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = srv.Shutdown(ctx)
	if err != nil {
		log.Error().Err(err).Msg("kv-service graceful shutdown failed")
	} else {
		log.Info().Msg("kv-service stopped gracefully")
	}

	err = logFile.Close()
    if err != nil {
         log.Error().Err(err).Msg("failed to close transaction log")
    } else {
        log.Info().Msg("transaction log closed")
    }
}
