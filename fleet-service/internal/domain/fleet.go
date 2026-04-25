package domain

import (
	"context"
	"time"
)

type Bus struct {
	ID           string    `json:"id"`
	PlateNumber  string    `json:"plate_number"`
	OperatorName string    `json:"operator_name"`
	Capacity     int32     `json:"capacity"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Route struct {
	ID                    string    `json:"id"`
	Origin                string    `json:"origin"`
	Destination           string    `json:"destination"`
	DistanceKm            int32     `json:"distance_km"`
	EstimatedDurationMins int32     `json:"estimated_duration_mins"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

type Schedule struct {
	ID            string    `json:"id"`
	RouteID       string    `json:"route_id"`
	BusID         string    `json:"bus_id"`
	DepartureTime time.Time `json:"departure_time"`
	ArrivalTime   time.Time `json:"arrival_time"`
	Price         float64   `json:"price"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type SeatLayout struct {
	ID         string    `json:"id"`
	BusID      string    `json:"bus_id"`
	SeatNumber string    `json:"seat_number"`
	IsWindow   bool      `json:"is_window"`
	IsAisle    bool      `json:"is_aisle"`
	CreatedAt  time.Time `json:"created_at"`
}

type FleetRepository interface {
	CreateBus(ctx context.Context, plate string, operator string, capacity int32) (*Bus, error)
	ListBuses(ctx context.Context, limit, offset int32) ([]*Bus, error)
	
	CreateRoute(ctx context.Context, req *Route) (*Route, error)
	ListRoutes(ctx context.Context, limit, offset int32) ([]*Route, error)

	CreateSchedule(ctx context.Context, req *Schedule) (*Schedule, error)
	ListSchedules(ctx context.Context, limit, offset int32) ([]*Schedule, error)

	CreateSeatLayout(ctx context.Context, req *SeatLayout) (*SeatLayout, error)
	ListBusSeats(ctx context.Context, busID string) ([]*SeatLayout, error)
}

type FleetService interface {
	CreateBus(ctx context.Context, plate string, operator string, capacity int32) (*Bus, error)
	ListBuses(ctx context.Context, limit, offset int32) ([]*Bus, error)
	
	CreateRoute(ctx context.Context, req *Route) (*Route, error)
	ListRoutes(ctx context.Context, limit, offset int32) ([]*Route, error)

	CreateSchedule(ctx context.Context, req *Schedule) (*Schedule, error)
	ListSchedules(ctx context.Context, limit, offset int32) ([]*Schedule, error)

	CreateSeatLayout(ctx context.Context, req *SeatLayout) (*SeatLayout, error)
	ListBusSeats(ctx context.Context, busID string) ([]*SeatLayout, error)
	
	GetManifest(ctx context.Context, scheduleID string) ([]string, error) // simple placeholder
}
