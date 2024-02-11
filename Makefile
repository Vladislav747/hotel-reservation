build:
	@go build -o cmd

run:
	@go run ./cmd

seed:
	@go run scripts/seed.go

test:
	@go test -v ./...