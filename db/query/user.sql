-- name: CreateUser :one
INSERT INTO users (username,hashed_password,full_name,email,created_at) VALUES ($1,$2,$3,$4,$5) RETURNING *;


-- name: GetUser :one
SELECT * FROM users WHERE username=$1;


-- name: GetUserID :one
SELECT user_id from users WHERE username=$1;

-- name: UpdateUserDescription :one
UPDATE docs
SET doc = $2
WHERE user_id = $1
RETURNING *;


-- name: DeleteUser :exec
DELETE FROM users WHERE user_id=$1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY user_id
LIMIT $1
OFFSET $2;


-- name: UpdateUser :one
UPDATE users
SET
  hashed_password = COALESCE(sqlc.narg(hashed_password), hashed_password),
  password_changed_at = COALESCE(sqlc.narg(password_changed_at), password_changed_at),
  full_name = COALESCE(sqlc.narg(full_name), full_name),
  email = COALESCE(sqlc.narg(email), email),
  is_email_verified = COALESCE(sqlc.narg(is_email_verified), is_email_verified)
WHERE
  username = sqlc.arg(username)
RETURNING *;