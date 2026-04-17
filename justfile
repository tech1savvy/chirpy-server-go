test:
    go test ./...

run:
    go run .

build:
    go build -o bin/chirpy .

lint:
    go vet ./...

tidy:
    go mod tidy

bruno:
    cd bruno && bru run . --env Local

db:
  docker compose up -d
