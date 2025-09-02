.PHONY: up run

DSN = "postgres://user:pass@localhost:5430/db?sslmode=disable"

up:
	docker compose -f deployments/compose.yaml --env-file configs/.env up

down:
	docker compose -f deployments/compose.yaml down

mig-up:
	goose -dir migrations/ postgres ${DSN} up

mig-down:
	goose -dir migrations/ postgres ${DSN} down

run:
	swag fmt
	swag init -g cmd/subscriptions/main.go
	go run ./...