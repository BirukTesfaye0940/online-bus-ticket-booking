package grpc_transport

import (
	"context"
	"time"

	"github.com/biruk/bus-ticket/fleet-service/internal/domain"
	pb "github.com/biruk/bus-ticket/fleet-service/proto"
)

type fleetHandler struct {
	pb.UnimplementedFleetServiceServer
	service domain.FleetService
}

func NewFleetHandler(service domain.FleetService) pb.FleetServiceServer {
	return &fleetHandler{
		service: service,
	}
}

func (h *fleetHandler) CreateBus(ctx context.Context, req *pb.CreateBusRequest) (*pb.BusResponse, error) {
	bus, err := h.service.CreateBus(ctx, req.PlateNumber, req.OperatorName, req.Capacity)
	if err != nil {
		return nil, err
	}

	return &pb.BusResponse{
		Id:           bus.ID,
		PlateNumber:  bus.PlateNumber,
		OperatorName: bus.OperatorName,
		Capacity:     bus.Capacity,
	}, nil
}

func (h *fleetHandler) ListBuses(ctx context.Context, req *pb.ListBusesRequest) (*pb.ListBusesResponse, error) {
	buses, err := h.service.ListBuses(ctx, req.Limit, req.Offset)
	if err != nil {
		return nil, err
	}

	var res []*pb.BusResponse
	for _, b := range buses {
		res = append(res, &pb.BusResponse{
			Id:           b.ID,
			PlateNumber:  b.PlateNumber,
			OperatorName: b.OperatorName,
			Capacity:     b.Capacity,
		})
	}

	return &pb.ListBusesResponse{Buses: res}, nil
}

func (h *fleetHandler) CreateRoute(ctx context.Context, req *pb.CreateRouteRequest) (*pb.RouteResponse, error) {
	route, err := h.service.CreateRoute(ctx, &domain.Route{
		Origin:                req.Origin,
		Destination:           req.Destination,
		DistanceKm:            req.DistanceKm,
		EstimatedDurationMins: req.EstimatedDurationMins,
	})
	if err != nil {
		return nil, err
	}

	return &pb.RouteResponse{
		Id:                    route.ID,
		Origin:                route.Origin,
		Destination:           route.Destination,
		DistanceKm:            route.DistanceKm,
		EstimatedDurationMins: route.EstimatedDurationMins,
	}, nil
}

func (h *fleetHandler) ListRoutes(ctx context.Context, req *pb.ListRoutesRequest) (*pb.ListRoutesResponse, error) {
	routes, err := h.service.ListRoutes(ctx, req.Limit, req.Offset)
	if err != nil {
		return nil, err
	}

	var res []*pb.RouteResponse
	for _, r := range routes {
		res = append(res, &pb.RouteResponse{
			Id:                    r.ID,
			Origin:                r.Origin,
			Destination:           r.Destination,
			DistanceKm:            r.DistanceKm,
			EstimatedDurationMins: r.EstimatedDurationMins,
		})
	}

	return &pb.ListRoutesResponse{Routes: res}, nil
}

func (h *fleetHandler) CreateSchedule(ctx context.Context, req *pb.CreateScheduleRequest) (*pb.ScheduleResponse, error) {
	dep, _ := time.Parse(time.RFC3339, req.DepartureTime)
	arr, _ := time.Parse(time.RFC3339, req.ArrivalTime)

	sched, err := h.service.CreateSchedule(ctx, &domain.Schedule{
		RouteID:       req.RouteId,
		BusID:         req.BusId,
		DepartureTime: dep,
		ArrivalTime:   arr,
		Price:         req.Price,
		Status:        req.Status,
	})
	if err != nil {
		return nil, err
	}

	return &pb.ScheduleResponse{
		Id:            sched.ID,
		RouteId:       sched.RouteID,
		BusId:         sched.BusID,
		DepartureTime: sched.DepartureTime.Format(time.RFC3339),
		ArrivalTime:   sched.ArrivalTime.Format(time.RFC3339),
		Price:         sched.Price,
		Status:        sched.Status,
	}, nil
}

func (h *fleetHandler) ListSchedules(ctx context.Context, req *pb.ListSchedulesRequest) (*pb.ListSchedulesResponse, error) {
	schedules, err := h.service.ListSchedules(ctx, req.Limit, req.Offset)
	if err != nil {
		return nil, err
	}

	var res []*pb.ScheduleResponse
	for _, s := range schedules {
		res = append(res, &pb.ScheduleResponse{
			Id:            s.ID,
			RouteId:       s.RouteID,
			BusId:         s.BusID,
			DepartureTime: s.DepartureTime.Format(time.RFC3339),
			ArrivalTime:   s.ArrivalTime.Format(time.RFC3339),
			Price:         s.Price,
			Status:        s.Status,
		})
	}

	return &pb.ListSchedulesResponse{Schedules: res}, nil
}

func (h *fleetHandler) CreateSeatLayout(ctx context.Context, req *pb.CreateSeatLayoutRequest) (*pb.SeatLayoutResponse, error) {
	seat, err := h.service.CreateSeatLayout(ctx, &domain.SeatLayout{
		BusID:      req.BusId,
		SeatNumber: req.SeatNumber,
		IsWindow:   req.IsWindow,
		IsAisle:    req.IsAisle,
	})
	if err != nil {
		return nil, err
	}

	return &pb.SeatLayoutResponse{
		Id:         seat.ID,
		BusId:      seat.BusID,
		SeatNumber: seat.SeatNumber,
		IsWindow:   seat.IsWindow,
		IsAisle:    seat.IsAisle,
	}, nil
}

func (h *fleetHandler) ListBusSeats(ctx context.Context, req *pb.ListBusSeatsRequest) (*pb.ListBusSeatsResponse, error) {
	seats, err := h.service.ListBusSeats(ctx, req.BusId)
	if err != nil {
		return nil, err
	}

	var res []*pb.SeatLayoutResponse
	for _, s := range seats {
		res = append(res, &pb.SeatLayoutResponse{
			Id:         s.ID,
			BusId:      s.BusID,
			SeatNumber: s.SeatNumber,
			IsWindow:   s.IsWindow,
			IsAisle:    s.IsAisle,
		})
	}

	return &pb.ListBusSeatsResponse{Seats: res}, nil
}

func (h *fleetHandler) GetManifest(ctx context.Context, req *pb.GetManifestRequest) (*pb.GetManifestResponse, error) {
	seats, err := h.service.GetManifest(ctx, req.ScheduleId)
	if err != nil {
		return nil, err
	}

	return &pb.GetManifestResponse{
		ScheduleId: req.ScheduleId,
		BookedSeats: seats,
	}, nil
}
