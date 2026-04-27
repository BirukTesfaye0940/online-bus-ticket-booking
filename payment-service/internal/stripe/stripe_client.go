package stripeclient

import (
	"fmt"

	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/paymentintent"
)

type StripeClient interface {
	CreatePaymentIntent(amount float64, currency string, metadata map[string]string) (*stripe.PaymentIntent, error)
}

type stripeImpl struct{}

func NewStripeClient(secretKey string) StripeClient {
	stripe.Key = secretKey
	return &stripeImpl{}
}

func (s *stripeImpl) CreatePaymentIntent(amount float64, currency string, metadata map[string]string) (*stripe.PaymentIntent, error) {
	// Scaling abstract value (e.g., dollars, birr) firmly into lowest units safely
	amountInScalingTarget := int64(amount * 100)

	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(amountInScalingTarget),
		Currency: stripe.String(currency),
	}
	
	for k, v := range metadata {
		params.AddMetadata(k, v)
	}

	intent, err := paymentintent.New(params)
	if err != nil {
		return nil, fmt.Errorf("failed executing native stripe payload dynamically smoothly cleanly accurately: %w", err)
	}

	return intent, nil
}
