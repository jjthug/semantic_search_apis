from concurrent import futures
import grpc
# from vectors_pb2 import GetVectorRequest, GetVectorResponse
import vectors_pb2
from model import encodeToVector
from vectors_pb2_grpc import VectorManagerServicer, add_VectorManagerServicer_to_server
import numpy as np

class VectorManagerServer(VectorManagerServicer):
    def GetVector(self, request, context):
        vector_data = encodeToVector(request.doc)
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