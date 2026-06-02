# Fase 1 - MVP (Mínimo Produto Viável)

**Duração Estimada:** 2 semanas

**Status:** Em Planejamento

---

## Objetivo

Implementar fluxo básico de entrada de estoque com:
- Interface webcam (React)
- Extração de texto (PaddleOCR)
- Part Number lookup (PostgreSQL)
- Registro em estoque (API Go)

## Fluxo

```
1. Usuário acessa interface Web
2. Captura foto do componente (webcam)
3. Sistema faz OCR da imagem
4. Extrai PN, SN e especificações
5. Busca PN no catálogo
6. Usuário aprova ou corrige
7. Sistema registra entrada em estoque
8. Exibe confirmação
```

---

## Tarefas

### 1.1 - Frontend React
- [ ] Setup: `npx create-react-app` com TypeScript
- [ ] Componente de captura de webcam
- [ ] Componente de exibição de imagem
- [ ] Componente de formulário de aprovação
- [ ] Integração com backend
- [ ] UI básica (Antd ou shadcn)

**Estimativa:** 3 dias

### 1.2 - Database PostgreSQL
- [ ] Criar tabela `catalog` (PN, fabricante, categoria, descrição)
- [ ] Criar tabela `inventory` (item_id, quantidade)
- [ ] Criar tabela `movements` (id, item_id, quantity, operation, timestamp)
- [ ] Popular catálogo básico com ~500 PNs de teste

**Estimativa:** 2 dias

### 1.3 - OCR Service (Python)
- [ ] Setup FastAPI + PaddleOCR
- [ ] Endpoint POST /ocr que recebe imagem
- [ ] Extração de texto com PaddleOCR
- [ ] Normalização de output JSON
- [ ] Docker para o serviço

**Estimativa:** 2 dias

### 1.4 - Inventory API (Go)
- [ ] Setup: `go mod init`
- [ ] Estrutura de projetos (handlers, models, db)
- [ ] Endpoints:
  - `GET /catalog/search?pn=XXXXX` - Buscar PN
  - `POST /inventory/in` - Entrada de estoque
  - `GET /inventory/items` - Listar estoque
- [ ] Conexão com PostgreSQL (database/sql ou gorm)
- [ ] Docker para o serviço

**Estimativa:** 4 dias

### 1.5 - Docker Compose
- [ ] `docker-compose.yml` orquestrando:
  - Frontend (React, porta 3000)
  - OCR Service (Python, porta 5001)
  - Inventory API (Go, porta 8080)
  - PostgreSQL (porta 5432)
  - Redis (porta 6379, opcional para MVP)

**Estimativa:** 1 dia

### 1.6 - Testes & Documentação
- [ ] Testes unitários básicos
- [ ] Testes de integração
- [ ] README.md com instruções de setup
- [ ] Exemplos de requisições (Postman/curl)

**Estimativa:** 2 dias

---

## Critérios de Aceitação

- ✅ Usuário consegue capturar foto via webcam
- ✅ Sistema consegue extrair texto com PaddleOCR
- ✅ Sistema consegue buscar PN no banco
- ✅ Usuário consegue aprovar e registrar entrada
- ✅ Estoque é atualizado corretamente
- ✅ Tudo roda em Docker Compose com um comando

---

## Estrutura de Pastas Esperada Após Fase 1

```
ComputationalVision/
├── frontend/
│   ├── src/
│   │   ├── components/
│   │   │   ├── WebcamCapture.tsx
│   │   │   ├── ImagePreview.tsx
│   │   │   ├── ApprovalForm.tsx
│   │   │   └── HistoryTable.tsx
│   │   ├── services/
│   │   │   └── api.ts
│   │   ├── App.tsx
│   │   └── index.tsx
│   ├── package.json
│   ├── Dockerfile
│   └── tsconfig.json
├── services/
│   ├── ocr/
│   │   ├── main.py
│   │   ├── ocr_service.py
│   │   ├── requirements.txt
│   │   └── Dockerfile
│   └── inventory/
│       ├── main.go
│       ├── handlers.go
│       ├── models.go
│       ├── go.mod
│       └── Dockerfile
├── infra/
│   ├── docker/
│   │   ├── docker-compose.yml
│   │   └── postgres/
│   │       └── init.sql
│   └── kubernetes/
├── docs/
│   ├── 00-ROADMAP.md
│   ├── 01-ARQUITETURA.md
│   ├── 02-TECNOLOGIAS.md
│   └── FASE-1-MVP.md
└── README.md
```

---

## Próximos Passos

Após Fase 1 estar completa:
1. Coleta de dados reais em produção
2. Popular catálogo com mais PNs
3. Fase 2: Implementar Parser (classificação automática)
