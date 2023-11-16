// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.23.0
// source: doc.sql

package db

import (
	"context"
)

const createDoc = `-- name: CreateDoc :one
INSERT INTO docs (user_id, doc)
VALUES ($1, $2)
ON CONFLICT (user_id) DO UPDATE
SET doc = EXCLUDED.doc
RETURNING user_id, doc
`

type CreateDocParams struct {
	UserID int64  `json:"user_id"`
	Doc    string `json:"doc"`
}

func (q *Queries) CreateDoc(ctx context.Context, arg CreateDocParams) (Doc, error) {
	row := q.db.QueryRow(ctx, createDoc, arg.UserID, arg.Doc)
	var i Doc
	err := row.Scan(&i.UserID, &i.Doc)
	return i, err
}

const getDoc = `-- name: GetDoc :one
SELECT user_id, doc FROM docs WHERE user_id=$1
`

func (q *Queries) GetDoc(ctx context.Context, userID int64) (Doc, error) {
	row := q.db.QueryRow(ctx, getDoc, userID)
	var i Doc
	err := row.Scan(&i.UserID, &i.Doc)
	return i, err
}

const getDocs = `-- name: GetDocs :many
SELECT user_id, doc FROM docs WHERE user_id IN ($1::bigserial[])
`

func (q *Queries) GetDocs(ctx context.Context, userIds []int64) ([]Doc, error) {
	rows, err := q.db.Query(ctx, getDocs, userIds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Doc{}
	for rows.Next() {
		var i Doc
		if err := rows.Scan(&i.UserID, &i.Doc); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
