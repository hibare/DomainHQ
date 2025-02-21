SHELL=/bin/bash

UI := $(shell id -u)
GID := $(shell id -g)
MAKEFLAGS += -s
DOCKER_COMPOSE_PREFIX = HOST_UID=${UID} HOST_GID=${GID} docker-compose -f docker-compose.dev.yml

# Bold
BCYAN=\033[1;36m
BBLUE=\033[1;34m

# No color (Reset)
NC=\033[0m

.DEFAULT_GOAL := help

.PHONY: db-up
db-up: ## Start DB Services
	${DOCKER_COMPOSE_PREFIX} up -d postgres adminer

.PHONY: db-down
db-down: ## Stop DB Services
	${DOCKER_COMPOSE_PREFIX} rm -fsv postgres adminer

.PHONY: dev
dev: ## Start App Services
	go mod tidy
	${DOCKER_COMPOSE_PREFIX} up 

.PHONY: clean	
clean: ## Clean up
	${DOCKER_COMPOSE_PREFIX} down
	go mod tidy

.PHONY: test
test: ## Run Tests
ifndef GITHUB_ACTIONS
	$(MAKE) db-up
endif
	go test ./... -cover

.PHONY: help
help: ## Disply this help
		echo -e "\n$(BBLUE)DomainHQ: Domain Admin Services$(NC)\n"
		@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "$(BCYAN)%-18s$(NC)%s\n", $$1, $$2}'