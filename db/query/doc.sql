-- name: CreateDoc :one
INSERT INTO docs (user_id,doc) VALUES ($1,$2) RETURNING *;


-- name: GetDoc :one
SELECT * FROM docs WHERE user_id=$1;

