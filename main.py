import grpc
import vectors_pb2  # Import the generated Python code
import vectors_pb2_grpc  # Import the generated gRPC code
import random
# Import necessary libraries
from fastapi import FastAPI
def generate_random_vector(dimension):
    return [random.uniform(0, 1) for _ in range(dimension)]

def run_grpc_client():
    # Connect to the gRPC server
    channel = grpc.insecure_channel('localhost:50051')
    stub = vectors_pb2_grpc.VectorManagerStub(channel)

    try:
        # Define the dimension of the vector
        dimension = 256

        # Generate a random 256-dimensional vector
        random_vector = generate_random_vector(dimension)

        # Create a request message
        request = vectors_pb2.AddVectorRequest(vectorReq=random_vector)

        # Call the remote gRPC method
        response = stub.AddVector(request)

        # Handle the response
        if response.result:
            print("Received: " + response.result)
        else:
            print("Error: Empty response received")
    except grpc.RpcError as e:
        print(f"Error calling gRPC method: {e.details()}")
    except Exception as e:
        print(f"An unexpected error occurred: {str(e)}")

if __name__ == '__main__':
    run_grpc_client()