package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	db "semantic_api/db/sqlc"
)

type createUserRequest struct {
	Name             string `json:"name" binding:"required"`
	Phone            string `json:"phone" binding:"required"`
	Email            string `json:"email" binding:"required"`
	PasswordHash     string `json:"password_hash" binding:"required"`
	PrivateContact   bool   `json:"private_contact" binding:"required"`
	AboutDescription string `json:"about_description" binding:"required"`
}

func (server *Server) CreateNewUser(ctx *gin.Context) {
	var req *createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Name:             req.Name,
		Phone:            req.Phone,
		Email:            req.Email,
		PasswordHash:     req.PasswordHash,
		PrivateContact:   req.PrivateContact,
		AboutDescription: req.AboutDescription,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}
