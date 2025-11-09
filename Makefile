# ========= Export all Make vars to child processes (cross-platform) =========
.EXPORT_ALL_VARIABLES:

# ========= Vars =========
COMPOSE ?= docker compose
DC_FILE ?= docker-compose.yml

BACKEND_SERVICE  ?= backend
FRONTEND_SERVICE ?= frontend

# Переменные окружения, которыми пользуется docker-compose.yml
MB_LOG_LEVEL  ?= debug
BACKEND_PORT  ?= 8080
FRONTEND_PORT ?= 5173

# ========= Phony =========
.PHONY: help \
        build-backend build-frontend build-all \
        run-backend run-frontend run-all \
        stop-backend stop-frontend stop-all \
        logs-backend logs-frontend logs-all \
        sh-backend sh-frontend ps

# ========= Help (default) =========
help:
	@echo ""
	@echo "=== Multibank Make Targets (no .env, cross-platform) ==="
	@echo "Build:"
	@echo "  build-backend        Build backend image"
	@echo "  build-frontend       Build frontend image"
	@echo "  build-all            Build all images"
	@echo ""
	@echo "Run / Stop:"
	@echo "  run-backend          Run backend container"
	@echo "  run-frontend         Run frontend container"
	@echo "  run-all              Run all containers"
	@echo "  stop-backend         Stop backend container"
	@echo "  stop-frontend        Stop frontend container"
	@echo "  stop-all             Stop and remove all containers"
	@echo ""
	@echo "Logs / Shell / Info:"
	@echo "  logs-backend         Show backend logs (follow)"
	@echo "  logs-frontend        Show frontend logs (follow)"
	@echo "  logs-all             Show all logs"
	@echo "  sh-backend           Open shell in backend container"
	@echo "  sh-frontend          Open shell in frontend container"
	@echo "  ps                   Show containers state"
	@echo ""
	@echo "Current env (override like: make run-all BACKEND_PORT=8081 MB_LOG_LEVEL=info):"
	@echo "  MB_LOG_LEVEL=$(MB_LOG_LEVEL)"
	@echo "  BACKEND_PORT=$(BACKEND_PORT)"
	@echo "  FRONTEND_PORT=$(FRONTEND_PORT)"
	@echo ""

# ========= Build =========
build-backend:
	$(COMPOSE) -f $(DC_FILE) build $(BACKEND_SERVICE)

build-frontend:
	$(COMPOSE) -f $(DC_FILE) build $(FRONTEND_SERVICE)

build-all:
	$(COMPOSE) -f $(DC_FILE) build

# ========= Run =========
run-backend:
	$(COMPOSE) -f $(DC_FILE) up -d $(BACKEND_SERVICE)

run-frontend:
	$(COMPOSE) -f $(DC_FILE) up -d $(FRONTEND_SERVICE)

run-all:
	$(COMPOSE) -f $(DC_FILE) up -d

# ========= Stop =========
stop-backend:
	$(COMPOSE) -f $(DC_FILE) stop $(BACKEND_SERVICE)

stop-frontend:
	$(COMPOSE) -f $(DC_FILE) stop $(FRONTEND_SERVICE)

stop-all:
	$(COMPOSE) -f $(DC_FILE) down

# ========= Logs / Shell / PS =========
logs-backend:
	$(COMPOSE) -f $(DC_FILE) logs -f $(BACKEND_SERVICE)

logs-frontend:
	$(COMPOSE) -f $(DC_FILE) logs -f $(FRONTEND_SERVICE)

logs-all:
	$(COMPOSE) -f $(DC_FILE) logs -f

sh-backend:
	$(COMPOSE) -f $(DC_FILE) exec $(BACKEND_SERVICE) sh

sh-frontend:
	$(COMPOSE) -f $(DC_FILE) exec $(FRONTEND_SERVICE) sh

ps:
	$(COMPOSE) -f $(DC_FILE) ps
