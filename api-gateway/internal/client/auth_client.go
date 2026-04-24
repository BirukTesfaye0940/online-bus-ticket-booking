package client

import (
	"fmt"
	"time"

	pb "github.com/biruk/bus-ticket/api-gateway/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

// NewAuthClient creates a shared, keep-alive gRPC connection to the Auth Service.
// This connection is created once at startup and reused across all requests.
func NewAuthClient(addr string) (pb.AuthServiceClient, *grpc.ClientConn, error) {
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
		return nil, nil, fmt.Errorf("failed to connect to auth-service at %s: %w", addr, err)
	}

	return pb.NewAuthServiceClient(conn), conn, nil
}
