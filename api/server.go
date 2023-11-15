package api

import (
	"github.com/gin-gonic/gin"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	db "semantic_api/db/sqlc"
	"semantic_api/pb"
)

type Server struct {
	store        db.Store
	router       *gin.Engine
	grpcClient   *pb.VectorManagerClient
	milvusClient *client.Client
}

func NewServer(store db.Store, client *pb.VectorManagerClient, milvusClient *client.Client) *Server {
	server := &Server{
		store:        store,
		grpcClient:   client,
		milvusClient: milvusClient,
	}

	router := gin.Default()
	router.POST("/create_user", server.CreateNewUser)
	router.GET("/get_user/:id", server.GetUser)
	router.POST("/create_doc", server.CreateDoc)
	router.GET("/get_doc/:id", server.GetDoc)

	server.router = router
	return server
}

func (server *Server) RunHTTPServer(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
