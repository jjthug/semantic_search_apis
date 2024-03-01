package vector_db

import (
	"fmt"
	"semantic_api/gapi"
	"semantic_api/pb"
	"time"
)

type VectorOp interface {
	AddToDb(userId int64, docVector []float32) error
	SearchInDb(queryVector []float32) ([]int64, error)
}

func AddToVectorDB(grpcClient *pb.VectorManagerClient, vectorOp VectorOp, userId int64, doc string) error {
	// get doc converted to vector from grpc server
	start := time.Now()

	docVector, err := gapi.GetDocAsVector(doc, grpcClient)
	if err != nil {
		fmt.Errorf("failed to get doc as vector %v", err.Error())
		return err
	}

	fmt.Println("GetVectorEmbedding from embedding server =>", time.Now().Sub(start))

	start = time.Now()

	// add to vector db
	err = vectorOp.AddToDb(userId, docVector)
	fmt.Println("AddToDb zilliz =>", time.Now().Sub(start))

	return err
}
