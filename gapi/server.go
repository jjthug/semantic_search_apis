package gapi

import (
	"api_with_milvus/pb"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
)

type Server struct {
	pb.UnimplementedVectorManagerServer
	client *client.Client
}

func NewServer(client *client.Client) (*Server, error) {
	server := &Server{client: client}

	return server, nil
}
