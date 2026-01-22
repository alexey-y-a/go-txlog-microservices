package http

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/alexey-y-a/go-txlog-microservices/libs/logger"
	"github.com/alexey-y-a/go-txlog-microservices/services/kv-service/internal/store"
)

type Handler struct {
    store *store.Store
}

func NewHandler(s *store.Store) *Handler {
    return &Handler {
        store: s,
    }
}

type healthResponse struct {
    Status string `json:"status"`
    Time string   `json:"time"`
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
    mux.HandleFunc("/health", healthHandler)

    mux.HandleFunc("/kv/set", h.setHandler)

    mux.HandleFunc("kv/get", h.getHandler)

    mux.HandleFunc("kv/del", h.deleteHandler)

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

type SetRequest struct {
    Key string `json:"key"`
    Value string `json:"value"`
}

type commonResponse struct {
    Status string `json:"status"`
    Message string `json:"message,omitempty"`
}

func (h *Handler) setHandler(w http.ResponseWriter, r *http.Request) {
    log := logger.L().With().Str("handler", "set").Logger()

    if r.Method != http.MethodPost {
        w.WriteHeader(http.StatusMethodNotAllowed)
        return
    }

    var req SetRequest
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

    err = h.store.Set(req.Key, req.Value)
    if err != nil {
        log.Error().Err(err).Str("key", req.Key).Msg("store set failed")

        w.WriteHeader(http.StatusInternalServerError)
        return
    }

     response := commonResponse {
         Status: "ok",
         Message: "value set",
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
    Value string `json:"value, omitempty"`
}

func (h *Handler) getHandler(w http.ResponseWriter, r *http.Request) {
	log := logger.L().With().Str("handler", "get").Logger()

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	key := r.URL.Query().Get("key")
	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	value, ok := h.store.Get(key)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	response := getResponse{
		Status: "ok",
		Value:  value,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Error().Err(err).Msg("failed to write get response")
	}
}

func (h *Handler) deleteHandler(w http.ResponseWriter, r *http.Request) {
	log := logger.L().With().Str("handler", "delete").Logger()

	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	key := r.URL.Query().Get("key")
	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := h.store.Delete(key)
	if err != nil {
		log.Error().Err(err).Str("key", key).Msg("store delete failed")

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := commonResponse{
		Status:  "ok",
		Message: "key deleted",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Error().Err(err).Msg("failed to write delete response")
	}
}