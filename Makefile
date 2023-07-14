SHELL=/bin/bash

UI := $(shell id -u)
GID := $(shell id -g)
MAKEFLAGS += -s
DOCKER_COMPOSE_PREFIX = HOST_UID=${UID} HOST_GID=${GID} docker-compose -f docker-compose.dev.yml

all: app-up

db-up:
	${DOCKER_COMPOSE_PREFIX} up -d postgres adminer

db-down:
	${DOCKER_COMPOSE_PREFIX} rm -fsv postgres adminer

app-up:
	go mod tidy
	${DOCKER_COMPOSE_PREFIX} up 
	
clean: 
	${DOCKER_COMPOSE_PREFIX} down
	go mod tidy

test: 
	$(MAKE) db-up
	go test ./... -cover

.PHONY = all clean app-up test db-up db-down