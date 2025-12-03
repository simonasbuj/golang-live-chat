# Load .env.example and .env
ifneq (,$(wildcard .env.example))
    include .env.example
    export $(shell sed 's/=.*//' .env.example)
endif

ifneq (,$(wildcard .env))
    include .env
    export $(shell sed 's/=.*//' .env)
endif

run:
	go run cmd/main.go

run-docker:
	docker compose up --build -d

# pre-commit
lint:
	golangci-lint run --verbose --max-issues-per-linter=0 --max-same-issues=0

lint-fix:
	golangci-lint run --verbose --fix

# test
.PHONY: test
test:
	go test -v ./...