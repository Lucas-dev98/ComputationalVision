.PHONY: help install build up down logs clean test lint format

help:
	@echo "Sistema de Entrada de Estoque - Comandos Disponíveis"
	@echo ""
	@echo "Infraestrutura:"
	@echo "  make up           - Inicia todos os serviços (Docker)"
	@echo "  make down         - Para todos os serviços"
	@echo "  make logs         - Visualiza logs"
	@echo "  make clean        - Remove containers, volumes e cache"
	@echo ""
	@echo "Desenvolvimento Local:"
	@echo "  make install      - Instala dependências de todos os serviços"
	@echo "  make frontend-dev - Inicia frontend em desenvolvimento"
	@echo "  make ocr-dev      - Inicia OCR Service em desenvolvimento"
	@echo "  make api-dev      - Inicia Inventory API em desenvolvimento"
	@echo ""
	@echo "Qualidade:"
	@echo "  make test         - Executa todos os testes"
	@echo "  make lint         - Executa linter"
	@echo "  make format       - Formata código"
	@echo ""
	@echo "Banco de Dados:"
	@echo "  make db-init      - Inicializa banco (via docker)"
	@echo "  make db-reset     - Reseta banco completamente"
	@echo ""

# Infraestrutura
up:
	docker-compose -f infra/docker/docker-compose.yml up -d
	@echo "✓ Serviços iniciados"
	@echo "  Frontend:      http://localhost:3000"
	@echo "  Inventory API: http://localhost:8080"
	@echo "  OCR Service:   http://localhost:5001"
	@echo "  PostgreSQL:    localhost:5432"
	@echo "  Redis:         localhost:6379"

down:
	docker-compose -f infra/docker/docker-compose.yml down
	@echo "✓ Serviços parados"

logs:
	docker-compose -f infra/docker/docker-compose.yml logs -f

logs-frontend:
	docker-compose -f infra/docker/docker-compose.yml logs -f frontend

logs-api:
	docker-compose -f infra/docker/docker-compose.yml logs -f inventory-service

logs-ocr:
	docker-compose -f infra/docker/docker-compose.yml logs -f ocr-service

logs-db:
	docker-compose -f infra/docker/docker-compose.yml logs -f postgres

clean:
	docker-compose -f infra/docker/docker-compose.yml down -v
	find . -type d -name __pycache__ -exec rm -rf {} +
	find . -type f -name "*.pyc" -delete
	find . -type d -name node_modules -exec rm -rf {} +
	@echo "✓ Limpeza completa realizada"

# Banco de Dados
db-init:
	docker-compose -f infra/docker/docker-compose.yml up -d postgres
	@sleep 3
	docker-compose -f infra/docker/docker-compose.yml exec postgres psql -U inventory -d inventory_db -f /docker-entrypoint-initdb.d/init.sql
	@echo "✓ Banco inicializado"

db-reset: clean up
	@echo "✓ Sistema resetado completamente"

# Desenvolvimento Local
install: install-frontend install-ocr install-api

install-frontend:
	cd frontend && npm install
	@echo "✓ Frontend dependências instaladas"

install-ocr:
	cd services/ocr && pip install -r requirements.txt
	@echo "✓ OCR dependências instaladas"

install-api:
	cd services/inventory && go mod download
	@echo "✓ API dependências instaladas"

frontend-dev:
	cd frontend && npm start

ocr-dev:
	cd services/ocr && python main.py

api-dev:
	cd services/inventory && go run .

# Qualidade de Código
test:
	@echo "Executando testes..."
	@echo "TODO: Implementar testes"

lint:
	@echo "Linting..."
	cd services/ocr && flake8 . --max-line-length=100
	cd frontend && npm run lint 2>/dev/null || echo "ESLint não configurado"
	@echo "✓ Linting concluído"

format:
	@echo "Formatando código..."
	cd services/ocr && black .
	cd services/inventory && gofmt -w .
	cd frontend && npx prettier --write "src/**/*.{ts,tsx,css}" 2>/dev/null || echo "Prettier não configurado"
	@echo "✓ Formatação concluída"

# Git
status:
	git status
	@echo ""
	git log --oneline -10

commits-history:
	@cat docs/COMMITS.md

# Defaults
.DEFAULT_GOAL := help
