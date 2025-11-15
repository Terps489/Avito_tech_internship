.PHONY: build docker-build up down logs seed load-test post-check

UNAME_S := $(shell uname -s 2> NUL)

ifeq ($(UNAME_S),Linux)
	IS_LINUX := 1
else
	IS_LINUX := 0
endif

post-check:
ifeq ($(IS_LINUX),1)
	bash ./scripts/check_after_load.sh
else
	powershell -ExecutionPolicy Bypass -File scripts/check_after_load_win.ps1
endif

build:
	go build -o bin/reviewer ./cmd/app

up:
	docker compose up --build

down:
	docker compose down

logs:
	docker compose logs -f app

seed:
	go run ./cmd/seed

load-test:
	go run ./cmd/loadtest -duration=30s -rps=5

lint:
	golangci-lint run ./...

check: post-check


