package api

import (
	"errors"
	"fmt"
	"net/http"
	db "semantic_api/db/sqlc"
	"semantic_api/val"

	"github.com/gin-gonic/gin"
)

type verifyEmailRequest struct {
	EmailId    int64
	SecretCode string
}

type verifyEmailResponse struct {
	User db.User
}

func (server *Server) VerifyEmail(ctx *gin.Context) {
	var req *verifyEmailRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	err := validateVerifyEmail(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("failed verify email validation")))
	}

	arg := db.VerifyEmailTxParams{
		EmailId:    req.EmailId,
		SecretCode: req.SecretCode,
	}

	user, err := server.store.VerifyEmailTx(ctx, arg)
	if err != nil {
		errCode := db.ErrorCode(err)
		if errCode == db.UniqueViolation || errCode == db.ForeignKeyViolation {
			ctx.JSON(http.StatusForbidden, errorResponse(err))
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	createUserResponse := &CreateUserResponse{
		Username: user.User.Username,
		FullName: user.User.FullName,
		Email:    user.User.Email,
	}

	ctx.JSON(http.StatusOK, createUserResponse)
}

func validateVerifyEmail(req *verifyEmailRequest) (err error) {
	if err := val.ValidateEmailId(req.EmailId); err != nil {
		return fmt.Errorf("failed email validation: %w", err)
	}
	if err = val.ValidateSecretCode(req.SecretCode); err != nil {
		return fmt.Errorf("failed secret validation: %w", err)
	}
	return nil
}
