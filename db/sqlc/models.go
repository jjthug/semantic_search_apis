// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.23.0

package db

import (
	"time"
)

type Doc struct {
	UserID int64  `json:"user_id"`
	Doc    string `json:"doc"`
}

type User struct {
	UserID         int64     `json:"user_id"`
	Username       string    `json:"username"`
	HashedPassword string    `json:"hashed_password"`
	CreatedAt      time.Time `json:"created_at"`
}
