package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/biruk/bus-ticket/auth-service/config"
	"github.com/biruk/bus-ticket/auth-service/internal/repository"
	"github.com/biruk/bus-ticket/auth-service/internal/service"
	transport "github.com/biruk/bus-ticket/auth-service/internal/transport/grpc"
	pb "github.com/biruk/bus-ticket/auth-service/proto"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	pool, err := pgxpool.New(context.Background(), cfg.DBURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	queries := repository.New(pool)
	userRepo := repository.NewUserRepository(queries)
	authService := service.NewAuthService(userRepo, cfg.JWTSecret, cfg.TokenDuration)
	authHandler := transport.NewAuthHandler(authService)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, authHandler)
	reflection.Register(grpcServer)

	log.Printf("starting gRPC server on port %s", cfg.GRPCPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
