from concurrent import futures
import grpc
# from vectors_pb2 import GetVectorRequest, GetVectorResponse
import vectors_pb2
from vectors_pb2_grpc import VectorManagerServicer, add_VectorManagerServicer_to_server
import numpy as np

class VectorManagerServer(VectorManagerServicer):
    def GetVector(self, request, context):
        # Replace this with your logic to fetch vector data based on the vector_id
        array_size = 10

        # Generate a random float32 array
        vector_data = np.random.rand(array_size).astype(np.float32)
        # vector_data = f"Vector data for ID: {request.vector_id}"
        return vectors_pb2.GetVectorResponse(docVector=vector_data)

def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    add_VectorManagerServicer_to_server(VectorManagerServer(), server)
    server.add_insecure_port('[::]:50051')
    server.start()
    print("Python gRPC server started on port 50051")
    server.wait_for_termination()

if __name__ == '__main__':
    serve()