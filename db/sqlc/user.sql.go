// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: user.sql

package db

import (
	"context"
	"time"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (username,hashed_password,created_at) VALUES ($1,$2,$3) RETURNING user_id, username, hashed_password, created_at
`

type CreateUserParams struct {
	Username       string    `json:"username"`
	HashedPassword string    `json:"hashed_password"`
	CreatedAt      time.Time `json:"created_at"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, createUser, arg.Username, arg.HashedPassword, arg.CreatedAt)
	var i User
	err := row.Scan(
		&i.UserID,
		&i.Username,
		&i.HashedPassword,
		&i.CreatedAt,
	)
	return i, err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM users WHERE user_id=$1
`

func (q *Queries) DeleteUser(ctx context.Context, userID int64) error {
	_, err := q.db.Exec(ctx, deleteUser, userID)
	return err
}

const getUser = `-- name: GetUser :one
SELECT user_id, username, hashed_password, created_at FROM users WHERE user_id=$1
`

func (q *Queries) GetUser(ctx context.Context, userID int64) (User, error) {
	row := q.db.QueryRow(ctx, getUser, userID)
	var i User
	err := row.Scan(
		&i.UserID,
		&i.Username,
		&i.HashedPassword,
		&i.CreatedAt,
	)
	return i, err
}

const listUsers = `-- name: ListUsers :many
SELECT user_id, username, hashed_password, created_at FROM users
ORDER BY id
LIMIT $1
OFFSET $2
`

type ListUsersParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListUsers(ctx context.Context, arg ListUsersParams) ([]User, error) {
	rows, err := q.db.Query(ctx, listUsers, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []User{}
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.UserID,
			&i.Username,
			&i.HashedPassword,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateUserDescription = `-- name: UpdateUserDescription :one
UPDATE docs
SET doc = $2
WHERE user_id = $1
RETURNING user_id, doc
`

type UpdateUserDescriptionParams struct {
	UserID int64  `json:"user_id"`
	Doc    string `json:"doc"`
}

func (q *Queries) UpdateUserDescription(ctx context.Context, arg UpdateUserDescriptionParams) (Doc, error) {
	row := q.db.QueryRow(ctx, updateUserDescription, arg.UserID, arg.Doc)
	var i Doc
	err := row.Scan(&i.UserID, &i.Doc)
	return i, err
}
