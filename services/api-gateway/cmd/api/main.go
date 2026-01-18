package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alexey-y-a/go-txlog-microservices/libs/logger"
	apihttp "github.com/alexey-y-a/go-txlog-microservices/services/api-gateway/internal/http"
)

func main() {
    logger.Init()
    log := logger.L().With().Str("service", "api-gateway").Logger()

     mux := http.NewServeMux()
     apihttp.RegisterRouters(mux)

     addr := "8080"
     srv :=  &http.Server {
         Addr: addr,
         Handler: mux,
     }

      log.Info().Str("addr", addr).Msg("starting api-gateway")

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
