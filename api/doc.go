package api

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
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

	docVector := getDocAsVector(req.Doc)

	// add to milvusdb
	vector_db.AddToDb(docVector)

	ctx.JSON(http.StatusOK, user)
}

func getDocAsVector(doc string) []float32 {
	// Set up a connection to the server
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Fatal("Failed to close connection", err)
		}
	}(conn)

	// Create a client using the generated code
	client := pb.NewVectorManagerClient(conn)

	// Call the GetVector method
	fmt.Print("calling grpc server")
	response, err := client.GetVector(context.Background(), &pb.GetVectorRequest{Doc: doc})
	if err != nil {
		log.Fatalf("Error calling GetVector: %v", err)
	}

	// Process the response
	fmt.Printf("Vector Data: %v\n", (*response).DocVector)

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
