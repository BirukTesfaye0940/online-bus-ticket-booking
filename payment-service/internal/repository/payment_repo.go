package repository

import (
	"context"

	"github.com/biruk/bus-ticket/payment-service/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type paymentRepository struct {
	queries *Queries
}

func NewPaymentRepository(q *Queries) domain.PaymentRepository {
	return &paymentRepository{queries: q}
}

func (r *paymentRepository) CreateTransaction(ctx context.Context, bookingID, userID string, amount float64, stripeIntentID, status string) (*domain.Transaction, error) {
	numericAmount := pgtype.Numeric{}
	_ = numericAmount.Scan(float32(amount))

	t, err := r.queries.CreateTransaction(ctx, CreateTransactionParams{
		BookingID:      r.stringToUUID(bookingID),
		UserID:         r.stringToUUID(userID),
		Amount:         numericAmount,
		StripeIntentID: stripeIntentID,
		Status:         status,
	})
	if err != nil {
		return nil, err
	}

	val, _ := t.Amount.Float64Value()

	return &domain.Transaction{
		ID:             r.uuidToString(t.ID),
		BookingID:      r.uuidToString(t.BookingID),
		UserID:         r.uuidToString(t.UserID),
		Amount:         val.Float64,
		StripeIntentID: t.StripeIntentID,
		Status:         t.Status,
		CreatedAt:      t.CreatedAt.Time,
		UpdatedAt:      t.UpdatedAt.Time,
	}, nil
}

func (r *paymentRepository) UpdateTransactionStatusByIntent(ctx context.Context, stripeIntentID, status string) (*domain.Transaction, error) {
	t, err := r.queries.UpdateTransactionStatusByIntent(ctx, UpdateTransactionStatusByIntentParams{
		StripeIntentID: stripeIntentID,
		Status:         status,
	})
	if err != nil {
		return nil, err
	}

	val, _ := t.Amount.Float64Value()

	return &domain.Transaction{
		ID:             r.uuidToString(t.ID),
		BookingID:      r.uuidToString(t.BookingID),
		UserID:         r.uuidToString(t.UserID),
		Amount:         val.Float64,
		StripeIntentID: t.StripeIntentID,
		Status:         t.Status,
		CreatedAt:      t.CreatedAt.Time,
		UpdatedAt:      t.UpdatedAt.Time,
	}, nil
}

func (r *paymentRepository) uuidToString(id pgtype.UUID) string {
	if !id.Valid {
		return ""
	}
	return uuid.UUID(id.Bytes).String()
}

func (r *paymentRepository) stringToUUID(str string) pgtype.UUID {
	u, err := uuid.Parse(str)
	if err != nil {
		return pgtype.UUID{Valid: false}
	}
	return pgtype.UUID{Bytes: u, Valid: true}
}
