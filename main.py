import grpc
import vectors_pb2  # Import the generated Python code
import vectors_pb2_grpc  # Import the generated gRPC code

def run_grpc_client():
    channel = grpc.insecure_channel('localhost:50051')  # Connect to the gRPC server
    stub = vectors_pb2_grpc.VectorManagerStub(channel)

    # Define the dimension of the vector
    dimension = 256

    # Generate a random 256-dimensional vector
    random_vector = [random.uniform(0, 1) for _ in range(dimension)]

    # Create a request message
    request = vectors_pb2.AddVectorRequest(vectorReq=[1.0,2.3,4.2])
    # Call the remote gRPC method
    response = stub.AddVector(request)
    print("Received: " + response.result)

if __name__ == '__main__':
    run_grpc_client()