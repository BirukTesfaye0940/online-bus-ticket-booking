-- name: CreateBus :one
INSERT INTO buses (
    plate_number, operator_name, capacity
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetBus :one
SELECT * FROM buses
WHERE id = $1 LIMIT 1;

-- name: ListBuses :many
SELECT * FROM buses
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CreateRoute :one
INSERT INTO routes (
    origin, destination, distance_km, estimated_duration_mins
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetRoute :one
SELECT * FROM routes
WHERE id = $1 LIMIT 1;

-- name: ListRoutes :many
SELECT * FROM routes
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CreateSchedule :one
INSERT INTO schedules (
    route_id, bus_id, departure_time, arrival_time, price, status
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetSchedule :one
SELECT * FROM schedules
WHERE id = $1 LIMIT 1;

-- name: ListSchedules :many
SELECT * FROM schedules
ORDER BY departure_time ASC
LIMIT $1 OFFSET $2;

-- name: CreateSeatLayout :one
INSERT INTO seat_layouts (
    bus_id, seat_number, is_window, is_aisle
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: ListBusSeats :many
SELECT * FROM seat_layouts
WHERE bus_id = $1;
