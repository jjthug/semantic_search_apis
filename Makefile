proto:
	rm -f pb/*.go
	protoc --proto_path=protobuf --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
    protobuf/*.proto