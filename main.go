package main

import (
	"api_with_milvus/gapi"
	"api_with_milvus/pb"
	"context"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func main() {
	//r := chi.NewRouter()
	//
	// Initialize Milvus client
	client, err := client.NewClient(context.Background(), client.Config{
		Address: "localhost:19530",
	})
	if err != nil {
		// handle error
	}
	defer client.Close()
	//
	//// Define your API route to add a vector
	////r.Post("/addVector", func(w http.ResponseWriter, r *http.Request) {
	////	// Parse the vector from the request body
	////	// You need to implement this part based on your request format
	////
	////	// For example, if your request sends a JSON object with a "vector" field
	////	// you can parse it like this:
	////	// var requestData struct {
	////	//     Vector []float32 `json:"vector"`
	////	// }
	////	// decoder := json.NewDecoder(r.Body)
	////	// if err := decoder.Decode(&requestData); err != nil {
	////	//     http.Error(w, "Invalid request body", http.StatusBadRequest)
	////	//     return
	////	// }
	////	// vector := requestData.Vector
	////
	////	// Add the vector to Milvus
	////	vector := []float32{1.0, 2.0, 3.0} // Replace with your vector data
	////	vectorIDs, err := client.Insert(context.Background(), "my_collection", []milvus.VecField{"my_vector_field"}, []milvus.Vector{vector})
	////	if err != nil {
	////		// Handle the error
	////		http.Error(w, "Failed to add the vector", http.StatusInternalServerError)
	////		return
	////	}
	////
	////	// You can return the vector IDs as a response
	////	w.Write([]byte("Vector added with ID: " + vectorIDs[0]))
	////})
	//
	//r.Post("/createColl", func(w http.ResponseWriter, r *http.Request) {
	//	if ok, _ := client.HasCollection(context.Background(), "book"); ok {
	//		w.Write([]byte("Collection already exists"))
	//		return
	//	}
	//	err := createColl(&client)
	//	if err != nil {
	//		log.Println(err)
	//		w.Write([]byte("Failed"))
	//	} else {
	//		w.Write([]byte("Created Collection"))
	//	}
	//})
	//
	//r.Post("/addRandomVecs", func(w http.ResponseWriter, r *http.Request) {
	//	vectors := utils.CreateNRandomVectorsDimM(20, 256)
	//	collectionName := "book"
	//	bookIDs := make([]int64, 0, 20)
	//
	//	for i := 0; i < 20; i++ {
	//		bookIDs = append(bookIDs, int64(i))
	//	}
	//
	//	idColumn := entity.NewColumnInt64("book_id", bookIDs)
	//	vectorColumn := entity.NewColumnFloatVector("word_count", 256, vectors)
	//
	//	_, err := client.Insert(
	//		context.Background(),
	//		collectionName,
	//		"",
	//		idColumn,
	//		vectorColumn,
	//	)
	//
	//	if err != nil {
	//		log.Println(err)
	//		w.Write([]byte("Failed to add vecs"))
	//	}
	//	w.Write([]byte("Added vectors"))
	//})
	//
	//r.Get("/hasPartition", func(w http.ResponseWriter, r *http.Request) {
	//	has, err := client.HasCollection(context.Background(), "book")
	//	if err != nil {
	//		w.Write([]byte("error"))
	//	}
	//	if has {
	//		w.Write([]byte("1"))
	//	} else {
	//		w.Write([]byte("0"))
	//	}
	//})
	//
	//http.ListenAndServe(":8080", r)

	runGrpcServer(&client)
}

func runGrpcServer(client *client.Client) {

	server, err := gapi.NewServer(client)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	grpcServer := grpc.NewServer()
	//server,err :=
	pb.RegisterVectorManagerServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatal("cannot create listener:", err)
	}

	log.Printf("start gRpc server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start gRpc server:", err)
	}

}
