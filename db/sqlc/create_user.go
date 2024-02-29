package db

import (
	"context"
	"semantic_api/vector_db"
	"time"
)

// TransferTxParams contains the input parameters of the transfer transaction
type CreateUserTxParams struct {
	Username     string              `json:"username"`
	PasswordHash string              `json:"password_hash"`
	Doc          string              `json:"doc"`
	VectorOp     *vector_db.VectorOp `json:"vector_op"`
	URL          string              `json:"url"`
	APIKEy       string              `json:"api_key"`
}

// TransferTxResult is the result of the transfer transaction
type CreateUserTxResult struct {
	UserID int64 `json:"user_id"`
}

// TransferTx performs a money transfer from one account to the other.
// It creates the transfer, add account entries, and update accounts' balance within a database transaction
func (store *SQLStore) CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error) {
	var result CreateUserTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		user, err := q.CreateUser(ctx, CreateUserParams{
			Username:       arg.Username,
			HashedPassword: arg.PasswordHash,
			CreatedAt:      time.Now().UTC(),
		})
		if err != nil {
			return err
		}

		_, err = q.CreateDoc(ctx, CreateDocParams{
			UserID: user.UserID,
			Doc:    arg.Doc,
		})
		if err != nil {
			return err
		}

		err = vector_db.AddToVectorDB((*(arg.VectorOp)), arg.Doc, arg.APIKEy, arg.URL, user.UserID)
		if err != nil {
			return err
		}

		result.UserID = user.UserID

		return err

	})

	return result, err
}
