package domain

import (
	"context"
	"time"
)

type Transaction struct {
	ID             string    `json:"id"`
	BookingID      string    `json:"booking_id"`
	UserID         string    `json:"user_id"`
	Amount         float64   `json:"amount"`
	StripeIntentID string    `json:"stripe_intent_id"`
	Status         string    `json:"status"` // PENDING, PAID, FAILED
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type PaymentRepository interface {
	CreateTransaction(ctx context.Context, bookingID, userID string, amount float64, stripeIntentID, status string) (*Transaction, error)
	UpdateTransactionStatusByIntent(ctx context.Context, stripeIntentID, status string) (*Transaction, error)
}
