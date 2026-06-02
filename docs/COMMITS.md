# Histórico de Commits

## Commit 1: Estrutura Base + Frontend + OCR Service + Inventory Service

**Hash:** `[será gerado pelo git]`

**Mensagem:** 
```
feat: estrutura base, frontend React, OCR Service e Inventory Service - Fase 1 MVP
```

**Data:** 01/06/2026

**O que foi implementado:**

### Estrutura do Projeto
- ✅ Organização em microserviços
- ✅ Pastas para Frontend, Services (OCR, Inventory, Vision, Parser, Catalog, Web-Research)
- ✅ Infra com Docker e Kubernetes
- ✅ Documentação em Obsidian

### Frontend (React + TypeScript)
- ✅ Componente WebcamCapture - Captura de foto via webcam
- ✅ Componente ImagePreview - Exibição de preview
- ✅ Componente ApprovalForm - Formulário para aprovação e entrada de estoque
- ✅ Componente HistoryTable - Tabela com histórico de estoque
- ✅ Serviço de API - Integração com backend
- ✅ Layout com Ant Design
- ✅ TypeScript configuration

### OCR Service (Python + FastAPI + PaddleOCR)
- ✅ Endpoint POST /ocr - Extração de texto
- ✅ Endpoint POST /ocr/batch - Processamento em lote
- ✅ Health check
- ✅ Processamento com PaddleOCR
- ✅ Suporte a português e inglês
- ✅ Retorno estruturado (texto, confiança, boxes)

### Inventory Service (Go + PostgreSQL)
- ✅ Endpoints:
  - GET /health - Status do serviço
  - GET /catalog/search?pn=XXXXX - Buscar Part Number
  - POST /inventory/in - Entrada de estoque
  - GET /inventory/items - Listar estoque
  - GET /inventory/items/{id} - Obter item
- ✅ Conexão com PostgreSQL
- ✅ Models: CatalogItem, InventoryItem, Movement
- ✅ CORS habilitado

### Banco de Dados (PostgreSQL)
- ✅ Tabela `catalog` - Part Numbers e especificações
- ✅ Tabela `inventory` - Itens em estoque
- ✅ Tabela `movements` - Histórico de movimentações
- ✅ Tabela `audit_log` - Logs de auditoria
- ✅ Índices para performance
- ✅ Dados de teste com 10 Part Numbers

### Docker & Compose
- ✅ docker-compose.yml orquestrando:
  - PostgreSQL
  - Redis
  - OCR Service
  - Inventory Service
  - Frontend React
- ✅ Dockerfile para cada serviço
- ✅ Network interno compartilhado
- ✅ Health checks

### Documentação
- ✅ README.md principal
- ✅ Roadmap com fases
- ✅ Arquitetura do sistema
- ✅ Justificativa de tecnologias
- ✅ Detalhes de Fase 1 (MVP)
- ✅ Arquivo de commits (este)

---

## Próximos Commits (Fase 1 continuação)

### Commit 2: Testes Unitários + CI/CD
- [ ] Testes do OCR Service
- [ ] Testes do Inventory Service
- [ ] GitHub Actions workflow
- [ ] Lint e formatação de código

### Commit 3: Melhorias de UX
- [ ] Validação de forma aprimorada
- [ ] Feedback visual melhorado
- [ ] Dark mode (opcional)
- [ ] Responsividade aprimorada

### Commit 4: Integração Completa
- [ ] Testes E2E
- [ ] Exemplo de uso
- [ ] Documentação de API (Swagger)
- [ ] Setup guide

---

## Fase 2: Parser + Classificação

### Commit 5: Parser Service (Go)
- [ ] Classificador de Memórias (DDR3/4/5, capacidade, etc)
- [ ] Classificador de Discos (SATA/SAS/NVMe)
- [ ] Classificador de Rede (RJ45/SFP)
- [ ] Endpoint de parsing

### Commit 6: Catalog Inteligente
- [ ] Expansão do banco com mais Part Numbers
- [ ] Normalização de descrições
- [ ] Integração com Parser

---

## Fase 3: Pesquisa Web

### Commit 7: Web Research Service
- [ ] Scraping de PNs
- [ ] Cache com Redis
- [ ] Integração com Catalog

---

## Fase 4: YOLO

### Commit 8: Vision Service
- [ ] Model YOLO treinado
- [ ] Detecção de componentes
- [ ] Localização de labels

---

## Fase 5: Produção

### Commit 9: Logging & Monitoramento
- [ ] Prometheus metrics
- [ ] Grafana dashboards
- [ ] Structured logging

### Commit 10: Deploy
- [ ] Kubernetes manifests
- [ ] Helm charts
- [ ] CI/CD pipeline completo
