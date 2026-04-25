package client

import (
	"fmt"
	"time"

	pb "github.com/biruk/bus-ticket/api-gateway/internal/proto/fleet"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

// NewFleetClient creates a shared, keep-alive gRPC connection to the Fleet Service.
func NewFleetClient(addr string) (pb.FleetServiceClient, *grpc.ClientConn, error) {
	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                10 * time.Second,
			Timeout:             3 * time.Second,
			PermitWithoutStream: true,
		}),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to fleet-service at %s: %w", addr, err)
	}

	return pb.NewFleetServiceClient(conn), conn, nil
}
