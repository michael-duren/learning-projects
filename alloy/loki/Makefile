.DEFAULT_GOAL := help

COMPOSE := docker compose

.PHONY: help up down logs logs-app logs-alloy raw-logs validate apply app status reset

help:
	@echo "Run from alloy/loki/:"
	@echo "  make up          Build and start all services in the background"
	@echo "  make down        Stop services, preserving data"
	@echo "  make logs        Follow all service logs"
	@echo "  make logs-app    Follow raw application output"
	@echo "  make logs-alloy  Follow processed entries and Alloy diagnostics"
	@echo "  make raw-logs    Follow the file Alloy reads"
	@echo "  make validate    Validate config.alloy"
	@echo "  make apply       Validate, then restart Alloy to apply changes"
	@echo "  make app         Rebuild and restart the Go app"
	@echo "  make status      Show service status"
	@echo "  make reset       Delete all playground data and start fresh"
	@echo "Grafana: http://localhost:3001/explore"

up:
	$(COMPOSE) up --build -d

down:
	$(COMPOSE) down

logs:
	$(COMPOSE) logs -f

logs-app:
	$(COMPOSE) logs -f app

logs-alloy:
	$(COMPOSE) logs -f alloy

raw-logs:
	$(COMPOSE) exec alloy tail -f /logs/app.log

validate:
	$(COMPOSE) run --rm --no-deps alloy validate /etc/alloy/config.alloy

apply: validate
	$(COMPOSE) restart alloy

app:
	$(COMPOSE) up -d --build app

status:
	$(COMPOSE) ps

reset:
	$(COMPOSE) down -v
	$(COMPOSE) up --build -d
