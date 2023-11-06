-- name: CreateUser :one
INSERT INTO users (name, phone, email, password_hash, private_contact, about_description) VALUES ($1,$2,$3,$4,$5,$6) RETURNING *;


-- name: GetUser :one
SELECT * FROM users WHERE id=$1;


-- name: UpdateUserDescription :one
UPDATE users
SET about_description = $1
WHERE id = $2
RETURNING *;


-- name: ListUsers :many
DELETE FROM users WHERE id=$1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY id
LIMIT $1
OFFSET $2;
