proto:
	python3 -m grpc_tools.protoc -I protobuf --python_out=pb --grpc_python_out=pb protobuf/vectors.proto