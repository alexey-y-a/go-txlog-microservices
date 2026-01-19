package http

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/alexey-y-a/go-txlog-microservices/libs/logger"
)

type healthResponse struct {
    Status string `json:"status"`
    Time string   `json:"time"`
}

func RegisterRoutes(mux *http.ServeMux) {
    mux.HandleFunc("/health", healthHandler)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
    log := logger.L().With().Str("handler", "health").Logger()

    response := healthResponse{
        Status: "ok",
        Time: time.Now().UTC().Format(time.RFC3339),
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)

    err := json.NewEncoder(w).Encode(response)
    if err != nil {
        log.Error().Err(err).Msg("failed to write health response")
    }
}