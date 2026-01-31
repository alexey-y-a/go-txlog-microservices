package server

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/alexey-y-a/go-txlog-microservices/libs/logger"
	"github.com/alexey-y-a/go-txlog-microservices/services/api-gateway/internal/client"
	apihttp "github.com/alexey-y-a/go-txlog-microservices/services/api-gateway/internal/http"
	apimetrics "github.com/alexey-y-a/go-txlog-microservices/services/api-gateway/internal/metrics"
)

func NewServer(kvBaseURL string) *http.Server {
	logger.Init()
	log := logger.L().With().Str("service", "api-gateway").Logger()

	kvTimeout := 3 * time.Second
	kvClient := client.NewKVClient(kvBaseURL, kvTimeout)

	mux := http.NewServeMux()

	handler := apihttp.NewHandler(kvClient)

	mux.Handle("/health", apimetrics.InstrumentHandler("health", http.HandlerFunc(handler.HealthHandler)))
	mux.Handle("/api/set", apimetrics.InstrumentHandler("api_set", http.HandlerFunc(handler.SetHandler)))
	mux.Handle("/api/get", apimetrics.InstrumentHandler("api_get", http.HandlerFunc(handler.GetHandler)))
	mux.Handle("/api/delete", apimetrics.InstrumentHandler("api_delete", http.HandlerFunc(handler.DeleteHandler)))

	mux.Handle("/metrics", promhttp.Handler())

	addr := ":8080"
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	log.Info().Str("addr", addr).Msg("api-gateway http server created")

	return server
}
