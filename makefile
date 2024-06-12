build:
	@go build -o bin/aiplus-test

run: build
	@./bin/aiplus-test

test:
	@go test -v ./...