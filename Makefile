NO_CACHE =

.PHONY:build

setup:
	go install github.com/air-verse/air@v1.52.2
	go install github.com/google/wire/cmd/wire@latest

build:
	docker build $(NO_CACHE) -f Dockerfile .

up:
	docker-compose -f docker-compose.yaml up --build

down:
	docker-compose -f docker-compose.yaml down --volumes

run:
	air server --port 9090

wire:
	$(GOPATH)/bin/wire ./di/wire.go
