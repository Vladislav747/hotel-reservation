build:
	@go build -o cmd

run:
	@go run ./cmd

test:
	@go test -v ./...