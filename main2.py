import grpc
import vectors_pb2  # Import the generated Python code
import vectors_pb2_grpc  # Import the generated gRPC code

def run_grpc_client():
    channel = grpc.insecure_channel('localhost:50051')  # Connect to the gRPC server
    stub = vectors_pb2_grpc.VectorManagerStub(channel)

    # Create a request message
    request = vectors_pb2.AddVectorRequest(data='your_data')
    # Call the remote gRPC method
    response = stub.MyMethod(request)
    print("Received: " + response.result)

if __name__ == '__main__':
    run_grpc_client()