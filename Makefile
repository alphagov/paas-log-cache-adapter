files = main.go handlers.go middleware.go server.go

build:
	go build $(files) -o bin/metrics

run:
	go run $(files)

test:
	go test -v -cover ./...
