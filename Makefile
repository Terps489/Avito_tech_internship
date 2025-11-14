.PHONY: run build docker-build up down logs

run:
	go run ./cmd/app

build:
	go build -o bin/reviewer ./cmd/app

up:
	docker compose up --build

down:
	docker compose down

logs:
	docker compose logs -f app
