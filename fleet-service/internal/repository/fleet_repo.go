package repository

import (
	"context"
	"fmt"

	"github.com/biruk/bus-ticket/fleet-service/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type fleetRepository struct {
	queries *Queries
}

func NewFleetRepository(q *Queries) domain.FleetRepository {
	return &fleetRepository{
		queries: q,
	}
}

func (r *fleetRepository) CreateBus(ctx context.Context, plate string, operator string, capacity int32) (*domain.Bus, error) {
	bus, err := r.queries.CreateBus(ctx, CreateBusParams{
		PlateNumber:  plate,
		OperatorName: operator,
		Capacity:     capacity,
	})
	if err != nil {
		return nil, err
	}

	return &domain.Bus{
		ID:           r.uuidToString(bus.ID),
		PlateNumber:  bus.PlateNumber,
		OperatorName: bus.OperatorName,
		Capacity:     bus.Capacity,
		CreatedAt:    bus.CreatedAt.Time,
		UpdatedAt:    bus.UpdatedAt.Time,
	}, nil
}

func (r *fleetRepository) ListBuses(ctx context.Context, limit, offset int32) ([]*domain.Bus, error) {
	buses, err := r.queries.ListBuses(ctx, ListBusesParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	res := make([]*domain.Bus, len(buses))
	for i, bus := range buses {
		res[i] = &domain.Bus{
			ID:           r.uuidToString(bus.ID),
			PlateNumber:  bus.PlateNumber,
			OperatorName: bus.OperatorName,
			Capacity:     bus.Capacity,
			CreatedAt:    bus.CreatedAt.Time,
			UpdatedAt:    bus.UpdatedAt.Time,
		}
	}
	return res, nil
}

func (r *fleetRepository) CreateRoute(ctx context.Context, req *domain.Route) (*domain.Route, error) {
	route, err := r.queries.CreateRoute(ctx, CreateRouteParams{
		Origin:                req.Origin,
		Destination:           req.Destination,
		DistanceKm:            req.DistanceKm,
		EstimatedDurationMins: req.EstimatedDurationMins,
	})
	if err != nil {
		return nil, err
	}

	return &domain.Route{
		ID:                    r.uuidToString(route.ID),
		Origin:                route.Origin,
		Destination:           route.Destination,
		DistanceKm:            route.DistanceKm,
		EstimatedDurationMins: route.EstimatedDurationMins,
		CreatedAt:             route.CreatedAt.Time,
		UpdatedAt:             route.UpdatedAt.Time,
	}, nil
}

func (r *fleetRepository) ListRoutes(ctx context.Context, limit, offset int32) ([]*domain.Route, error) {
	routes, err := r.queries.ListRoutes(ctx, ListRoutesParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	res := make([]*domain.Route, len(routes))
	for i, route := range routes {
		res[i] = &domain.Route{
			ID:                    r.uuidToString(route.ID),
			Origin:                route.Origin,
			Destination:           route.Destination,
			DistanceKm:            route.DistanceKm,
			EstimatedDurationMins: route.EstimatedDurationMins,
			CreatedAt:             route.CreatedAt.Time,
			UpdatedAt:             route.UpdatedAt.Time,
		}
	}
	return res, nil
}

func (r *fleetRepository) CreateSchedule(ctx context.Context, req *domain.Schedule) (*domain.Schedule, error) {
	numericPrice := pgtype.Numeric{}
	if err := numericPrice.Scan(fmt.Sprintf("%f", req.Price)); err != nil {
		return nil, err
	}

	schedule, err := r.queries.CreateSchedule(ctx, CreateScheduleParams{
		RouteID:       r.stringToUUID(req.RouteID),
		BusID:         r.stringToUUID(req.BusID),
		DepartureTime: pgtype.Timestamptz{Time: req.DepartureTime, Valid: true},
		ArrivalTime:   pgtype.Timestamptz{Time: req.ArrivalTime, Valid: true},
		Price:         numericPrice,
		Status:        req.Status,
	})
	if err != nil {
		return nil, err
	}

	p, _ := schedule.Price.Float64Value()

	return &domain.Schedule{
		ID:            r.uuidToString(schedule.ID),
		RouteID:       r.uuidToString(schedule.RouteID),
		BusID:         r.uuidToString(schedule.BusID),
		DepartureTime: schedule.DepartureTime.Time,
		ArrivalTime:   schedule.ArrivalTime.Time,
		Price:         p.Float64,
		Status:        schedule.Status,
		CreatedAt:     schedule.CreatedAt.Time,
		UpdatedAt:     schedule.UpdatedAt.Time,
	}, nil
}

func (r *fleetRepository) ListSchedules(ctx context.Context, limit, offset int32) ([]*domain.Schedule, error) {
	schedules, err := r.queries.ListSchedules(ctx, ListSchedulesParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	res := make([]*domain.Schedule, len(schedules))
	for i, schedule := range schedules {
		p, _ := schedule.Price.Float64Value()
		res[i] = &domain.Schedule{
			ID:            r.uuidToString(schedule.ID),
			RouteID:       r.uuidToString(schedule.RouteID),
			BusID:         r.uuidToString(schedule.BusID),
			DepartureTime: schedule.DepartureTime.Time,
			ArrivalTime:   schedule.ArrivalTime.Time,
			Price:         p.Float64,
			Status:        schedule.Status,
			CreatedAt:     schedule.CreatedAt.Time,
			UpdatedAt:     schedule.UpdatedAt.Time,
		}
	}
	return res, nil
}

func (r *fleetRepository) CreateSeatLayout(ctx context.Context, req *domain.SeatLayout) (*domain.SeatLayout, error) {
	seat, err := r.queries.CreateSeatLayout(ctx, CreateSeatLayoutParams{
		BusID:      r.stringToUUID(req.BusID),
		SeatNumber: req.SeatNumber,
		IsWindow:   req.IsWindow,
		IsAisle:    req.IsAisle,
	})
	if err != nil {
		return nil, err
	}

	return &domain.SeatLayout{
		ID:         r.uuidToString(seat.ID),
		BusID:      r.uuidToString(seat.BusID),
		SeatNumber: seat.SeatNumber,
		IsWindow:   seat.IsWindow,
		IsAisle:    seat.IsAisle,
		CreatedAt:  seat.CreatedAt.Time,
	}, nil
}

func (r *fleetRepository) ListBusSeats(ctx context.Context, busID string) ([]*domain.SeatLayout, error) {
	seats, err := r.queries.ListBusSeats(ctx, r.stringToUUID(busID))
	if err != nil {
		return nil, err
	}

	res := make([]*domain.SeatLayout, len(seats))
	for i, seat := range seats {
		res[i] = &domain.SeatLayout{
			ID:         r.uuidToString(seat.ID),
			BusID:      r.uuidToString(seat.BusID),
			SeatNumber: seat.SeatNumber,
			IsWindow:   seat.IsWindow,
			IsAisle:    seat.IsAisle,
			CreatedAt:  seat.CreatedAt.Time,
		}
	}
	return res, nil
}

func (r *fleetRepository) uuidToString(id pgtype.UUID) string {
	if !id.Valid {
		return ""
	}
	return uuid.UUID(id.Bytes).String()
}

func (r *fleetRepository) stringToUUID(str string) pgtype.UUID {
	u, err := uuid.Parse(str)
	if err != nil {
		return pgtype.UUID{Valid: false}
	}
	return pgtype.UUID{Bytes: u, Valid: true}
}
