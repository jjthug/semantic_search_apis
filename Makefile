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

.PHONY: postgres createdb dropdb migrateup migratedown server