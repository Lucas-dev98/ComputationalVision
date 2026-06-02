# Sistema Inteligente de Entrada de Estoque por Visão Computacional

**Status:** Em desenvolvimento - Fase 1 (MVP)

## Objetivo

Automatizar a entrada de estoque de componentes de TI através de visão computacional, OCR e catálogo inteligente.

## Arquitetura

```
Camera → OCR Service → Parser Service → Catalog Engine → Approval UI → Inventory API
```

## Microserviços

| Serviço | Tecnologia | Status |
|---------|-----------|--------|
| Inventory Service | Go + PostgreSQL | Planejado |
| Vision Service | Python + YOLO | Planejado |
| OCR Service | Python + PaddleOCR | MVP |
| Parser Service | Go | MVP Fase 2 |
| Catalog Service | Go + PostgreSQL | Planejado |
| Web Research Service | Go | Planejado |
| Frontend | React + TypeScript | MVP |

## Fases

- **Fase 1 (MVP):** React + Webcam + PaddleOCR + PostgreSQL (2 semanas)
- **Fase 2:** Parser + Classificação automática (implementada)
- **Fase 3:** Pesquisa web automática (2 semanas)
- **Fase 4:** YOLO - Detecção visual (4 semanas)
- **Fase 5:** Produção - Docker, Logs, Auditoria (contínuo)

## Como começar

Para rodar localmente sem Docker, siga o guia em [LOCAL_SETUP.md](LOCAL_SETUP.md).

Resumo rápido:

```bash
cd frontend
npm install
npm start

cd services/ocr
python main_mock.py

cd services/parser
$env:PORT='8082'
go run .

cd services/inventory
$env:DATABASE_URL='postgres://postgres@localhost:5434/inventory_db?sslmode=disable'
$env:PORT='8081'
go run .
```

Abra a aplicação em http://localhost:3000.

## Documentação

Ver pasta `/docs` para histórico de implementação.
