package api

import (
	"github.com/gin-gonic/gin"
	db "semantic_api/db/sqlc"
	"semantic_api/pb"
)

type Server struct {
	store      db.Store
	router     *gin.Engine
	grpcClient *pb.VectorManagerClient
}

func NewServer(store db.Store, client *pb.VectorManagerClient) *Server {
	server := &Server{
		store:      store,
		grpcClient: client,
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
