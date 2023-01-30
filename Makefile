postgres:
	docker run --name postgres12 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=hypersus -p 5432:5432 -d postgres:12-alpine

createdb:
	docker exec postgres12 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec postgres12 dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:hypersus@localhost:5432/simple_bank?sslmode=disable" --verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:hypersus@localhost:5432/simple_bank?sslmode=disable" --verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

clean:
	docker stop postgres12
	docker rm postgres12

.PHONY:
	postgres createdb dropdb migrateup migratedown test