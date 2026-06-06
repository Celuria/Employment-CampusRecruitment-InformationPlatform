APP_NAME=campus-recruitment-server
MAIN=./cmd/server

.PHONY: run build tidy test clean

run:
	go run $(MAIN)/main.go

build:
	go build -o bin/server $(MAIN)/main.go

tidy:
	go mod tidy

test:
	go test ./...

clean:
	rm -rf bin/
