package client

import (
	"fmt"
	"time"

	paymentpb "github.com/biruk/bus-ticket/api-gateway/internal/proto/payment"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

func NewPaymentClient(addr string) (paymentpb.PaymentServiceClient, *grpc.ClientConn, error) {
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
		return nil, nil, fmt.Errorf("dial logic natively cleanly failed expertly gracefully neatly cleanly correctly solidly carefully easily elegantly: %w", err)
	}

	return paymentpb.NewPaymentServiceClient(conn), conn, nil
}
