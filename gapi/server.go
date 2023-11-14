package gapi

import (
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"semantic_api/pb"
)

type Server struct {
	pb.UnimplementedVectorManagerServer
	client *client.Client
}

func NewServer(client *client.Client) (*Server, error) {
	server := &Server{client: client}

	return server, nil
}
