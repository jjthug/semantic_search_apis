package api

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	db "semantic_api/db/sqlc"
	"semantic_api/pb"
	"semantic_api/vector_db"
)

type createDocRequest struct {
	UserId int64  `json:"user_id" binding:"required"`
	Doc    string `json:"doc" binding:"required"`
}

const collectionName = "people_docs"

func (server *Server) CreateDoc(ctx *gin.Context) {
	var req *createDocRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateDocParams{
		UserID: req.UserId,
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

	if !has {
		err := vector_db.CreateColl(server.milvusClient, collectionName)
		if err != nil {
			log.Fatal("failed to create collection", err.Error())
		}
	}

	// add to milvusdb
	vector_db.AddToDb(server.milvusClient, req.UserId, docVector, collectionName)

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
