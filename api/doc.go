package api

import (
	"context"
	"fmt"
	"github.com/cockroachdb/errors"
	"github.com/gin-gonic/gin"
	"net/http"
	db "semantic_api/db/sqlc"
	"semantic_api/token"
	"semantic_api/vectorEmbeddingAPI"
	"semantic_api/vector_db"
)

type createDocRequest struct {
	Doc string `json:"doc" binding:"required"`
}

const collectionName = "people_docs5"

func (server *Server) CreateDoc(ctx *gin.Context) {
	var req *createDocRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	userId, err := server.store.GetUserID(ctx, authPayload.Username)

	arg := db.CreateDocParams{
		UserID: userId,
		Doc:    req.Doc,
	}

	user, err := server.store.CreateDoc(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = vector_db.AddToVectorDB(server.vectorOp, req.Doc, server.config.OpenAIAPIKey, server.config.OpenAIURL, userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	ctx.JSON(http.StatusOK, user)
}

type GetDocRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) GetDoc(ctx *gin.Context) {
	var req GetDocRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	doc, err := server.store.GetDoc(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, doc)
}

type SearchSimilarDocsRequest struct {
	QueryDoc string `json:"query_doc" binding:"required"`
}

func (server *Server) SearchSimilarDocs(ctx *gin.Context) {
	var req *SearchSimilarDocsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// TODO handle
	if milvusOp, ok := server.vectorOp.(*vector_db.MilvusVectorOp); ok {
		// server.vectorOp is of type *vector_db.MilvusVectorOp
		has, err := (*(milvusOp.MilvusClient)).HasCollection(context.Background(), collectionName)
		if err != nil {
			fmt.Errorf("failed to get Has collection %w", err.Error())
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		}

		if !has {
			err := milvusOp.CreateColl()
			if err != nil {
				fmt.Errorf("failed to create collection %w", err.Error())
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			}

			// TODO handle
			//err = createIndex(server.milvusClient)
			if err != nil {
				fmt.Errorf("failed to create index %w", err.Error())
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			}
		}
	} else {
		fmt.Println("milvusOp incorrect")
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("milvusOp incorrect")))
	}

	// get queryDoc as vector
	queryVector, err := vectorEmbeddingAPI.GetVectorEmbedding(req.QueryDoc, server.config.OpenAIAPIKey, server.config.OpenAIURL)
	if err != nil {
		fmt.Errorf("error getting query vector: %v", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	// search in milvusdb
	similarDocsIds, err := server.vectorOp.SearchInDb(queryVector)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	// get docs
	similarDocs, err := server.store.GetDocs(ctx, similarDocsIds)
	if err != nil {
		fmt.Errorf("failed to get similar docs %w", err.Error())
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	ctx.JSON(http.StatusOK, similarDocs)
}
