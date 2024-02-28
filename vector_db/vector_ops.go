package vector_db

type VectorOp interface {
	AddToDb(...interface{}) error
	SearchInDb(...interface{}) ([]int64, error)
}
