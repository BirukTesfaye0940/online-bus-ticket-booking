-- name: CreateTransaction :one
INSERT INTO transactions (
    booking_id, user_id, amount, stripe_intent_id, status
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: UpdateTransactionStatusByIntent :one
UPDATE transactions 
SET status = $2, updated_at = CURRENT_TIMESTAMP
WHERE stripe_intent_id = $1 RETURNING *;
