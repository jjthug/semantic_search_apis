package main

import (
	"context"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

func createColl(client *client.Client) error {
	//client, err := client.NewClient(context.Background(), client.Config{
	//	Address: "localhost:19530",
	//})
	//if err != nil {
	//	// handle error
	//}
	//defer client.Close()

	var (
		collectionName = "book"
	)

	schema := &entity.Schema{
		CollectionName: collectionName,
		Description:    "Book search",
		Fields: []*entity.Field{
			{
				Name:       "book_id",
				DataType:   entity.FieldTypeInt64,
				PrimaryKey: true,
				AutoID:     false,
			},
			{
				Name:     "word_count",
				DataType: entity.FieldTypeFloatVector,
				TypeParams: map[string]string{
					"dim": "256",
				},
			},
		},
		EnableDynamicField: true,
	}

	err := (*client).CreateCollection(context.Background(), schema, 2)

	return err

}
