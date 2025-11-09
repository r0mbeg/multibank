# ========= Vars =========
COMPOSE ?= docker compose
DC_FILE ?= docker-compose.yml

BACKEND_SERVICE   ?= backend
FRONTEND_SERVICE  ?= frontend

# Переопределяемые переменные окружения
MB_LOG_LEVEL ?= debug
BACKEND_PORT ?= 8080
FRONTEND_PORT ?= 5173

# ========= Phony =========
.PHONY: build-backend build-frontend build-all \
        run-backend run-frontend run-all \
        stop-backend stop-frontend stop-all \
        logs-backend logs-frontend logs-all \
        sh-backend sh-frontend ps

# ========= Build =========
build-backend:
	$(COMPOSE) -f $(DC_FILE) build $(BACKEND_SERVICE)

build-frontend:
	$(COMPOSE) -f $(DC_FILE) build $(FRONTEND_SERVICE)

build-all:
	$(COMPOSE) -f $(DC_FILE) build

# ========= Run =========
run-backend:
	MB_LOG_LEVEL=$(MB_LOG_LEVEL) BACKEND_PORT=$(BACKEND_PORT) $(COMPOSE) -f $(DC_FILE) up -d $(BACKEND_SERVICE)

run-frontend:
	MB_LOG_LEVEL=$(MB_LOG_LEVEL) FRONTEND_PORT=$(FRONTEND_PORT) $(COMPOSE) -f $(DC_FILE) up -d $(FRONTEND_SERVICE)

run-all:
	MB_LOG_LEVEL=$(MB_LOG_LEVEL) BACKEND_PORT=$(BACKEND_PORT) FRONTEND_PORT=$(FRONTEND_PORT) $(COMPOSE) -f $(DC_FILE) up -d

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
