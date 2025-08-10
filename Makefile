.PHONY: build test lint docker-build

build:
go build ./...

test:
go test ./...

lint:
gofmt -w $(shell find . -name '*.go')

docker-build:
docker build -t slipway-controller -f core/controller/Dockerfile .
docker build -t slipway-cli -f cli/Dockerfile .
