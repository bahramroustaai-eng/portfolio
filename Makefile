include .env
export

.PHONY: up down migrate-status migrate-up migrate-down migrate-new sqlc test

up:
	docker compose up -d

down:
	docker compose down

migrate-new:
	goose -dir migrations create $(name) sql

migrate-up:
	goose -dir migrations postgres "$(DATABASE_URL)" up

migrate-down:
	goose -dir migrations postgres "$(DATABASE_URL)" down

migrate-status:
	goose -dir migrations postgres "$(DATABASE_URL)" status

sqlc:
	sqlc generate

test:
	go test ./... -race