postgres:
	sudo docker run --name postgres16 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:16

createdb:
	sudo docker exec -it postgres16 createdb --username=root --owner=root users_semantic

dropdb:
	sudo docker exec -it postgres16 dropdb users_semantic

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/users_semantic?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/users_semantic?sslmode=disable" -verbose down

sqlc:
	sqlc generate
server:
	go run main.go

proto:
	rm -f pb/*.go
	protoc --proto_path=protobuf --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
    protobuf/*.proto

.PHONY: postgres createdb dropdb migrateup sqlc migratedown server proto