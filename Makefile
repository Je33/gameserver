.PHONY: dc mocks test build run lint

dc:
	docker-compose up  --remove-orphans --build

mocks:
	go generate ./...

test:
	go test -v -coverprofile cover.out ./... && go tool cover -html=cover.out

build:
	go build -o build/server cmd/server/main.go

run:
	go run -race cmd/server/main.go

lint:
	golangci-lint run