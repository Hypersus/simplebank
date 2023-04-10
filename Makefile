#!make
include app.env

postgres:
	docker run --name postgres12 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=hypersus -p 5432:5432 -d postgres:12-alpine

createdb:
	docker exec postgres12 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec postgres12 dropdb simple_bank

migrateup:
	migrate -path db/migration -database "${DB_SOURCE}" --verbose up

migrateup1:
	migrate -path db/migration -database "${DB_SOURCE}" --verbose up 1

migratedown:
	migrate -path db/migration -database "${DB_SOURCE}" --verbose down

migratedown1:
	migrate -path db/migration -database "${DB_SOURCE}" --verbose down 1
	
sqlc:
	sqlc generate

test:
	go test -v -cover ./...

clean:
	docker stop postgres12
	docker rm postgres12

server:
	go run main.go

mock:
	mockgen -destination db/mock/store.go -package mockdb github.com/Hypersus/simplebank/db/sqlc Store

.PHONY:
	postgres createdb dropdb migrateup migratedown test server mock