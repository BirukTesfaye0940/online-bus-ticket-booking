package router

import (
	"github.com/biruk/bus-ticket/api-gateway/internal/handler"
	"github.com/biruk/bus-ticket/api-gateway/internal/middleware"
	pb "github.com/biruk/bus-ticket/api-gateway/internal/proto"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

// New constructs the chi router with the full middleware pipeline and all routes.
func New(authClient pb.AuthServiceClient, logger *zap.Logger, rps float64, burst int) *chi.Mux {
	r := chi.NewRouter()

	// --- Internal: Prometheus metrics scrape endpoint ---
	// Exposed on the same port; in prod you can move this to a separate internal port.
	r.Handle("/metrics", promhttp.Handler())

	// --- Global middleware (applied to every request) ---
	r.Use(middleware.RequestID)
	r.Use(middleware.Metrics) // ← Prometheus counter/histogram per request
	r.Use(middleware.Logger(logger))
	r.Use(middleware.Recoverer(logger))
	r.Use(chiMiddleware.StripSlashes)
	r.Use(middleware.RateLimiter(rps, burst))

	// --- Handlers ---
	authHandler := handler.NewAuthHandler(authClient)

	// --- Routes ---
	r.Route("/api/v1", func(r chi.Router) {

		// Auth routes — mixed public/protected
		r.Route("/auth", func(r chi.Router) {
			// Public
			r.Post("/register", authHandler.Register)
			r.Post("/login", authHandler.Login)

			// Protected (requires valid JWT)
			r.Group(func(r chi.Router) {
				r.Use(middleware.Auth(authClient))
				r.Get("/me", authHandler.GetMe)
			})
		})
	})

	return r
}
