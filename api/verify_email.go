package api

import (
	"errors"
	"fmt"
	"net/http"
	db "semantic_api/db/sqlc"
	"semantic_api/val"
	"strconv"

	"github.com/gin-gonic/gin"
)

type verifyEmailRequest struct {
	EmailId    int64
	SecretCode string
}

type verifyEmailResponse struct {
	Username        string `json:"username"`
	Email           string `json:"email"`
	IsEmailVerified bool   `json:"is_email_verified"`
}

func (server *Server) VerifyEmail(ctx *gin.Context) {
	id_string := ctx.Query("id")           // Equivalent to ctx.Request.URL.Query().Get("id")
	secretCode := ctx.Query("secret_code") // Equivalent to ctx.Request.URL.Query().Get("secret_code")

	id, err := strconv.ParseInt(id_string, 10, 64)
	if err != nil {
		// Handle error if the conversion fails
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	req := &verifyEmailRequest{
		EmailId:    int64(id),
		SecretCode: secretCode,
	}

	err = validateVerifyEmail(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("failed verify email validation")))
	}

	arg := db.VerifyEmailTxParams{
		EmailId:    req.EmailId,
		SecretCode: req.SecretCode,
		VectorOp:   &server.vectorOp,
		URL:        server.config.OpenAIURL,
		APIKEy:     server.config.OpenAIAPIKey,
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

	verifyEmailResponse := &verifyEmailResponse{
		Username:        user.User.Username,
		Email:           user.User.Email,
		IsEmailVerified: user.User.IsEmailVerified,
	}

	ctx.JSON(http.StatusOK, verifyEmailResponse)
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
