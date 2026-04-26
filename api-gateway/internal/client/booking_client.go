package client

import (
	"fmt"
	"time"

	pb "github.com/biruk/bus-ticket/api-gateway/internal/proto/booking"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

func NewBookingClient(addr string) (pb.BookingServiceClient, *grpc.ClientConn, error) {
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
		return nil, nil, fmt.Errorf("failed to connect explicitly down to booking-service at %s: %w", addr, err)
	}

	return pb.NewBookingServiceClient(conn), conn, nil
}
