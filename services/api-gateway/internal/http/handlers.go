package http

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/alexey-y-a/go-txlog-microservices/libs/logger"
	"github.com/alexey-y-a/go-txlog-microservices/services/api-gateway/internal/client"
)

type Handler struct {
    kvClient *client.KVClient
}

func NewHandler(kvClient *client.KVClient) *Handler {
    return &Handler {
        kvClient: kvClient,
    }
}

type healthResponse struct {
    Status string `json:"status"`
    Time   string `json:"time"`
}

func (h *Handler) RegisterRouters(mux *http.ServeMux) {
    mux.HandleFunc("/health", h.healthHandler)

    mux.HandleFunc("/api/set", h.setHandler)
    mux.HandleFunc("/api/get", h.getHandler)
    mux.HandleFunc("/api/delete", h.deleteHandler)
}


func (h *Handler) healthHandler(w http.ResponseWriter, r *http.Request) {
    log := logger.L().With().Str("handler", "health").Logger()

    response := healthResponse {
        Status: "ok",
        Time:   time.Now().UTC().Format(time.RFC3339),
    }

     w.Header().Set("Content-Type", "application/json")
     w.WriteHeader(http.StatusOK)

     err := json.NewEncoder(w).Encode(response)
     if err != nil {
         log.Error().Err(err).Msg("failed to write health response")
     }
}

type setRequest struct {
    Key   string `json:"key"`
    Value string `json:"value"`
}

type commonResponse struct {
    Status string `json:"status"`
    Message string `json:"message,omitempty"`
}

func (h *Handler) setHandler(w http.ResponseWriter, r *http.Request) {
    log := logger.L().With().Str("handler", "api_set").Logger()

    if r.Method != http.MethodPost {
        w.WriteHeader(http.StatusMethodNotAllowed)
        return
    }

    var req setRequest
    decoder := json.NewDecoder(r.Body)
    err := decoder.Decode(&req)
    if err != nil {
        log.Error().Err(err).Msg("failed to decode set request")
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    if req.Key == "" {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    err = h.kvClient.Set(req.Key, req.Value)
    if err != nil {
        log.Error().Err(err).Str("key", req.Key).Msg("kv-client set failed")
        w.WriteHeader(http.StatusBadGateway)
        return
    }

    response := commonResponse {
        Status: "ok",
        Message: "value set via api-gateway",
    }

     w.Header().Set("Content-Type", "application/json")
     w.WriteHeader(http.StatusOK)

     err = json.NewEncoder(w).Encode(response)
     if err != nil {
         log.Error().Err(err).Msg("failed to write set response")
     }
}

type getResponse struct {
    Status string `json:"status"`
    Value string `json:"value,omitempty"`
}

func (h *Handler) getHandler(w http.ResponseWriter, r *http.Request) {
    log := logger.L().With().Str("handler", "api_get").Logger()

    if r.Method != http.MethodGet {
        w.WriteHeader(http.StatusMethodNotAllowed)
        return
    }

    key := r.URL.Query().Get("key")
    if key == "" {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    value, ok, err := h.kvClient.Get(key)
    if err != nil {
        log.Error().Err(err).Str("key", key).Msg("kv-client get failed")
        w.WriteHeader(http.StatusBadGateway)
        return
    }

    if !ok {
        w.WriteHeader(http.StatusNotFound)
        return
    }

    response := getResponse {
        Status: "ok",
        Value: value,
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)

    err = json.NewEncoder(w).Encode(response)
    if err != nil {
        log.Error().Err(err).Msg("Failed to write get response")
    }
}

func (h *Handler) deleteHandler(w http.ResponseWriter, r *http.Request) {
	log := logger.L().With().Str("handler", "api_delete").Logger()

	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	key := r.URL.Query().Get("key")
	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := h.kvClient.Delete(key)
	if err != nil {
		log.Error().Err(err).Str("key", key).Msg("kv-client delete failed")
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	response := commonResponse{
		Status:  "ok",
		Message: "key deleted via api-gateway",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Error().Err(err).Msg("failed to write delete response")
	}
}