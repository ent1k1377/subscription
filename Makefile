.PHONY: up

DSN = "postgres://user:pass@localhost:5430/db?sslmode=disable"

up:
	docker compose -f deployments/compose.yaml up

down:
	docker compose -f deployments/compose.yaml down

mig-up:
	goose -dir migrations/ postgres ${DSN} up

mig-down:
	goose -dir migrations/ postgres ${DSN} down