// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package db

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Doc struct {
	UserID int64  `json:"user_id"`
	Doc    string `json:"doc"`
}

type Session struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	RefreshToken string    `json:"refresh_token"`
	UserAgent    string    `json:"user_agent"`
	ClientIp     string    `json:"client_ip"`
	IsBlocked    bool      `json:"is_blocked"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
}

type User struct {
	UserID            int64            `json:"user_id"`
	Username          string           `json:"username"`
	HashedPassword    string           `json:"hashed_password"`
	FullName          string           `json:"full_name"`
	Email             string           `json:"email"`
	PasswordChangedAt pgtype.Timestamp `json:"password_changed_at"`
	CreatedAt         time.Time        `json:"created_at"`
}
