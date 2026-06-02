# Sistema Inteligente de Entrada de Estoque por Visão Computacional

**Status:** Em desenvolvimento - Fase 1 (MVP)

## Objetivo

Automatizar a entrada de estoque de componentes de TI através de visão computacional, OCR e catálogo inteligente.

## Arquitetura

```
Camera → Vision Service → OCR Service → Parser → Catalog Engine → Approval UI → Inventory API
```

## Microserviços

| Serviço | Tecnologia | Status |
|---------|-----------|--------|
| Inventory Service | Go + PostgreSQL | Planejado |
| Vision Service | Python + YOLO | Planejado |
| OCR Service | Python + PaddleOCR | MVP |
| Parser Service | Go | Planejado |
| Catalog Service | Go + PostgreSQL | Planejado |
| Web Research Service | Go | Planejado |
| Frontend | React + TypeScript | MVP |

## Fases

- **Fase 1 (MVP):** React + Webcam + PaddleOCR + PostgreSQL (2 semanas)
- **Fase 2:** Parser + Classificação automática (3 semanas)
- **Fase 3:** Pesquisa web automática (2 semanas)
- **Fase 4:** YOLO - Detecção visual (4 semanas)
- **Fase 5:** Produção - Docker, Logs, Auditoria (contínuo)

## Como começar

```bash
cd /home/lucasbastos/ComputationalVision
docker-compose up
```

## Documentação

Ver pasta `/docs` para histórico de implementação.
