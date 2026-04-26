package repository

import (
	"context"

	"github.com/biruk/bus-ticket/booking-service/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type bookingRepository struct {
	queries *Queries
}

func NewBookingRepository(q *Queries) domain.BookingRepository {
	return &bookingRepository{queries: q}
}

func (r *bookingRepository) CreateBooking(ctx context.Context, userID, scheduleID, seatNumber, status string, price float64) (*domain.Booking, error) {
	numericPrice := pgtype.Numeric{}
	// Basic float implementation explicitly bound
	_ = numericPrice.Scan(float32(price))

	b, err := r.queries.CreateBooking(ctx, CreateBookingParams{
		UserID:     r.stringToUUID(userID),
		ScheduleID: r.stringToUUID(scheduleID),
		SeatNumber: seatNumber,
		Status:     status,
		Price:      numericPrice,
	})
	if err != nil {
		return nil, err
	}

	p, _ := b.Price.Float64Value()

	return &domain.Booking{
		ID:         r.uuidToString(b.ID),
		UserID:     r.uuidToString(b.UserID),
		ScheduleID: r.uuidToString(b.ScheduleID),
		SeatNumber: b.SeatNumber,
		Status:     b.Status,
		Price:      p.Float64,
		CreatedAt:  b.CreatedAt.Time,
		UpdatedAt:  b.UpdatedAt.Time,
	}, nil
}

func (r *bookingRepository) UpdateBookingStatus(ctx context.Context, bookingID string, status string) (*domain.Booking, error) {
	b, err := r.queries.UpdateBookingStatus(ctx, UpdateBookingStatusParams{
		ID:     r.stringToUUID(bookingID),
		Status: status,
	})
	if err != nil {
		return nil, err
	}

	p, _ := b.Price.Float64Value()

	return &domain.Booking{
		ID:         r.uuidToString(b.ID),
		UserID:     r.uuidToString(b.UserID),
		ScheduleID: r.uuidToString(b.ScheduleID),
		SeatNumber: b.SeatNumber,
		Status:     b.Status,
		Price:      p.Float64,
		CreatedAt:  b.CreatedAt.Time,
		UpdatedAt:  b.UpdatedAt.Time,
	}, nil
}

func (r *bookingRepository) ListUserBookings(ctx context.Context, userID string, limit, offset int32) ([]*domain.Booking, error) {
	bookings, err := r.queries.ListUserBookings(ctx, ListUserBookingsParams{
		UserID: r.stringToUUID(userID),
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	res := make([]*domain.Booking, len(bookings))
	for i, b := range bookings {
		p, _ := b.Price.Float64Value()
		res[i] = &domain.Booking{
			ID:         r.uuidToString(b.ID),
			UserID:     r.uuidToString(b.UserID),
			ScheduleID: r.uuidToString(b.ScheduleID),
			SeatNumber: b.SeatNumber,
			Status:     b.Status,
			Price:      p.Float64,
			CreatedAt:  b.CreatedAt.Time,
			UpdatedAt:  b.UpdatedAt.Time,
		}
	}

	return res, nil
}

func (r *bookingRepository) uuidToString(id pgtype.UUID) string {
	if !id.Valid {
		return ""
	}
	return uuid.UUID(id.Bytes).String()
}

func (r *bookingRepository) stringToUUID(str string) pgtype.UUID {
	u, err := uuid.Parse(str)
	if err != nil {
		return pgtype.UUID{Valid: false}
	}
	return pgtype.UUID{Bytes: u, Valid: true}
}
