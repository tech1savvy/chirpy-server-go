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