APP:=readit

all: build test lint run

build:
	go build -o $(APP)

test:
	go test ./...

lint:
	golangci-lint run

run:
	./$(APP)