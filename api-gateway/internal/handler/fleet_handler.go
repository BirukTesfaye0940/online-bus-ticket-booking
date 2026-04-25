package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	pb "github.com/biruk/bus-ticket/api-gateway/internal/proto/fleet"
)

type FleetHandler struct {
	fleetClient pb.FleetServiceClient
}

func NewFleetHandler(fleetClient pb.FleetServiceClient) *FleetHandler {
	return &FleetHandler{fleetClient: fleetClient}
}

// --- Requests ---

type createBusRequest struct {
	PlateNumber  string `json:"plate_number"`
	OperatorName string `json:"operator_name"`
	Capacity     int32  `json:"capacity"`
}

type createRouteRequest struct {
	Origin                string `json:"origin"`
	Destination           string `json:"destination"`
	DistanceKm            int32  `json:"distance_km"`
	EstimatedDurationMins int32  `json:"estimated_duration_mins"`
}

type createScheduleRequest struct {
	RouteId       string  `json:"route_id"`
	BusId         string  `json:"bus_id"`
	DepartureTime string  `json:"departure_time"`
	ArrivalTime   string  `json:"arrival_time"`
	Price         float64 `json:"price"`
	Status        string  `json:"status"`
}

// --- Handlers ---

func (h *FleetHandler) CreateBus(w http.ResponseWriter, r *http.Request) {
	var req createBusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	resp, err := h.fleetClient.CreateBus(r.Context(), &pb.CreateBusRequest{
		PlateNumber:  req.PlateNumber,
		OperatorName: req.OperatorName,
		Capacity:     req.Capacity,
	})
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeProtoJSON(w, http.StatusCreated, resp)
}

func (h *FleetHandler) ListBuses(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")
	
	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)
	if limit <= 0 {
		limit = 10
	}

	resp, err := h.fleetClient.ListBuses(r.Context(), &pb.ListBusesRequest{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *FleetHandler) CreateRoute(w http.ResponseWriter, r *http.Request) {
	var req createRouteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	resp, err := h.fleetClient.CreateRoute(r.Context(), &pb.CreateRouteRequest{
		Origin:                req.Origin,
		Destination:           req.Destination,
		DistanceKm:            req.DistanceKm,
		EstimatedDurationMins: req.EstimatedDurationMins,
	})
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeProtoJSON(w, http.StatusCreated, resp)
}

func (h *FleetHandler) ListRoutes(w http.ResponseWriter, r *http.Request) {
	resp, err := h.fleetClient.ListRoutes(r.Context(), &pb.ListRoutesRequest{
		Limit:  10,
		Offset: 0,
	})
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *FleetHandler) CreateSchedule(w http.ResponseWriter, r *http.Request) {
	var req createScheduleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	resp, err := h.fleetClient.CreateSchedule(r.Context(), &pb.CreateScheduleRequest{
		RouteId:       req.RouteId,
		BusId:         req.BusId,
		DepartureTime: req.DepartureTime,
		ArrivalTime:   req.ArrivalTime,
		Price:         req.Price,
		Status:        req.Status,
	})
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeProtoJSON(w, http.StatusCreated, resp)
}

func (h *FleetHandler) ListSchedules(w http.ResponseWriter, r *http.Request) {
	resp, err := h.fleetClient.ListSchedules(r.Context(), &pb.ListSchedulesRequest{
		Limit:  10,
		Offset: 0,
	})
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeProtoJSON(w, http.StatusOK, resp)
}
