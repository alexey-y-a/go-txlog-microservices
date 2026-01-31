package server

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/alexey-y-a/go-txlog-microservices/libs/logger"
	"github.com/alexey-y-a/go-txlog-microservices/libs/txlog"
	kvhttp "github.com/alexey-y-a/go-txlog-microservices/services/kv-service/internal/http"
	kvmetrics "github.com/alexey-y-a/go-txlog-microservices/services/kv-service/internal/metrics"
	"github.com/alexey-y-a/go-txlog-microservices/services/kv-service/internal/store"
)

func NewServer() (*http.Server, *txlog.FileLog, error) {
	logger.Init()
	log := logger.L().With().Str("service", "kv-service").Logger()

	logFile, err := txlog.NewFileLog("kv.log")
	if err != nil {
		return nil, nil, err
	}

	kvStore := store.NewStore(logFile)

	mux := http.NewServeMux()

	handler := kvhttp.NewHandler(kvStore)

	mux.Handle("/health", kvmetrics.InstrumentHandler("health", http.HandlerFunc(handler.HealthHandler)))
	mux.Handle("/kv/set", kvmetrics.InstrumentHandler("kv_set", http.HandlerFunc(handler.SetHandler)))
	mux.Handle("/kv/get", kvmetrics.InstrumentHandler("kv_get", http.HandlerFunc(handler.GetHandler)))
	mux.Handle("/kv/delete", kvmetrics.InstrumentHandler("kv_delete", http.HandlerFunc(handler.DeleteHandler)))

	mux.Handle("/metrics", promhttp.Handler())

	addr := ":8081"
	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	log.Info().Str("addr", addr).Msg("kv-service http server created")

	return srv, logFile, nil
}
