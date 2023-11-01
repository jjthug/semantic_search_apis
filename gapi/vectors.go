package gapi

import (
	"api_with_milvus/pb"
	"api_with_milvus/utils"
	"context"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

func (server *Server) AddVector(context.Context, *pb.AddVectorRequest) (*pb.AddVectorResponse, error) {
	//return nil, status.Errorf(codes.Unimplemented, "method AddVector not implemented")
	vectors := utils.CreateNRandomVectorsDimM(20, 256)
	collectionName := "book"
	bookIDs := make([]int64, 0, 20)

	for i := 0; i < 20; i++ {
		bookIDs = append(bookIDs, int64(i))
	}

	idColumn := entity.NewColumnInt64("book_id", bookIDs)
	vectorColumn := entity.NewColumnFloatVector("word_count", 256, vectors)

	_, err := (*server.client).Insert(
		context.Background(),
		collectionName,
		"",
		idColumn,
		vectorColumn,
	)

	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &pb.AddVectorResponse{
		Resp: true,
	}, nil

}
func (server *Server) SearchVector(context.Context, *pb.SearchVectorRequest) (*pb.SearchVectorResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SearchVector not implemented")
}
