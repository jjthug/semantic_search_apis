-- name: CreateUser :one
INSERT INTO users (username,hashed_password,created_at) VALUES ($1,$2,$3) RETURNING *;


-- name: GetUser :one
SELECT * FROM users WHERE user_id=$1;


-- name: UpdateUserDescription :one
UPDATE docs
SET doc = $2
WHERE user_id = $1
RETURNING *;


-- name: DeleteUser :exec
DELETE FROM users WHERE user_id=$1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY id
LIMIT $1
OFFSET $2;
