package vector_db

import (
	"semantic_api/vectorEmbeddingAPI"
	"time"

	"github.com/rs/zerolog/log"
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
		log.Error().Msgf("failed to get doc as vector %v", err.Error())
		return err
	}

	log.Info().Msgf("GetVectorEmbedding from API =>", time.Since(start))

	start = time.Now()

	// add to vector db
	err = vectorOp.AddToDb(userId, docVector)
	log.Info().Msgf("AddToDb zilliz =>", time.Since(start))

	return err
}
