package vector_db

import (
	"context"
	"fmt"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
	"log"
)

func AddToDb(milvusClient *client.Client, userId int64, docVector []float32, collectionName string) {

	idColumn := entity.NewColumnInt64("user_id", []int64{userId})
	c := [][]float32{}
	c = append(c, docVector)
	docColumn := entity.NewColumnFloatVector("doc_vector", 384, c)

	_, err := (*milvusClient).Insert(
		context.Background(), // ctx
		collectionName,       // CollectionName
		"",                   // partitionName
		idColumn,             // columnarData
		docColumn,            // columnarData
	)
	if err != nil {
		log.Fatal("failed to insert data:", err.Error())
	}
}

func CreateColl(client *client.Client, collectionName string) error {

	schema := &entity.Schema{
		CollectionName: collectionName,
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
					"dim": "384",
				},
			},
		},
		EnableDynamicField: true,
	}

	err := (*client).CreateCollection(context.Background(), schema, 2)

	if err != nil {
		log.Fatal("failed to create collection:", err.Error())
	}

	return err

}

func CreateIndex(milvusClient *client.Client, collectionName string) {
	idx, err := entity.NewIndexIvfFlat( // NewIndex func
		entity.L2, // metricType
		1024,      // ConstructParams
	)
	if err != nil {
		log.Fatal("fail to create ivf flat index parameter:", err.Error())
	}

	err = (*milvusClient).CreateIndex(
		context.Background(), // ctx
		collectionName,       // CollectionName
		"doc_vector",         // fieldName
		idx,                  // entity.Index
		false,                // async
	)
	if err != nil {
		log.Fatal("fail to create index:", err.Error())
	}
}

func SearchInDb(milvusClient *client.Client, collectionName string, queryVector []float32) []int64 {
	// first load collection to memory
	err := (*milvusClient).LoadCollection(
		context.Background(), // ctx
		collectionName,       // CollectionName
		false,                // async
	)

	if err != nil {
		log.Fatal("failed to load collection:", err.Error())
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

	searchResult, err := (*milvusClient).Search(
		context.Background(),                             // ctx
		collectionName,                                   // CollectionName
		[]string{},                                       // partitionNames
		"",                                               // expr
		[]string{"user_id"},                              // outputFields
		[]entity.Vector{entity.FloatVector(queryVector)}, // vectors
		"doc_vector",                                     // vectorField
		entity.L2,                                        // metricType
		10,                                               // topK
		sp,                                               // sp
		opt,
	)

	if err != nil {
		log.Fatal("fail to search collection:", err.Error())
	}

	fmt.Printf("%#v\n", searchResult)

	val1, _ := searchResult[0].IDs.GetAsInt64(0)

	// smaller the scores the more similar
	println(searchResult[0].Scores)

	err = (*milvusClient).ReleaseCollection(
		context.Background(), // ctx
		collectionName,       // CollectionName
	)

	if err != nil {
		log.Fatal("failed to release collection:", err.Error())
	}

	return []int64{val1}
}
