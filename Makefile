NO_CACHE =

.PHONY:build hack setup up down run wire test

setup:
	go install github.com/air-verse/air@v1.52.2	go install github.com/google/wire/cmd/wire@latest

build:
	docker build $(NO_CACHE) -f Dockerfile .

up:
	docker-compose -f docker-compose.local.yaml up --build -d
	sh scripts/up.sh

build-vision:
	sh scripts/build.sh

down:
	docker-compose -f docker-compose.local.yaml down --volumes

run:
	air server --port 9090

wire:
	$(GOPATH)/bin/wire ./di/wire.go

test:
	go test -v ./...

visonary:
	cd tools && go run visonary.go