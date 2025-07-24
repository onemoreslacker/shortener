.PHONY: up
up:
	@docker compose -f deploy/compose.yaml up --build

.PHONY: up-silent
up-silent:
	@docker compose -f deploy/compose.yaml up --build -d

.PHONY: down
down:
	@docker compose -f deploy/compose.yaml down -v