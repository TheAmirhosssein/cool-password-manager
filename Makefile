include .env
export $(shell sed 's/=.*//' .env)
POSTGRES_CONNECTION_STRING = postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable

run: run
up: up
up-build: up-build
down: down
test: test
migrate: migrate
migrate-down: migrate-down
migration: migration

run:
	@ air -c .air.toml

up:
	@ docker compose up  -d

up-build:
	@ docker compose up --build

down:
	@ docker compose down

test:
	@ find . -type d -name 'test*' -exec go test {}/... \;

generate-doc:
	@ swag init -g ./cmd/main.go 

migrate:
	@ goose -dir ./internal/infrastructure/database/migrations postgres $(POSTGRES_CONNECTION_STRING) up

migrate-down:
	@ goose -dir ./internal/infrastructure/database/migrations postgres $(POSTGRES_CONNECTION_STRING) down

migration:
	@read -p "Enter your migration name: " input; \
	goose -dir ./internal/infrastructure/database/migrations create "$$input" sql