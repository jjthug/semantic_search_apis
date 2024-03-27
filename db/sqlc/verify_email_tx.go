package db

import (
	"context"
	"semantic_api/vector_db"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog/log"
)

type VerifyEmailTxParams struct {
	EmailId    int64
	SecretCode string
	VectorOp   *vector_db.VectorOp `json:"vector_op"`
	URL        string              `json:"url"`
	APIKEy     string              `json:"api_key"`
}

type VerifyEmailTxResult struct {
	User        User
	VerifyEmail VerifyEmail
}

func (store *SQLStore) VerifyEmailTx(ctx context.Context, arg VerifyEmailTxParams) (VerifyEmailTxResult, error) {
	var result VerifyEmailTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.VerifyEmail, err = q.UpdateVerifyEmail(ctx, UpdateVerifyEmailParams{
			ID:         arg.EmailId,
			SecretCode: arg.SecretCode,
		})
		if err != nil {
			return err
		}

		result.User, err = q.UpdateUser(ctx, UpdateUserParams{
			Username: result.VerifyEmail.Username,
			IsEmailVerified: pgtype.Bool{
				Bool:  true,
				Valid: true,
			},
		})
		if err != nil {
			return err
		}

		startTime := time.Now()
		docs, err := q.GetDoc(ctx, result.User.UserID)
		if err != nil {
			return err
		}
		log.Info().Msgf("GetDoc time =>", time.Since(startTime))

		startTime = time.Now()

		err = vector_db.AddToVectorDB((*(arg.VectorOp)), docs[0].Doc, arg.APIKEy, arg.URL, result.User.UserID)
		if err != nil {
			return err
		}
		log.Info().Msgf("AddToVectorDB time =>", time.Since(startTime))
		return err
	})

	return result, err
}
