package api

import (
	"fmt"
	db "semantic_api/db/sqlc"
	"semantic_api/token"
	"semantic_api/util"
	"semantic_api/vector_db"
	"semantic_api/worker"

	"github.com/gin-gonic/gin"
)

type Server struct {
	config          util.Config
	store           db.Store
	tokenMaker      token.Maker
	router          *gin.Engine
	vectorOp        vector_db.VectorOp
	taskDistributor worker.TaskDistributor
}

func NewServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker([]byte(config.TokenSymmetric))
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	// Using Zillis, might change
	vectorOp := vector_db.NewZillisOp(config.VectorDBCollectionName, config.ZillisAPIKey, config.ZillisEndpoint)

	server := &Server{
		config:          config,
		store:           store,
		tokenMaker:      tokenMaker,
		vectorOp:        vectorOp,
		taskDistributor: taskDistributor,
	}

	router := gin.Default()

	router.POST("/user", server.CreateNewUser)
	//router.GET("/get_user/:id", server.GetUser)
	router.POST("/user/login", server.LoginUser)
	router.GET("/verify_email", server.VerifyEmail)
	router.POST("/tokens/renew_access", server.renewAccessToken)

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
