postgres:
	docker run --name postgres18 -p 5432:5432 -e POSTGRES_USER=app_user -e POSTGRES_PASSWORD=secret -d postgres:18-alpine
createdb:
	docker exec -it postgres18 createdb --username=app_user --owner=app_user simple_bank
dropdb:
	docker exec -it postgres18 dropdb --username=app_user simple_bank 
migrateup:
	migrate -path db/migration -database "postgresql://app_user:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up
migrateup1:
	migrate -path db/migration -database "postgresql://app_user:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up 1
migratedown:
	migrate -path db/migration -database "postgresql://app_user:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down
migratedown1:
	migrate -path db/migration -database "postgresql://app_user:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down 1
sqlc:
	sqlc generate
test:
	go test -v -cover ./...
server:
	go run main.go
mock:
	mockgen --destination db/mock/store.go --package mockdb github.com/dralos/simplebank/db/sqlc Store 

.PHONY: postgres createdb dropdb migrateup migrateup1 migratedown migratedown1 sqlc test server mock