DB_URL=postgresql://root:secret@localhost:5432/yoyo_store?sslmode=disable

network:
	docker network create yoyo-network

postgres:
	docker run --name postgres --network yoyo-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:17-alpine

createdb:
	docker exec -it postgres createdb --username=root --owner=root yoyo_store

dropdb:
	docker exec -it postgres dropdb yoyo_store

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migrateup1:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

migratedown1:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1

new_migration:
	migrate create -ext sql -dir db/migration -seq $(name)

## build: builds all binaries
build: clean build_back build_invoice build_front
	@printf "All binaries built!\n"

## clean: cleans all binaries and runs go clean
clean:
	@echo "Cleaning..."
	@- rm -f dist/*
	@go clean
	@echo "Cleaned!"

## build_back: builds the back end
build_back:
	@echo "Building back end..."
	@go build -o dist/server_main ./server_main
	@echo "Back end built!"

## build_invoice: builds the invoice microservice
build_invoice:
	@echo "Building invoice microservice..."
	@go build -o dist/server_invoice ./server_invoice
	@echo "Invoice microservice built!"

## build_front: builds the front end
build_front:
	@echo "Building front end..."
	@go build -o dist/frontend ./frontend
	@echo "Front end built!"

## start_built: starts front, invoice and back end
start_built: start_built_back start_built_invoice start_built_front

## start_built_back: starts the back end built
start_built_back: build_back
	@echo "Starting the back end..."
	@./dist/server_main &
	@echo "Back end running!"

## start_built_invoice: starts the invoice microservice built
start_built_invoice: build_invoice
	@echo "Starting the invoice microservice..."
	@./dist/server_invoice &
	@echo "Invoice microservice running!"

## start_built_front: starts the front end built
start_built_front: build_front
	@echo "Starting the front end..."
	@./dist/frontend &
	@echo "Front end running!"

## stop: stops the front and back end
stop_built: stop_front stop_back stop_invoice
	@echo "All applications stopped"

## stop_front: stops the front end
stop_built_front:
	@echo "Stopping the front end..."
	@-pkill -SIGTERM -f "frontend"
	@echo "Stopped front end"

## stop_invoice: stops the invoice microservice
stop_built_invoice:
	@echo "Stopping the invoice microservice..."
	@-pkill -SIGTERM -f "server_invoice"
	@echo "Stopped invoice microservice"

## stop_back: stops the back end
stop_built_back:
	@echo "Stopping the back end..."
	@-pkill -SIGTERM -f "server_main"
	@echo "Stopped back end"

start_back:
	@go run server_main/main.go

start_invoice:
	@go run server_invoice/main.go

start_front:
	@go run frontend/main.go

proto:
	rm -f internal/pb/*.go
	protoc --proto_path=internal/proto --go_out=internal/pb --go_opt=paths=source_relative \
	--go-grpc_out=internal/pb --go-grpc_opt=paths=source_relative \
	internal/proto/*.proto

test:
	go test -v -cover -short ./...

mock:
	mockgen -package pb -destination internal/pb/mock_invoice_service.go github.com/LamThanhNguyen/yoyo-store-backend/internal/pb InvoiceServiceClient
	mockgen -package api -destination server_main/api/mock_interfaces_test.go github.com/LamThanhNguyen/yoyo-store-backend/server_main/api customerInserter,orderInserter,transactionInserter

.PHONY: network postgres createdb dropdb migrateup migrateup1 migratedown migratedown1 new_migration build clean build_back build_invoice build_front start start_back start_invoice start_front stop stop_front stop_invoice stop_back proto mock test
