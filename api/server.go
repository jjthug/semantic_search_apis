package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	db "semantic_api/db/sqlc"
	"semantic_api/pb"
	"semantic_api/token"
	"semantic_api/util"
)

type Server struct {
	config       util.Config
	store        db.Store
	tokenMaker   token.Maker
	router       *gin.Engine
	grpcClient   *pb.VectorManagerClient
	milvusClient *client.Client
}

func NewServer(config util.Config, store db.Store, client *pb.VectorManagerClient, milvusClient *client.Client) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker([]byte(config.TokenSymmetric))
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:       config,
		store:        store,
		grpcClient:   client,
		tokenMaker:   tokenMaker,
		milvusClient: milvusClient,
	}

	router := gin.Default()
	router.POST("/user", server.CreateNewUser)
	//router.GET("/get_user/:id", server.GetUser)
	router.POST("/user/login", server.LoginUser)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	authRoutes.POST("/create_doc", server.CreateDoc)
	//router.GET("/get_doc/:id", server.GetDoc)
	authRoutes.GET("/search_doc", server.SearchSimilarDocs)

	server.router = router
	return server, nil
}

func (server *Server) RunHTTPServer(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
