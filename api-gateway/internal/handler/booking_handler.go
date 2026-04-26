package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/biruk/bus-ticket/api-gateway/internal/middleware"
	pb "github.com/biruk/bus-ticket/api-gateway/internal/proto/booking"
)

type BookingHandler struct {
	bookingClient pb.BookingServiceClient
}

func NewBookingHandler(bc pb.BookingServiceClient) *BookingHandler {
	return &BookingHandler{bookingClient: bc}
}

type initiateBookingRequest struct {
	ScheduleId string  `json:"schedule_id"`
	SeatNumber string  `json:"seat_number"`
	Price      float64 `json:"price"`
}

func (h *BookingHandler) InitiateBooking(w http.ResponseWriter, r *http.Request) {
	var req initiateBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid mapping syntax locally"})
		return
	}

	// Assuming a JWT auth middleware drops the authenticated User ID explicitly into context manually
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Authentication token context dropped seamlessly"})
		return
	}

	resp, err := h.bookingClient.InitiateBooking(r.Context(), &pb.InitiateBookingRequest{
		UserId:     userID,
		ScheduleId: req.ScheduleId,
		SeatNumber: req.SeatNumber,
		Price:      req.Price,
	})

	if err != nil {
		// Example of passing gRPC error safely structurally wrapped down.
		writeJSON(w, http.StatusConflict, map[string]string{"error": err.Error()})
		return
	}

	writeProtoJSON(w, http.StatusCreated, resp)
}

func (h *BookingHandler) ListBookings(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")
	
	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)
	if limit <= 0 {
		limit = 10
	}

	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Authentication required securely explicit"})
		return
	}

	resp, err := h.bookingClient.ListUserBookings(r.Context(), &pb.ListUserBookingsRequest{
		UserId: userID,
		Limit:  int32(limit),
		Offset: int32(offset),
	})

	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeProtoJSON(w, http.StatusOK, resp)
}
