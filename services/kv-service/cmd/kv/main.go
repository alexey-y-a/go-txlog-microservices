package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alexey-y-a/go-txlog-microservices/libs/logger"
	"github.com/alexey-y-a/go-txlog-microservices/libs/txlog"
	kvhttp "github.com/alexey-y-a/go-txlog-microservices/services/kv-service/internal/http"
	"github.com/alexey-y-a/go-txlog-microservices/services/kv-service/internal/store"
)


func main() {
	logger.Init()
	log := logger.L().With().Str("service", "kv-service").Logger()

	logFile, err := txlog.NewFileLog("kv.log")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create transaction log")
	}

	kvStore := store.NewStore(logFile)

	mux := http.NewServeMux()

	handler := kvhttp.NewHandler(kvStore)
	handler.RegisterRoutes(mux)

	addr := ":8081"
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	log.Info().Str("addr", addr).Msg("starting kv-service")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		err = server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Error().Err(err).Msg("kv-service stopped with error")
		}
	}()

	sig := <-sigChan
	log.Info().Str("signal", sig.String()).Msg("shutting down kv-service")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Error().Err(err).Msg("kv-service graceful shutdown failed")
	} else {
		log.Info().Msg("kv-service stopped gracefully")
	}

}


