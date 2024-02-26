// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0

package db

import (
	"context"
)

type Querier interface {
	CreateDoc(ctx context.Context, arg CreateDocParams) (Doc, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteUser(ctx context.Context, userID int64) error
	GetDoc(ctx context.Context, userID int64) (Doc, error)
	GetDocs(ctx context.Context, userIds []int64) ([]Doc, error)
	GetUser(ctx context.Context, username string) (User, error)
	GetUserID(ctx context.Context, username string) (int64, error)
	ListUsers(ctx context.Context, arg ListUsersParams) ([]User, error)
	UpdateUserDescription(ctx context.Context, arg UpdateUserDescriptionParams) (Doc, error)
}

var _ Querier = (*Queries)(nil)
