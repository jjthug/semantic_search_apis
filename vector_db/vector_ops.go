package vector_db

type VectorOp interface {
	AddToDb(userId int64, docVector []float32) error
	SearchInDb(queryVector []float32) ([]int64, error)
}
