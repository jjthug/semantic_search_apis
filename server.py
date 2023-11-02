# # Import necessary libraries
# from fastapi import FastAPI
# from pydantic import BaseModel
# from sentence_transformers import SentenceTransformer
#
# # Create a FastAPI app
# app = FastAPI()
#
# # Load a pre-trained Sentence Transformers model
# model = SentenceTransformer('paraphrase-MiniLM-L6-v2')
#
#
# # Create a Pydantic model for input data
# class InputData(BaseModel):
#     text: str
#
#
# # Create an endpoint to process the input string
# @app.post("/transform")
# async def transform_text(data: InputData):
#     input_text = data.text
#     # Transform the input text using the Sentence Transformers model
#     embeddings = model.encode([input_text])
#
#     # Return the transformed embeddings
#     return {"embeddings": embeddings[0].tolist()}
#
#
# def run():
#     # NOTE(gRPC Python Team): .close() is possible on a channel and should be
#     # used in circumstances in which the with statement does not fit the needs
#     # of the code.
#     with grpc.insecure_channel("localhost:50051") as channel:
#         stub = route_guide_pb2_grpc.RouteGuideStub(channel)
#         print("-------------- GetFeature --------------")
#         guide_get_feature(stub)
#         print("-------------- ListFeatures --------------")
#         guide_list_features(stub)
#         print("-------------- RecordRoute --------------")
#         guide_record_route(stub)
#         print("-------------- RouteChat --------------")
#         guide_route_chat(stub)
#
# # Run the FastAPI app
# if __name__ == "__main__":
#     import uvicorn
#
#     uvicorn.run(app, host="0.0.0.0", port=8000)
