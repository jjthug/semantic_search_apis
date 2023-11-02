proto:
	python3 -m grpc_tools.protoc -I protobuf --python_out=. --grpc_python_out=. protobuf/vectors.proto