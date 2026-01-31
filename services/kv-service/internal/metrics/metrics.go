package metrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var httpRequestsTotal = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Subsystem: "http",
		Name:      "requests_total",
		Help:      "Total number of HTTP requests handled by the kv-service.",
	},
	[]string{"handler", "method", "status"},
)

func InstrumentHandler(handlerName string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := &responseWriterWrapper{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		start := time.Now()

		next.ServeHTTP(ww, r)

		_ = start

		statusStr := strconv.Itoa(ww.statusCode)
		httpRequestsTotal.WithLabelValues(handlerName, r.Method, statusStr).Inc()
	})
}

type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

func (w *responseWriterWrapper) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}
