up: up
up-build: up-build
down: down
test: test
migrate: migrate
migrate-down: migrate-down
migration: migration

up:
	@ docker compose up 

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