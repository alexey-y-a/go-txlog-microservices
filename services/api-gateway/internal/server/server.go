package server

import (
	"net/http"
	"time"

	"github.com/alexey-y-a/go-txlog-microservices/libs/logger"
	"github.com/alexey-y-a/go-txlog-microservices/services/api-gateway/internal/client"
	apihttp "github.com/alexey-y-a/go-txlog-microservices/services/api-gateway/internal/http"
)

func NewServer(kvBaseURL string) *http.Server {
	logger.Init()
	log := logger.L().With().Str("service", "api-gateway").Logger()

	kvTimeout := 3 * time.Second
	kvClient := client.NewKVClient(kvBaseURL, kvTimeout)

	mux := http.NewServeMux()

	handler := apihttp.NewHandler(kvClient)
	handler.RegisterRoutes(mux)

	addr := ":8080"
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	log.Info().Str("addr", addr).Msg("api-gateway http server created")

	return server
}
