package domain

import (
	"context"
	"time"
)

type Booking struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	ScheduleID string    `json:"schedule_id"`
	SeatNumber string    `json:"seat_number"`
	Status     string    `json:"status"` // PENDING, CONFIRMED, FAILED
	Price      float64   `json:"price"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type BookingRepository interface {
	CreateBooking(ctx context.Context, userID, scheduleID, seatNumber, status string, price float64) (*Booking, error)
	UpdateBookingStatus(ctx context.Context, bookingID string, status string) (*Booking, error)
	ListUserBookings(ctx context.Context, userID string, limit, offset int32) ([]*Booking, error)
}

type BookingService interface {
	InitiateBooking(ctx context.Context, userID, scheduleID, seatNumber string, price float64) (*Booking, string, error)
	ConfirmBooking(ctx context.Context, bookingID, userID string) (*Booking, error)
	CancelBooking(ctx context.Context, bookingID, userID string) (*Booking, error)
	ListUserBookings(ctx context.Context, userID string, limit, offset int32) ([]*Booking, error)
}
