FROM golang:1.25

WORKDIR /app

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

CMD sh -c "\
goose -dir migrations/ postgres \"postgres://${DATABASE_USERNAME}:${DATABASE_PASSWORD}@postgres:5432/subscription_db?sslmode=disable\" up && \
exec go run ./... \
"