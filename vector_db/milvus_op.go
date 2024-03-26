package vector_db

import (
	"context"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
	"github.com/rs/zerolog/log"
)

type MilvusVectorOp struct {
	MilvusClient   *client.Client
	collectionName string
}

func (milvusOp *MilvusVectorOp) AddToDb(userId int64, docVector []float32) error {
	milvusClient := (*(milvusOp.MilvusClient))

	has, err := milvusClient.HasCollection(context.Background(), milvusOp.collectionName)

	if err != nil {
		log.Error().Msgf("failed to get Has collection %s", err.Error())
		return err
	}

	if !has {
		err := milvusOp.CreateColl()
		if err != nil {
			log.Error().Msgf("failed to create collection %s", err.Error())
			return err
		}

		err = milvusOp.CreateIndex()
		if err != nil {
			log.Error().Msgf("failed to create index %s", err.Error())
			return err
		}
	}

	idColumn := entity.NewColumnInt64("user_id", []int64{userId})
	c := [][]float32{}
	c = append(c, docVector)
	docColumn := entity.NewColumnFloatVector("doc_vector", 768, c)

	_, err = milvusClient.Insert(
		context.Background(),    // ctx
		milvusOp.collectionName, // CollectionName
		"",                      // partitionName
		idColumn,                // columnarData
		docColumn,               // columnarData
	)
	if err != nil {
		//log.Fatal("failed to insert data:", err.Error())
		log.Error().Msgf("failed to insert data: %v", err.Error())
	}

	return err
}

func (milvusOp *MilvusVectorOp) CreateColl() error {

	schema := &entity.Schema{
		CollectionName: milvusOp.collectionName,
		Description:    "People docs",
		Fields: []*entity.Field{
			{
				Name:       "user_id",
				DataType:   entity.FieldTypeInt64,
				PrimaryKey: true,
				AutoID:     false,
			},
			{
				Name:     "doc_vector",
				DataType: entity.FieldTypeFloatVector,
				TypeParams: map[string]string{
					"dim": "768",
				},
			},
		},
		EnableDynamicField: true,
	}

	err := (*(milvusOp.MilvusClient)).CreateCollection(context.Background(), schema, 2)

	if err != nil {
		log.Error().Msgf("failed to create collection: %s", err.Error())
	}

	return err

}

func (milvusOp *MilvusVectorOp) CreateIndex() error {
	idx, err := entity.NewIndexIvfFlat( // NewIndex func
		entity.L2, // metricType
		1024,      // ConstructParams
	)
	if err != nil {
		log.Error().Msgf("fail to create ivf flat index parameter: %s", err.Error())
		return err
	}

	err = (*(milvusOp.MilvusClient)).CreateIndex(
		context.Background(),    // ctx
		milvusOp.collectionName, // CollectionName
		"doc_vector",            // fieldName
		idx,                     // entity.Index
		false,                   // async
	)
	if err != nil {
		log.Error().Msgf("fail to create index: %s", err.Error())
	}
	return err
}

func (milvusOp *MilvusVectorOp) SearchInDb(queryVector []float32) ([]int64, error) {
	// first load collection to memory
	milvusClient := (*(milvusOp.MilvusClient))
	err := milvusClient.LoadCollection(
		context.Background(),    // ctx
		milvusOp.collectionName, // CollectionName
		false,                   // async
	)

	if err != nil {
		log.Error().Msgf("failed to load collection: %s", err.Error())
		return nil, err
	}

	sp, _ := entity.NewIndexIvfFlatSearchParam( // NewIndex*SearchParam func
		10, // searchParam
	)

	opt := client.SearchQueryOptionFunc(func(option *client.SearchQueryOption) {
		option.Limit = 3
		option.Offset = 0
		option.ConsistencyLevel = entity.ClStrong
		option.IgnoreGrowing = false
	})

	searchResult, err := milvusClient.Search(
		context.Background(),    // ctx
		milvusOp.collectionName, // CollectionName
		[]string{},              // partitionNames
		"",                      // expr
		[]string{"user_id"},     // outputFields
		[]entity.Vector{entity.FloatVector(queryVector)}, // vectors
		"doc_vector", // vectorField
		entity.L2,    // metricType
		10,           // topK
		sp,           // sp
		opt,
	)

	if err != nil {
		log.Error().Msgf("fail to search collection: %s", err.Error())
		return nil, err
	}

	log.Info().Msgf("%#v\n", searchResult)

	val1, err := searchResult[0].IDs.GetAsInt64(0)
	if err != nil {
		log.Error().Msgf("failed to release collection: %s", err.Error())
		return nil, err
	}
	// smaller the scores the more similar
	println(searchResult[0].Scores)

	err = milvusClient.ReleaseCollection(
		context.Background(),    // ctx
		milvusOp.collectionName, // CollectionName
	)

	if err != nil {
		log.Error().Msgf("failed to release collection: %s", err.Error())
		return nil, err
	}

	return []int64{val1}, err
}

func NewMilvusVectorOp(milvusClient *client.Client, collectionName string) VectorOp {
	milvusOp := &MilvusVectorOp{
		MilvusClient:   milvusClient,
		collectionName: collectionName,
	}

	return milvusOp
}
