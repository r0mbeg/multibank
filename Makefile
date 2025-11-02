# ===========
# Settings
# ===========
APP_NAME        ?= multibank-backend
IMAGE           ?= $(APP_NAME)
TAG             ?= dev
CONTAINER_NAME  ?= $(APP_NAME)-$(TAG)

BACKEND_DIR     := backend
DOCKERFILE      := $(BACKEND_DIR)/Dockerfile
BUILD_CONTEXT   := $(BACKEND_DIR)

# Абсолютные пути для volume-монтов (важно для Windows)
CFG_FILE        := $(abspath $(BACKEND_DIR)/config/local.yaml)
STORAGE_DIR     := $(abspath $(BACKEND_DIR)/storage)
LOGS_DIR        := $(abspath $(BACKEND_DIR)/logs)

# Параметры приложения
PORT            ?= 8080
LOG_LEVEL       ?= debug

# Docker Buildx платформа (можно закомментировать, если не нужна)
PLATFORM        ?= linux/amd64

# ===========
# Phony
# ===========
.PHONY: build run stop restart logs sh rm image-rm clean test swagger swagger-install help

help:
	@echo "Targets:"
	@echo "  build           - docker build образ $(IMAGE):$(TAG)"
	@echo "  run             - запустить контейнер (порт $(PORT), volumes config/storage/logs)"
	@echo "  stop            - остановить контейнер"
	@echo "  restart         - перезапустить контейнер"
	@echo "  logs            - показать логи контейнера (follow)"
	@echo "  sh              - зайти в shell внутри контейнера"
	@echo "  rm              - удалить контейнер (если есть)"
	@echo "  image-rm        - удалить образ $(IMAGE):$(TAG)"
	@echo "  clean           - stop + rm + image-rm"
	@echo "  test            - прогнать backend-тесты (пакет ./backend/tests)"
	@echo "  swagger-install - установить swag CLI"
	@echo "  swagger         - сгенерировать Swagger (./backend/docs)"
	@echo ""
	@echo "Vars (override via make VAR=...):"
	@echo "  TAG=$(TAG) PORT=$(PORT) LOG_LEVEL=$(LOG_LEVEL)"

# ========== Build ==========
build:
	@echo "==> Building docker image $(IMAGE):$(TAG)"
	@docker build --platform $(PLATFORM) -f $(DOCKERFILE) -t $(IMAGE):$(TAG) $(BUILD_CONTEXT)

# ========== Run / Stop ==========
run: # предварительно создаём директории, чтобы монтирование не падало
	@mkdir -p "$(STORAGE_DIR)" "$(LOGS_DIR)"
	@echo "==> Running $(CONTAINER_NAME) from $(IMAGE):$(TAG)"
	@docker run --rm -d \
		--name $(CONTAINER_NAME) \
		-p $(PORT):8080 \
		-e MB_LOG_LEVEL=$(LOG_LEVEL) \
		-v "$(CFG_FILE):/etc/multibank/config.yaml:ro" \
		-v "$(STORAGE_DIR):/app/storage" \
		-v "$(LOGS_DIR):/app/logs" \
		$(IMAGE):$(TAG)
	@echo "-> http://localhost:$(PORT)   (Swagger: /swagger/index.html)"

stop:
	@docker stop $(CONTAINER_NAME) 2>/dev/null || true

restart: stop run

logs:
	@docker logs -f $(CONTAINER_NAME)

sh:
	@docker exec -it $(CONTAINER_NAME) /bin/sh || docker exec -it $(CONTAINER_NAME) sh

rm:
	@docker rm -f $(CONTAINER_NAME) 2>/dev/null || true

image-rm:
	@docker rmi $(IMAGE):$(TAG) 2>/dev/null || true

clean: stop rm image-rm

# ========== Dev helpers ==========
test:
	@echo "==> go test ./backend/tests"
	@cd backend && go test -v ./tests

swagger-install:
	@echo "==> Installing swag CLI"
	@go install github.com/swaggo/swag/cmd/swag@latest

swagger:
	@echo "==> Generating Swagger docs"
	@cd backend && swag init -g cmd/backend/main.go -o docs
