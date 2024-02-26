package api

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"log"
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

	// get doc converted to vector from grpc server

	docVector := getDocAsVector(req.Doc, server.grpcClient)

	has, err := (*server.milvusClient).HasCollection(context.Background(), collectionName)

	if err != nil {
		log.Fatal("failed to get Has collection", err.Error())
	}

	if !has {
		err := vector_db.CreateColl(server.milvusClient, collectionName)
		if err != nil {
			log.Fatal("failed to create collection", err.Error())
		}
	}

	// add to milvusdb
	vector_db.AddToDb(server.milvusClient, userId, docVector, collectionName)

	ctx.JSON(http.StatusOK, user)
}

func getDocAsVector(doc string, grpcClient *pb.VectorManagerClient) []float32 {

	// Call the GetVector method
	fmt.Print("calling grpc server")
	response, err := (*grpcClient).GetVector(context.Background(), &pb.GetVectorRequest{Doc: doc})
	if err != nil {
		log.Fatalf("Error calling GetVector: %v", err)
	}

	// Process the response
	fmt.Printf("Vector Data: %v\n", response.DocVector)

	return response.DocVector
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

func createIndex(milvusClient *client.Client) {
	vector_db.CreateIndex(milvusClient, collectionName)
}

func (server *Server) SearchSimilarDocs(ctx *gin.Context) {
	var req *SearchSimilarDocsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	has, err := (*server.milvusClient).HasCollection(context.Background(), collectionName)
	if err != nil {
		log.Fatal("failed to get Has collection", err.Error())
	}

	if !has {
		err := vector_db.CreateColl(server.milvusClient, collectionName)
		if err != nil {
			log.Fatal("failed to create collection", err.Error())
		}

		createIndex(server.milvusClient)
	}

	// get queryDoc as vector
	queryVector := getDocAsVector(req.QueryDoc, server.grpcClient)

	// search in milvusdb
	similarDocsIds := vector_db.SearchInDb(server.milvusClient, collectionName, queryVector)

	// get docs
	similarDocs, err := server.store.GetDoc(ctx, similarDocsIds[0])

	if err != nil {
		log.Fatal("failed to get similar docs", err.Error())
	}

	ctx.JSON(http.StatusOK, similarDocs)
}
