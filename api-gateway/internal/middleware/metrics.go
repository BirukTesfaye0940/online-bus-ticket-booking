package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// httpRequestsTotal counts the total number of HTTP requests.
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gateway_http_requests_total",
			Help: "Total number of HTTP requests processed by the API Gateway.",
		},
		[]string{"method", "path", "status"},
	)

	// httpRequestDuration tracks request latency as a histogram.
	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "gateway_http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds.",
			Buckets: prometheus.DefBuckets, // [.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10]
		},
		[]string{"method", "path"},
	)

	// activeRequests tracks how many requests are currently in-flight.
	activeRequests = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "gateway_http_active_requests",
		Help: "Number of HTTP requests currently being processed.",
	})
)

// Metrics is a middleware that records Prometheus metrics for every request:
//   - gateway_http_requests_total (counter, labeled by method/path/status)
//   - gateway_http_request_duration_seconds (histogram, labeled by method/path)
//   - gateway_http_active_requests (gauge of in-flight requests)
func Metrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		activeRequests.Inc()
		defer activeRequests.Dec()

		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(rw, r)

		duration := time.Since(start).Seconds()
		status := strconv.Itoa(rw.statusCode)

		httpRequestsTotal.WithLabelValues(r.Method, r.URL.Path, status).Inc()
		httpRequestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration)
	})
}
