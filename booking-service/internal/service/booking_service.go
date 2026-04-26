package service

import (
	"context"
	"fmt"
	"time"

	"github.com/biruk/bus-ticket/booking-service/internal/domain"
	"github.com/biruk/bus-ticket/booking-service/internal/repository"
	paymentpb "github.com/biruk/bus-ticket/booking-service/proto"
)

type bookingService struct {
	repo          domain.BookingRepository
	redisLock     repository.RedisLock
	paymentClient paymentpb.PaymentServiceClient
}

func NewBookingService(repo domain.BookingRepository, rLock repository.RedisLock, pc paymentpb.PaymentServiceClient) domain.BookingService {
	return &bookingService{
		repo:          repo,
		redisLock:     rLock,
		paymentClient: pc,
	}
}

func (s *bookingService) InitiateBooking(ctx context.Context, userID, scheduleID, seatNumber string, price float64) (*domain.Booking, error) {
	lockKey := fmt.Sprintf("seat_lock:%s:%s", scheduleID, seatNumber)

	// 1. Acquire Redis Lock (10 minute hold)
	locked, err := s.redisLock.AcquireLock(ctx, lockKey, 10*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("system error acquiring lock: %w", err)
	}
	if !locked {
		return nil, fmt.Errorf("seat already locked or taken by another user")
	}

	// 2. Insert PENDING booking to Postgres safely
	booking, err := s.repo.CreateBooking(ctx, userID, scheduleID, seatNumber, "PENDING", price)
	if err != nil {
		// Release lock gracefully if initial DB fails
		_ = s.redisLock.ReleaseLock(ctx, lockKey)
		return nil, fmt.Errorf("failed to record pending database entry: %w", err)
	}

	// 3. Initiate Payment Service Checkout
	// In a complete architecture, this might just pass payment intent to a frontend.
	// But enforcing synchronous processing ensures the atomic layout requested perfectly.
	payReq := &paymentpb.ProcessPaymentRequest{
		UserId:    userID,
		BookingId: booking.ID,
		Amount:    price,
		Provider:  "TELEBIRR",
	}

	// We wrap gRPC natively inside our domain logic as required. Wait, we ignore deadline. We use background context.
	payRes, err := s.paymentClient.ProcessPayment(ctx, payReq)

	// Evaluate rollback formats
	if err != nil || !payRes.Success {
		// Update DB to failed and strip lock instantly
		s.repo.UpdateBookingStatus(context.Background(), booking.ID, "FAILED")
		_ = s.redisLock.ReleaseLock(context.Background(), lockKey)
		return nil, fmt.Errorf("payment failed upstream: dropped seat lock naturally")
	}

	// Payment Succeeded explicitly: Confirm booking permanently.
	confirmed, err := s.repo.UpdateBookingStatus(context.Background(), booking.ID, "CONFIRMED")
	if err != nil {
		return nil, err
	}
	
	// We retain the Redis lock natively allowing it to simply expire after 10m naturally 
	// OR clear it securely since DB acts as absolute reality.
	_ = s.redisLock.ReleaseLock(context.Background(), lockKey)

	return confirmed, nil
}

func (s *bookingService) ConfirmBooking(ctx context.Context, bookingID, userID string) (*domain.Booking, error) {
	return s.repo.UpdateBookingStatus(ctx, bookingID, "CONFIRMED")
}

func (s *bookingService) CancelBooking(ctx context.Context, bookingID, userID string) (*domain.Booking, error) {
	return s.repo.UpdateBookingStatus(ctx, bookingID, "CANCELLED")
}

func (s *bookingService) ListUserBookings(ctx context.Context, userID string, limit, offset int32) ([]*domain.Booking, error) {
	return s.repo.ListUserBookings(ctx, userID, limit, offset)
}
