-- name: CreateBooking :one
INSERT INTO bookings (
    user_id, schedule_id, seat_number, status, price
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: UpdateBookingStatus :one
UPDATE bookings 
SET status = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $1 RETURNING *;

-- name: ListUserBookings :many
SELECT * FROM bookings
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;
