.PHONY: run docs swagger-install

run:
	go run ./cmd/backend --config=./config/local.yaml

swagger-install:
	go install github.com/swaggo/swag/cmd/swag@latest

docs:
	swag init -g cmd/backend/main.go -o docs
