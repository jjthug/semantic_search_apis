package api

import (
	"github.com/gin-gonic/gin"
	db "semantic_api/db/sqlc"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{
		store: store,
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
