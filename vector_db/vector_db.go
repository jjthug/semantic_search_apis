package vector_db

import (
	"context"
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
