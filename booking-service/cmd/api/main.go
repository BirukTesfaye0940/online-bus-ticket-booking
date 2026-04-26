package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/biruk/bus-ticket/booking-service/config"
	"github.com/biruk/bus-ticket/booking-service/internal/repository"
	"github.com/biruk/bus-ticket/booking-service/internal/service"
	grpc_transport "github.com/biruk/bus-ticket/booking-service/internal/transport/grpc"
	pb "github.com/biruk/bus-ticket/booking-service/proto"
	paymentpb "github.com/biruk/bus-ticket/booking-service/proto"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load booking config natively: %v", err)
	}

	logger, _ := zap.NewProduction()
	defer logger.Sync() //nolint:errcheck

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 1. Setup Database Connection
	dbpool, err := pgxpool.New(ctx, cfg.DBURL)
	if err != nil {
		logger.Fatal("cannot connect dynamically to db natively", zap.Error(err))
	}
	defer dbpool.Close()
	logger.Info("connected onto posgres database seamlessly")

	// 2. Setup Redis Connection
	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.RedisURL,
	})
	if err := rdb.Ping(ctx).Err(); err != nil {
		logger.Fatal("failed redis validation", zap.Error(err))
	}
	defer rdb.Close()
	logger.Info("connected onto redis locks seamlessly")

	// 3. Mount Repositories natively
	queries := repository.New(dbpool)
	bookingRepo := repository.NewBookingRepository(queries)
	redisLockLayer := repository.NewRedisLock(rdb)

	// 4. Form Payment Service connection dynamically leveraging native gRPC stubs securely
	paymentConn, err := grpc.NewClient(
		cfg.PaymentServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logger.Fatal("failed binding the upstream payment client gracefully", zap.Error(err))
	}
	defer paymentConn.Close()
	paymentClient := paymentpb.NewPaymentServiceClient(paymentConn)

	// 5. Mount Core Booking Service bridging structurally mappings logically over domains sequentially
	bookingSvc := service.NewBookingService(bookingRepo, redisLockLayer, paymentClient)

	// 6. Connect Transports natively exposing local mapping logic
	grpcServer := grpc.NewServer()
	bookingHandler := grpc_transport.NewBookingHandler(bookingSvc)
	pb.RegisterBookingServiceServer(grpcServer, bookingHandler)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCPort))
	if err != nil {
		logger.Fatal("failed TCP establishing logic naturally", zap.Error(err))
	}

	// Start Serving Gracefully
	go func() {
		logger.Info("Booking gRPC listening cleanly on mapping binds naturally", zap.String("port", cfg.GRPCPort))
		if err := grpcServer.Serve(lis); err != nil {
			logger.Fatal("Server dropped execution dynamically natively", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down booking instance dynamically mapping stops securely")
	grpcServer.GracefulStop()
	logger.Info("Exited cleanly stopping")
}
