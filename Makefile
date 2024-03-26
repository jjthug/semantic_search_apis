postgres:
	docker run --name postgres --network=trovi_network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:16.2-alpine

createdb:
	docker exec -it postgres createdb --username=root --owner=root trovi_db

add_migration:
	migrate create -ext sql -dir db/migration -seq add_sessions

dropdb:
	docker exec -it postgres dropdb trovi_db

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/trovi_db?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/trovi_db?sslmode=disable" -verbose down

sqlc:
	sqlc generate

server:
	go run main.go

local:
	docker start c4a7459f5434 fb8211168695 64cec05ba392 6a2da40e8575 39bd4f9330cc 149102c59987

milvus_up:
	docker-compose up -d

build_trovi_image:
	docker build --no-cache -t trovi:latest .

run_trovi:
	docker run --name trovi --network trovi_network -e GIN_MODE=debug -e DB_SOURCE=postgresql://root:secret@postgres:5432/users_semantic?sslmode=disable -p 8080:8080 trovi:latest

find_5432:
	sudo lsof -i :5432

kill_pid:
	sudo kill -9 <pid>

gen_rand_sym_key:
	openssl rand -hex 64 | head -c 32

redis:
	docker run --name redis -p 6379:6379 -d redis:7-alpine

redis_ping:
	docker exec -it redis redis-cli ping

proto:
	rm -f pb/*.go
	protoc --proto_path=protobuf --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
    protobuf/*.proto

.PHONY: postgres createdb dropdb migrateup sqlc migratedown server proto kill_pid find_5432 milvus_up build_trovi_image run_trovi