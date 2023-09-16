postgres:
	docker run --name simple_bank -e POSTGRES_PASSWORD=secret -e POSTGRES_USER=root -p 5432:5432 -d postgres

createdb:
	docker exec -it simple_bank createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it simple_bank dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@127.0.0.1:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@127.0.0.1:5432/simple_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen --package mockdb --destination db/mock/store.go github.com/zsoltardai/simple_bank/db/sqlc Store

.PHONY: postgres createdb dropdb migrateup migratedown sqlc server mock
