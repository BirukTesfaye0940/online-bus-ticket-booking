package service

import (
	"context"

	"github.com/biruk/bus-ticket/payment-service/internal/domain"
	stripeclient "github.com/biruk/bus-ticket/payment-service/internal/stripe"
)

type PaymentService interface {
	CreatePaymentIntent(ctx context.Context, userID, bookingID string, amount float64) (string, error)
	ConfirmPayment(ctx context.Context, stripeIntentID string) error
}

type paymentService struct {
	repo         domain.PaymentRepository
	stripeClient stripeclient.StripeClient
}

func NewPaymentService(r domain.PaymentRepository, sc stripeclient.StripeClient) PaymentService {
	return &paymentService{repo: r, stripeClient: sc}
}

func (s *paymentService) CreatePaymentIntent(ctx context.Context, userID, bookingID string, amount float64) (string, error) {
	intent, err := s.stripeClient.CreatePaymentIntent(amount, "etb", map[string]string{
		"booking_id": bookingID,
		"user_id":    userID,
	})
	if err != nil {
		return "", err
	}

	_, err = s.repo.CreateTransaction(ctx, bookingID, userID, amount, intent.ID, "PENDING")
	if err != nil {
		return "", err
	}

	return intent.ClientSecret, nil
}

func (s *paymentService) ConfirmPayment(ctx context.Context, stripeIntentID string) error {
	_, err := s.repo.UpdateTransactionStatusByIntent(ctx, stripeIntentID, "PAID")
	return err
}
