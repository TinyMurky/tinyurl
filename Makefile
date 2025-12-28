ENV ?= dev
ENV_PATH := ./env/.env.$(ENV)
COMPOSE_PATH := ./builders/docker-compose.yml

up:
	@echo "Starting Docker with environment: $(ENV)..."
	docker compose -f $(COMPOSE_PATH) --env-file $(ENV_PATH) up -d

down:
	docker compose -f $(COMPOSE_PATH) --env-file $(ENV_PATH) down

restart:
	docker compose -f $(COMPOSE_PATH) --env-file $(ENV_PATH) restart

ps:
	docker compose -f $(COMPOSE_PATH) --env-file $(ENV_PATH) ps

log:
	docker compose -f $(COMPOSE_PATH) --env-file $(ENV_PATH) logs -f
