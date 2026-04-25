package service

import (
	"context"

	"github.com/biruk/bus-ticket/fleet-service/internal/domain"
)

type fleetService struct {
	repo domain.FleetRepository
}

func NewFleetService(repo domain.FleetRepository) domain.FleetService {
	return &fleetService{
		repo: repo,
	}
}

func (s *fleetService) CreateBus(ctx context.Context, plate string, operator string, capacity int32) (*domain.Bus, error) {
	return s.repo.CreateBus(ctx, plate, operator, capacity)
}

func (s *fleetService) ListBuses(ctx context.Context, limit, offset int32) ([]*domain.Bus, error) {
	return s.repo.ListBuses(ctx, limit, offset)
}

func (s *fleetService) CreateRoute(ctx context.Context, req *domain.Route) (*domain.Route, error) {
	return s.repo.CreateRoute(ctx, req)
}

func (s *fleetService) ListRoutes(ctx context.Context, limit, offset int32) ([]*domain.Route, error) {
	return s.repo.ListRoutes(ctx, limit, offset)
}

func (s *fleetService) CreateSchedule(ctx context.Context, req *domain.Schedule) (*domain.Schedule, error) {
	return s.repo.CreateSchedule(ctx, req)
}

func (s *fleetService) ListSchedules(ctx context.Context, limit, offset int32) ([]*domain.Schedule, error) {
	return s.repo.ListSchedules(ctx, limit, offset)
}

func (s *fleetService) CreateSeatLayout(ctx context.Context, req *domain.SeatLayout) (*domain.SeatLayout, error) {
	return s.repo.CreateSeatLayout(ctx, req)
}

func (s *fleetService) ListBusSeats(ctx context.Context, busID string) ([]*domain.SeatLayout, error) {
	return s.repo.ListBusSeats(ctx, busID)
}

func (s *fleetService) GetManifest(ctx context.Context, scheduleID string) ([]string, error) {
	// A placeholder for now. 
	// To get a full manifest, we might need a join or complex query with Bookings.
	// We'll return an empty list for now until Booking Service connects.
	return []string{}, nil
}
