package gapi

import (
	"context"
	"fmt"
	"semantic_api/pb"
)

func GetDocAsVector(doc string, grpcClient *pb.VectorManagerClient) ([]float32, error) {

	// Call the GetVector method
	fmt.Print("calling grpc server")
	response, err := (*grpcClient).GetVector(context.Background(), &pb.GetVectorRequest{Doc: doc})
	if err != nil {
		fmt.Errorf("error calling GetVector: %w", err)
		return nil, err
	}

	// Process the response
	//fmt.Printf("Vector Data: %v\n", response.DocVector)

	return response.DocVector, nil
}
