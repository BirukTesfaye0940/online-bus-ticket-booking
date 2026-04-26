package grpc_transport

import (
	"context"

	"github.com/biruk/bus-ticket/booking-service/internal/domain"
	pb "github.com/biruk/bus-ticket/booking-service/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BookingHandler struct {
	pb.UnimplementedBookingServiceServer
	service domain.BookingService
}

func NewBookingHandler(service domain.BookingService) *BookingHandler {
	return &BookingHandler{service: service}
}

func (h *BookingHandler) InitiateBooking(ctx context.Context, req *pb.InitiateBookingRequest) (*pb.BookingResponse, error) {
	b, err := h.service.InitiateBooking(ctx, req.UserId, req.ScheduleId, req.SeatNumber, req.Price)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "booking failed structurally: %v", err)
	}

	return &pb.BookingResponse{
		Id:         b.ID,
		UserId:     b.UserID,
		ScheduleId: b.ScheduleID,
		SeatNumber: b.SeatNumber,
		Status:     b.Status,
		Price:      b.Price,
		CreatedAt:  b.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

func (h *BookingHandler) ConfirmBooking(ctx context.Context, req *pb.ConfirmBookingRequest) (*pb.BookingResponse, error) {
	b, err := h.service.ConfirmBooking(ctx, req.BookingId, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "confirmation failed: %v", err)
	}

	return &pb.BookingResponse{
		Id:         b.ID,
		UserId:     b.UserID,
		ScheduleId: b.ScheduleID,
		SeatNumber: b.SeatNumber,
		Status:     b.Status,
		Price:      b.Price,
		CreatedAt:  b.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

func (h *BookingHandler) CancelBooking(ctx context.Context, req *pb.CancelBookingRequest) (*pb.BookingResponse, error) {
	b, err := h.service.CancelBooking(ctx, req.BookingId, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cancellation failed: %v", err)
	}

	return &pb.BookingResponse{
		Id:         b.ID,
		UserId:     b.UserID,
		ScheduleId: b.ScheduleID,
		SeatNumber: b.SeatNumber,
		Status:     b.Status,
		Price:      b.Price,
		CreatedAt:  b.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

func (h *BookingHandler) ListUserBookings(ctx context.Context, req *pb.ListUserBookingsRequest) (*pb.ListBookingsResponse, error) {
	bookings, err := h.service.ListUserBookings(ctx, req.UserId, req.Limit, req.Offset)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "query failed: %v", err)
	}

	var res []*pb.BookingResponse
	for _, b := range bookings {
		res = append(res, &pb.BookingResponse{
			Id:         b.ID,
			UserId:     b.UserID,
			ScheduleId: b.ScheduleID,
			SeatNumber: b.SeatNumber,
			Status:     b.Status,
			Price:      b.Price,
			CreatedAt:  b.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return &pb.ListBookingsResponse{Bookings: res}, nil
}
