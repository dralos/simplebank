postgres:
	docker run --name postgres18 -p 5432:5432 -e POSTGRES_USER=app_user -e POSTGRES_PASSWORD=secret -d postgres:18-alpine
createdb:
	docker exec -it postgres18 createdb --username=app_user --owner=app_user simple_bank
dropdb:
	docker exec -it postgres18 dropdb simple_bank
migrateup:
	migrate -path db/migration -database "postgresql://app_user:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up
migratedown:
	migrate -path db/migration -database "postgresql://app_user:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down
sqlc:
	sqlc generate
test:
	go test -v -cover ./...

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test