package db

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
)

// TransferTxParams contains the input parameters of the transfer transaction
type CreateUserTxParams struct {
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
	FullName     string `json:"full_name"`
	Email        string `json:"email"`
	Doc          string `json:"doc"`
	AfterCreate  func(user User) error
}

// TransferTxResult is the result of the transfer transaction
type CreateUserTxResult struct {
	User User
}

// TransferTx performs a money transfer from one account to the other.
// It creates the transfer, add account entries, and update accounts' balance within a database transaction
func (store *SQLStore) CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error) {
	var result CreateUserTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		startTime := time.Now()
		result.User, err = q.CreateUser(ctx, CreateUserParams{
			Username:       arg.Username,
			HashedPassword: arg.PasswordHash,
			FullName:       arg.FullName,
			Email:          arg.Email,
			CreatedAt:      time.Now().UTC(),
		})
		if err != nil {
			return err
		}
		log.Info().Msgf("Created user time =>", time.Since(startTime))

		_, err = q.CreateDoc(ctx, CreateDocParams{
			UserID: result.User.UserID,
			Doc:    arg.Doc,
		})
		if err != nil {
			return err
		}

		return arg.AfterCreate(result.User)
	})

	return result, err
}
