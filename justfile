set dotenv-load := true

env:
  set -a && source ./.env

test:
    go test -v ./...

run:
    go run .

build:
    go build -o bin/chirpy .

lint:
    go vet ./...

tidy:
    go mod tidy

psql:
    psql

migrate-up:
    goose up

migrate-down:
    goose down
