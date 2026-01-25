package server

import (
	"net/http"

	"github.com/alexey-y-a/go-txlog-microservices/libs/logger"
	"github.com/alexey-y-a/go-txlog-microservices/libs/txlog"
	kvhttp "github.com/alexey-y-a/go-txlog-microservices/services/kv-service/internal/http"
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
	handler.RegisterRoutes(mux)

	addr := ":8081"
	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	log.Info().Str("addr", addr).Msg("kv-service http server created")

	return srv, logFile, nil
}
