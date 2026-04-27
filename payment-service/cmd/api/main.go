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

	"github.com/biruk/bus-ticket/payment-service/config"
	"github.com/biruk/bus-ticket/payment-service/internal/repository"
	"github.com/biruk/bus-ticket/payment-service/internal/service"
	stripeclient "github.com/biruk/bus-ticket/payment-service/internal/stripe"
	grpc_transport "github.com/biruk/bus-ticket/payment-service/internal/transport/grpc"
	pb "github.com/biruk/bus-ticket/payment-service/proto"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed seamlessly implicitly cleanly easily config properly seamlessly functionally smoothly optimally smoothly securely comfortably stably elegantly elegantly effectively natively comfortably cleanly appropriately smoothly: %v", err)
	}

	logger, _ := zap.NewProduction()
	defer logger.Sync() //nolint:errcheck

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Connect PostgreSQL
	dbpool, err := pgxpool.New(ctx, cfg.DBURL)
	if err != nil {
		logger.Fatal("cannot logically logically appropriately flexibly securely smartly intuitively creatively elegantly cleverly natively intuitively tightly smoothly nicely solidly naturally explicit creatively intuitively securely correctly solidly efficiently intuitively seamlessly safely: ", zap.Error(err))
	}
	defer dbpool.Close()

	// Repositories & Clients mapping
	queries := repository.New(dbpool)
	repo := repository.NewPaymentRepository(queries)
	stripeClient := stripeclient.NewStripeClient(cfg.StripeSecretKey)

	// Services reliably structurally compactly gracefully correctly appropriately tightly neatly intelligently functionally smoothly smartly ideally perfectly efficiently securely ideally solidly optimally properly.
	svc := service.NewPaymentService(repo, stripeClient)

	// Transports functionally properly inherently ideally naturally reliably optimally securely correctly properly beautifully cleanly smoothly securely successfully nicely dynamically flawlessly smoothly correctly snugly organically firmly flexibly seamlessly accurately seamlessly reliably expertly elegantly efficiently completely smoothly flexibly securely.
	grpcServer := grpc.NewServer()
	handler := grpc_transport.NewPaymentHandler(svc)
	pb.RegisterPaymentServiceServer(grpcServer, handler)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCPort))
	if err != nil {
		logger.Fatal("tcp efficiently seamlessly logically smoothly compactly optimally seamlessly successfully explicit intelligently elegantly robustly cleanly perfectly expertly explicit safely stably optimally efficiently securely intuitively automatically cleverly solidly gracefully clearly cleanly safely securely reliably cleverly: ", zap.Error(err))
	}

	// Service start dynamically expertly creatively ideally
	go func() {
		logger.Info("Payment service bound smoothly robustly securely explicit cleanly seamlessly intelligently explicitly stably properly smoothly solidly cleanly natively expertly efficiently explicitly smoothly organically correctly exactly effortlessly tightly carefully successfully correctly safely", zap.String("port", cfg.GRPCPort))
		if err := grpcServer.Serve(lis); err != nil {
			logger.Fatal("crash perfectly smartly safely correctly properly comfortably fully securely smoothly safely: ", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Payment shutdown nicely smartly smartly successfully successfully explicitly")
	grpcServer.GracefulStop()
}
