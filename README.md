# Sistema Inteligente de Entrada de Estoque por Visão Computacional

**Status:** Em desenvolvimento - Fase 3 (Pesquisa web automática implementada)

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
| Web Research Service | Go | MVP Fase 3 |
| Frontend | React + TypeScript | MVP |

## Fases

- **Fase 1 (MVP):** React + Webcam + PaddleOCR + PostgreSQL (2 semanas)
- **Fase 2:** Parser + Classificação automática (implementada)
- **Fase 3:** Pesquisa web automática (implementada)
- **Fase 4:** YOLO - Detecção visual (4 semanas)
- **Fase 5:** Produção - Docker, Logs, Auditoria (contínuo)

## Como começar

Para rodar localmente sem Docker, siga o guia em [LOCAL_SETUP.md](LOCAL_SETUP.md).

Fluxo sem Docker (minimo para Fase 3): Frontend + OCR + Parser + Web Research.
O Inventory Service fica opcional enquanto o PostgreSQL local nao estiver disponivel.

Resumo rápido:

```bash
cd frontend
npm install
npm start

cd services/ocr
python main.py

cd services/parser
$env:PORT='8082'
go run .

cd services/web-research
$env:PORT='8083'
go run .

cd services/inventory
$env:DATABASE_DRIVER='sqlite'
$env:DATABASE_URL='file:./inventory-dev.db?_pragma=foreign_keys(1)'
$env:PORT='8081'
go run .
```

Abra a aplicação em http://localhost:3000.

## Documentação

Ver pasta `/docs` para histórico de implementação.

## Fase 3 - Pesquisa Web Automática

Quando o catálogo não encontra o Part Number, o frontend aciona automaticamente o serviço `web-research` para tentar enriquecer:

- fabricante sugerido
- categoria sugerida
- descrição normalizada sugerida
- fontes de pesquisa e confiança

Endpoint do serviço:

```bash
POST http://localhost:8083/research
Content-Type: application/json

{
	"part_number": "M393A4K40DB3-CWE",
	"manufacturer": "",
	"category": "unknown",
	"normalized_description": "",
	"tokens": ["M393A4K40DB3-CWE", "DDR4", "32GB"]
}
```

## Modo Local Com Banco Embutido (sem PostgreSQL)

O Inventory Service agora suporta SQLite para desenvolvimento local.

```bash
cd services/inventory
$env:DATABASE_DRIVER='sqlite'
$env:DATABASE_URL='file:./inventory-dev.db?_pragma=foreign_keys(1)'
$env:PORT='8081'
go run .
```

No primeiro start ele cria schema e insere seed inicial automaticamente.
