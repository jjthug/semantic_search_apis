package api

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"net/http"
	db "semantic_api/db/sqlc"
	"semantic_api/pb"
	"semantic_api/token"
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

	err = AddToVectorDB(server.milvusClient, server.grpcClient, req.Doc, userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	ctx.JSON(http.StatusOK, user)
}

func AddToVectorDB(milvusClient *client.Client, grpcClient *pb.VectorManagerClient, doc string, userId int64) error {
	// get doc converted to vector from grpc server

	docVector, err := getDocAsVector(doc, grpcClient)
	if err != nil {
		fmt.Errorf("failed to get doc as vector %w", err.Error())
		return err
	}

	has, err := (*milvusClient).HasCollection(context.Background(), collectionName)

	if err != nil {
		fmt.Errorf("failed to get Has collection %w", err.Error())
		return nil
	}

	if !has {
		err := vector_db.CreateColl(milvusClient, collectionName)
		if err != nil {
			fmt.Errorf("failed to create collection %w", err.Error())
			return nil
		}
	}

	// add to milvusdb
	err = vector_db.AddToDb(milvusClient, userId, docVector, collectionName)

	return err
}

func getDocAsVector(doc string, grpcClient *pb.VectorManagerClient) ([]float32, error) {

	// Call the GetVector method
	fmt.Print("calling grpc server")
	response, err := (*grpcClient).GetVector(context.Background(), &pb.GetVectorRequest{Doc: doc})
	if err != nil {
		fmt.Errorf("Error calling GetVector: %w", err)
		return nil, err
	}

	// Process the response
	fmt.Printf("Vector Data: %v\n", response.DocVector)

	return response.DocVector, nil
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

func createIndex(milvusClient *client.Client) error {
	return vector_db.CreateIndex(milvusClient, collectionName)
}

func (server *Server) SearchSimilarDocs(ctx *gin.Context) {
	var req *SearchSimilarDocsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	has, err := (*server.milvusClient).HasCollection(context.Background(), collectionName)
	if err != nil {
		fmt.Errorf("failed to get Has collection %w", err.Error())
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	if !has {
		err := vector_db.CreateColl(server.milvusClient, collectionName)
		if err != nil {
			fmt.Errorf("failed to create collection %w", err.Error())
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		}

		err = createIndex(server.milvusClient)
		if err != nil {
			fmt.Errorf("failed to create index %w", err.Error())
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		}
	}

	// get queryDoc as vector
	queryVector, err := getDocAsVector(req.QueryDoc, server.grpcClient)
	if err != nil {
		fmt.Errorf("error getting query vector: %v", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	// search in milvusdb
	similarDocsIds, err := vector_db.SearchInDb(server.milvusClient, collectionName, queryVector)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	// get docs
	similarDocs, err := server.store.GetDoc(ctx, similarDocsIds[0])
	if err != nil {
		fmt.Errorf("failed to get similar docs %w", err.Error())
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	ctx.JSON(http.StatusOK, similarDocs)
}
