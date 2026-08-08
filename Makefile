.PHONY: up down build logs ps test backend-test backend-build frontend-build clean

## Start the full stack (build + daemon)
up:
	docker compose up -d --build

## Stop and remove containers
down:
	docker compose down

## Rebuild images without recreating running containers
build:
	docker compose build

## Tail logs of all services
logs:
	docker compose logs -f

## List running services and their status
ps:
	docker compose ps

## Run all backend tests
test: backend-test

## Run Go unit tests locally
backend-test:
	cd backend && go test ./...

## Vet and build the backend locally
backend-build:
	cd backend && go vet ./... && go build ./...

## Build the frontend bundle locally
frontend-build:
	cd frontend && npm run build

## Stop containers and remove volumes
clean:
	docker compose down -v
