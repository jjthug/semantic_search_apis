-- name: CreateDoc :one
INSERT INTO docs (user_id, doc)
VALUES ($1, $2)
ON CONFLICT (user_id) DO UPDATE
SET doc = EXCLUDED.doc
RETURNING *;

-- name: GetDoc :one
SELECT * FROM docs WHERE user_id=$1;

-- name: GetDocs :many
SELECT * FROM docs WHERE user_id IN (@user_ids::bigserial[]);