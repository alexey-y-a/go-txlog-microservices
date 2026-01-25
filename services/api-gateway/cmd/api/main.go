package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alexey-y-a/go-txlog-microservices/libs/logger"
	"github.com/alexey-y-a/go-txlog-microservices/services/api-gateway/internal/server"
)

func main() {
	logger.Init()
	log := logger.L().With().Str("service", "api-gateway").Logger()

	kvBaseURL := "http://kv-service:8081"

	srv := server.NewServer(kvBaseURL)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Error().Err(err).Msg("api-gateway stopped with error")
		}
	}()

	sig := <-sigChan
	log.Info().Str("signal", sig.String()).Msg("shutting down api-gateway")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := srv.Shutdown(ctx)
	if err != nil {
		log.Error().Err(err).Msg("api-gateway graceful shutdown failed")
	} else {
		log.Info().Msg("api-gateway stopped gracefully")
	}
}
