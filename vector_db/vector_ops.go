package vector_db

import (
	"fmt"
	"semantic_api/vectorEmbeddingAPI"
	"time"
)

type VectorOp interface {
	AddToDb(userId int64, docVector []float32) error
	SearchInDb(queryVector []float32) ([]int64, error)
}

func AddToVectorDB(vectorOp VectorOp, doc, apiKey, url string, userId int64) error {
	// get doc converted to vector from grpc server
	start := time.Now()

	docVector, err := vectorEmbeddingAPI.GetVectorEmbedding(doc, apiKey, url)
	if err != nil {
		fmt.Errorf("failed to get doc as vector %v", err.Error())
		return err
	}

	fmt.Println("GetVectorEmbedding from API =>", time.Now().Sub(start))

	start = time.Now()

	// add to vector db
	err = vectorOp.AddToDb(userId, docVector)
	fmt.Println("AddToDb zilliz =>", time.Now().Sub(start))

	return err
}
