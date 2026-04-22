-- name: CreateUser :one
INSERT INTO users (
    email,
    password_hash,
    full_name,
    phone_number,
    role
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: UpdateUser :one
UPDATE users
SET
    full_name = COALESCE(sqlc.narg('full_name'), full_name),
    phone_number = COALESCE(sqlc.narg('phone_number'), phone_number),
    role = COALESCE(sqlc.narg('role'), role),
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;
