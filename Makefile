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
	@go build -o dist/server_main ./cmd/server_main
	@echo "Back end built!"

## build_invoice: builds the invoice microservice
build_invoice:
	@echo "Building invoice microservice..."
	@go build -o dist/invoice ./cmd/micro/invoice
	@echo "Invoice microservice built!"

## build_front: builds the front end
build_front:
	@echo "Building front end..."
	@go build -o dist/gostripe ./cmd/web
	@echo "Front end built!"

## start: starts front and back end
start: start_back start_invoice start_front

## start_back: starts the back end
start_back: build_back
	@echo "Starting the back end..."
	@env STRIPE_KEY=${STRIPE_KEY} STRIPE_SECRET=${STRIPE_SECRET} ./dist/gostripe_api -port=${API_PORT} &
	@echo "Back end running!"

## start_invoice: starts the invoice microservice
start_invoice: build_invoice
	@echo "Starting the invoice microservice..."
	@./dist/invoice &
	@echo "Invoice microservice running!"

## start_front: starts the front end
start_front: build_front
	@echo "Starting the front end..."
	@env STRIPE_KEY=${STRIPE_KEY} STRIPE_SECRET=${STRIPE_SECRET} ./dist/gostripe -port=${GOSTRIPE_PORT} &
	@echo "Front end running!"

## stop: stops the front and back end
stop: stop_front stop_back stop_invoice
	@echo "All applications stopped"

## stop_front: stops the front end
stop_front:
	@echo "Stopping the front end..."
	@-pkill -SIGTERM -f "gostripe -port=${GOSTRIPE_PORT}"
	@echo "Stopped front end"

## stop_invoice: stops the invoice microservice
stop_invoice:
	@echo "Stopping the invoice microservice..."
	@-pkill -SIGTERM -f "invoice"
	@echo "Stopped invoice microservice"

## stop_back: stops the back end
stop_back:
	@echo "Stopping the back end..."
	@-pkill -SIGTERM -f "gostripe_api -port=${API_PORT}"
	@echo "Stopped back end"

.PHONY: network postgres createdb dropdb migrateup migrateup1 migratedown migratedown1 new_migration build clean build_back build_invoice build_front start start_back start_invoice start_front stop stop_front stop_invoice stop_back