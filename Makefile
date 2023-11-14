proto:
	python3 -m grpc_tools.protoc -I protobuf --python_out=. --grpc_python_out=. protobuf/vectors.proto

proto2:
	python -m grpc_tools.protoc --proto_path protobuf/ --python_out=pb --grpc_python_out=pb protobuf/vectors.proto