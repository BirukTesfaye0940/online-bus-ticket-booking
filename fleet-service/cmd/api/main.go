package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/biruk/bus-ticket/fleet-service/config"
	"github.com/biruk/bus-ticket/fleet-service/internal/repository"
	"github.com/biruk/bus-ticket/fleet-service/internal/service"
	grpc_transport "github.com/biruk/bus-ticket/fleet-service/internal/transport/grpc"
	pb "github.com/biruk/bus-ticket/fleet-service/proto"
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
	fleetRepo := repository.NewFleetRepository(queries)
	fleetService := service.NewFleetService(fleetRepo)
	fleetHandler := grpc_transport.NewFleetHandler(fleetService)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterFleetServiceServer(grpcServer, fleetHandler)
	reflection.Register(grpcServer)

	log.Printf("Starting fleet-service gRPC server on port %s", cfg.GRPCPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
