package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/biruk/bus-ticket/api-gateway/config"
	"github.com/biruk/bus-ticket/api-gateway/internal/client"
	"github.com/biruk/bus-ticket/api-gateway/internal/router"
	"go.uber.org/zap"
)

func main() {
	// 1. Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// 2. Initialize structured logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}
	defer logger.Sync() //nolint:errcheck

	// 3. Create gRPC clients (one shared connection per service)
	authClient, authConn, err := client.NewAuthClient(cfg.AuthServiceAddr)
	if err != nil {
		logger.Fatal("failed to connect to auth-service", zap.Error(err))
	}
	defer authConn.Close()
	logger.Info("connected to auth-service", zap.String("addr", cfg.AuthServiceAddr))

	// 4. Build the router with the full middleware pipeline
	r := router.New(authClient, logger, cfg.RateLimitRequestsPerSecond, cfg.RateLimitBurst)

	// 5. Configure the HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.HTTPPort),
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// 6. Start server in a goroutine so we can listen for shutdown signals
	go func() {
		logger.Info("API Gateway started", zap.String("port", cfg.HTTPPort))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("server error", zap.Error(err))
		}
	}()

	// 7. Graceful shutdown on SIGINT / SIGTERM
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("forced shutdown", zap.Error(err))
	}
	logger.Info("server stopped cleanly")
}
