NO_CACHE =

.PHONY:build

build:
	docker build $(NO_CACHE) -f Dockerfile .