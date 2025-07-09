.PHONY: help dev test build docker-up docker-down

help: ## Mostrar ajuda
	@echo "Comandos disponíveis:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

dev: ## Rodar em desenvolvimento
	go run cmd/api/main.go

test: ## Rodar testes
	go test ./...

build: ## Build da aplicação
	go build -o bin/user-api cmd/api/main.go

docker-up: ## Subir containers
	docker-compose up -d

docker-down: ## Parar containers
	docker-compose down

docker-logs: ## Ver logs
	docker-compose logs -f

clean: ## Limpar containers e volumes
	docker-compose down -v
	docker system prune -f