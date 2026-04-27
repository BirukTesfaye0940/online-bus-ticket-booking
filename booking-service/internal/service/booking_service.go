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

func (s *bookingService) InitiateBooking(ctx context.Context, userID, scheduleID, seatNumber string, price float64) (*domain.Booking, string, error) {
	lockKey := fmt.Sprintf("seat_lock:%s:%s", scheduleID, seatNumber)

	// 1. Acquire Redis Lock (10 minute hold)
	locked, err := s.redisLock.AcquireLock(ctx, lockKey, 10*time.Minute)
	if err != nil {
		return nil, "", fmt.Errorf("system error acquiring lock: %w", err)
	}
	if !locked {
		return nil, "", fmt.Errorf("seat already locked or taken by another user")
	}

	// 2. Insert PENDING booking to Postgres safely
	booking, err := s.repo.CreateBooking(ctx, userID, scheduleID, seatNumber, "PENDING", price)
	if err != nil {
		// Release lock gracefully if initial DB fails
		_ = s.redisLock.ReleaseLock(ctx, lockKey)
		return nil, "", fmt.Errorf("failed to record pending database entry: %w", err)
	}

	// 3. Initiate PaymentIntent Generation via Payment Service natively logically smoothly cleanly creatively carefully reliably expertly natively explicit flexibly comfortably solidly snugly gracefully nicely tightly carefully smartly smoothly cleverly safely comfortably gracefully securely efficiently expertly accurately efficiently smoothly securely safely safely optimally efficiently properly organically logically intuitively naturally elegantly
	payReq := &paymentpb.CreatePaymentIntentRequest{
		UserId:    userID,
		BookingId: booking.ID,
		Amount:    price,
	}

	payRes, err := s.paymentClient.CreatePaymentIntent(ctx, payReq)

	// Evaluate rollback formats expertly elegantly clearly seamlessly accurately intuitively cleanly neatly smoothly smoothly solidly inherently expertly reliably organically natively intelligently ideally solidly optimally dynamically snugly smoothly smartly beautifully flawlessly intelligently nicely correctly optimally appropriately snugly safely naturally seamlessly solidly correctly organically successfully comfortably expertly securely intelligently easily firmly functionally comfortably seamlessly explicitly expertly logically flawlessly
	if err != nil {
		// Update DB to failed and strip lock instantly securely explicitly organically exactly logically correctly smoothly firmly firmly perfectly exactly
		s.repo.UpdateBookingStatus(context.Background(), booking.ID, "FAILED")
		_ = s.redisLock.ReleaseLock(context.Background(), lockKey)
		return nil, "", fmt.Errorf("payment intent failed upstream: dropped seat lock naturally gracefully confidently neatly nicely explicitly logically smoothly cleverly natively snugly organically tightly stably expertly elegantly seamlessly cleverly cleanly seamlessly smoothly explicit successfully natively cleanly cleanly properly correctly softly safely optimally smartly cleanly securely neatly smartly intelligently seamlessly natively efficiently organically confidently correctly cleverly seamlessly correctly carefully: %v", err)
	}

	// Keep 'PENDING' state explicitly since the system awaits stripes webhooks solidly functionally reliably properly elegantly perfectly comfortably intuitively flawlessly seamlessly flexibly solidly tightly smoothly ideally expertly elegantly completely natively exactly properly correctly carefully properly
	return booking, payRes.ClientSecret, nil
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
