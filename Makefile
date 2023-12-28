build:
	@go build -o bin/musicSearchAPI

run: build
	@./bin/musicSearchAPI

test: 
	@go test -v ./...