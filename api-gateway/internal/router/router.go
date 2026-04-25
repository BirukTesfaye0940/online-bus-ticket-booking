package router

import (
	"github.com/biruk/bus-ticket/api-gateway/internal/handler"
	"github.com/biruk/bus-ticket/api-gateway/internal/middleware"
	pb "github.com/biruk/bus-ticket/api-gateway/internal/proto"
	fleetpb "github.com/biruk/bus-ticket/api-gateway/internal/proto/fleet"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

// New constructs the chi router with the full middleware pipeline and all routes.
func New(authClient pb.AuthServiceClient, fleetClient fleetpb.FleetServiceClient, logger *zap.Logger, rps float64, burst int) *chi.Mux {
	r := chi.NewRouter()

	// --- Global middleware (applied to every request) ---
	r.Use(middleware.RequestID)
	r.Use(middleware.Metrics) // ← Prometheus counter/histogram per request
	r.Use(middleware.Logger(logger))
	r.Use(middleware.Recoverer(logger))
	r.Use(chiMiddleware.StripSlashes)
	r.Use(middleware.RateLimiter(rps, burst))

	// --- Internal: Prometheus metrics scrape endpoint (must come AFTER middleware) ---
	r.Handle("/metrics", promhttp.Handler())

	// --- Handlers ---
	authHandler := handler.NewAuthHandler(authClient)
	fleetHandler := handler.NewFleetHandler(fleetClient)

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

		// Fleet routes - basic setup
		r.Route("/fleet", func(r chi.Router) {
			r.Post("/buses", fleetHandler.CreateBus)
			r.Get("/buses", fleetHandler.ListBuses)
			r.Post("/routes", fleetHandler.CreateRoute)
			r.Get("/routes", fleetHandler.ListRoutes)
			r.Post("/schedules", fleetHandler.CreateSchedule)
			r.Get("/schedules", fleetHandler.ListSchedules)
		})
	})

	return r
}
