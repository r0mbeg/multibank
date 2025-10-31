# -------------------- Build stage --------------------
FROM golang:1.25-alpine AS builder

WORKDIR /src

# Кэшируем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходники
COPY . .

# На всякий случай привести зависимости в порядок (если chi/lumberjack помечены как indirect)
RUN go mod tidy

# Сборка статического бинаря (CGO не нужен для modernc.org/sqlite)
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -ldflags "-s -w" -o /out/app ./cmd/backend

# -------------------- Runtime stage --------------------
FROM alpine:3.20

# TLS/время
RUN apk add --no-cache ca-certificates tzdata

# Непривилегированный пользователь
RUN addgroup -S app && adduser -S app -G app

WORKDIR /app
COPY --from=builder /out/app /app/app

# Рабочие директории
RUN mkdir -p /app/logs /app/storage /etc/multibank && \
    chown -R app:app /app /etc/multibank

USER app

EXPOSE 8080

# Ожидаем, что конфиг примонтируют сюда:
ENTRYPOINT ["/app/app"]
CMD ["--config=/etc/multibank/config.yaml"]

# docker run --rm -p 8080:8080 -v "$PWD/config/local.yaml:/etc/multibank/config.yaml:ro" -v "$PWD/storage:/app/storage" -v "$PWD/logs:/app/logs" -e MB_LOG_LEVEL=debug multibank-backend:dev