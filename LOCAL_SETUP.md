# 🚀 Setup Local - Sistema de Entrada de Estoque

## Visão Geral

Este documento descreve como rodar toda a stack do projeto localmente para desenvolvimento e testes.

### Arquitetura

```
┌─────────────────────────────────────────────────────────────┐
│                    Frontend React                           │
│              (localhost:3000 - Static Server)               │
└──────────────────────────┬──────────────────────────────────┘
                           │
                    ┌──────▼──────┐
                    │  JavaScript │
                    │  Calls API  │
                    └──────┬──────┘
                           │
    ┌──────────────────────┼──────────────────────┐
    │                      │                      │
    ▼                      ▼                      ▼
┌──────────┐  ┌──────────────────┐  ┌──────────────────┐
  │ Go API   │  │ OCR Service      │  │ PostgreSQL       │
  │:8081     │  │ (Mock):5001      │  │ :5434            │
└──────────┘  └──────────────────┘  └──────────────────┘
  ┌──────────┐
  │ Parser   │
  │:8082     │
  └──────────┘
```

## 📋 Pré-requisitos

- **Node.js** 18+ e npm
- **Go** 1.21+
- **Python** 3.9+
- **PostgreSQL** 15 (opcional nesta fase)
- **Make** (para build tasks)

Docker e Docker Compose sao opcionais. O fluxo principal abaixo roda sem containers.

## ⚡ Fluxo sem Docker (recomendado para este computador)

Abra 4 terminais para executar os servicos base:

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
```

Opcionalmente, abra um 5o terminal para o Inventory Service se tiver PostgreSQL local configurado:

```bash
cd services/inventory
$env:DATABASE_URL='postgres://postgres@localhost:5434/inventory_db?sslmode=disable'
$env:PORT='8081'
go run .
```

Mesmo sem Inventory, o frontend agora continua com parser + pesquisa web automatica.

### Opcao recomendada: Inventory com banco embutido (SQLite)

Se quiser testar catalogo e entrada de estoque sem PostgreSQL, rode:

```bash
cd services/inventory
$env:DATABASE_DRIVER='sqlite'
$env:DATABASE_URL='file:./inventory-dev.db?_pragma=foreign_keys(1)'
$env:PORT='8081'
go run .
```

No primeiro start, o servico cria as tabelas e um seed minimo no arquivo local `inventory-dev.db`.

## 🔧 Configuração Inicial

### 1. Clone e Dependências

```bash
cd /home/lucasbastos/ComputationalVision
npm install --prefix frontend
cd services/inventory && go mod tidy && cd ../..
cd services/ocr && pip install -r requirements-simple.txt && cd ../..
```

### 2. PostgreSQL (Docker)

Se não tiver PostgreSQL instalado localmente, use Docker:

```bash
# Criar container
docker rm -f postgres-dev 2>/dev/null
docker run -d \
  --name postgres-dev \
  -e POSTGRES_USER=inventory \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=inventory \
  -p 5434:5432 \
  postgres:15-alpine

# Verificar se está rodando
sleep 3 && docker ps | grep postgres-dev
```

### 3. Inicializar Schema do Banco

```bash
psql -h localhost -p 5434 -U inventory -d inventory <<'EOF'
CREATE TABLE IF NOT EXISTS catalog (
  id SERIAL PRIMARY KEY,
  part_number VARCHAR(255) UNIQUE NOT NULL,
  serial_pattern VARCHAR(255),
  manufacturer VARCHAR(255),
  category VARCHAR(255),
  normalized_description TEXT,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS inventory (
  id SERIAL PRIMARY KEY,
  catalog_id INTEGER NOT NULL REFERENCES catalog(id),
  serial_number VARCHAR(255) UNIQUE NOT NULL,
  quantity INTEGER,
  location VARCHAR(255),
  status VARCHAR(50) DEFAULT 'active',
  received_at TIMESTAMP DEFAULT NOW(),
  last_updated TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS movement (
  id SERIAL PRIMARY KEY,
  inventory_id INTEGER NOT NULL REFERENCES inventory(id),
  quantity INTEGER NOT NULL,
  operation VARCHAR(50),
  reason VARCHAR(255),
  user_id VARCHAR(255),
  created_at TIMESTAMP DEFAULT NOW()
);

\dt  -- Listar tabelas
EOF
```

## 🎯 Executar os Serviços

Abra 4 terminais diferentes:

### Terminal 1 - Go API Service

```bash
cd /home/lucasbastos/ComputationalVision/services/inventory

# Build (opcional, só primeira vez)
go build -o inventory-service

# Rodar
DATABASE_URL="postgres://inventory:password@localhost:5434/inventory?sslmode=disable" \
PORT=8081 \
go run .
```

**Esperado:**
```
Server running on port 8081
```

### Terminal 2 - OCR Service (Mock)

```bash
cd /home/lucasbastos/ComputationalVision/services/ocr

# Rodar
python3 main_mock.py
```

**Esperado:**
```
INFO:     Uvicorn running on http://0.0.0.0:5001 (Press CTRL+C to quit)
```

### Terminal 3 - Parser Service

```bash
cd services/parser
$env:PORT='8082'
go run .
```

**Esperado:**
```
Parser Service escutando em 0.0.0.0:8082
```

### Terminal 4 - Frontend Static Server

```bash
cd /home/lucasbastos/ComputationalVision/frontend

# Build React (primeira vez ou após mudanças)
npm run build

# Servir arquivos estáticos
cd build
python3 -m http.server 3000
```

**Esperado:**
```
Serving HTTP on 0.0.0.0:3000 (http://0.0.0.0:3000/) ...
```

## ✅ Verificações de Health

Em um quinto terminal, verifique se tudo está rodando:

```bash
# Frontend
curl http://localhost:3000/ | head -1

# Go API
curl http://localhost:8081/health

# OCR Service
curl http://localhost:5001/health

# Parser Service
curl http://localhost:8082/health

# Inventory items
curl http://localhost:8081/inventory/items
```

**Saída esperada:**

```
✅ Frontend: <!DOCTYPE html>
✅ Go API: {"service":"inventory-service","status":"healthy","version":"1.0.0"}
✅ OCR: {"status":"healthy","service":"ocr-service","version":"1.0.0-mock"}
✅ Parser: {"status":"healthy","service":"parser-service","version":"1.0.0"}
✅ Items: {"total":0,"items":[]}
```

## 🌐 Acessar a Aplicação

Abra o navegador em: **http://localhost:3000**

### Componentes Principais

1. **Câmera** - Widget para capturar fotos (mock no ambiente virtual)
2. **Preview** - Visualizar imagem capturada
3. **Aprovação** - Confirmação dos dados OCR
4. **Histórico** - Lista de itens no estoque

## 🛠️ Desenvolvimento

### Modificar Frontend

```bash
cd frontend
npm start  # Inicia dev server com hot reload na porta 3000 (alternativamente)
# OU editar e fazer:
npm run build
cd build && python3 -m http.server 3000
```

### Modificar Go API

```bash
cd services/inventory
go run .  # Auto-recompila com mudanças (requer `go install github.com/cosmtrek/air@latest`)
```

### Modificar OCR Service

```bash
cd services/ocr
python3 main_mock.py  # Reinicia manualmente após mudanças
```

## 📊 Estrutura de Pastas

```
/home/lucasbastos/ComputationalVision/
├── frontend/                    # React app
│   ├── src/
│   ├── build/                  # Build estático (produção)
│   ├── package.json
│   └── tsconfig.json
├── services/
│   ├── inventory/              # Go API
│   │   ├── main.go
│   │   ├── handlers.go
│   │   └── models.go
│   └── ocr/                    # OCR Service (Python)
│       ├── main_mock.py
│       └── requirements.txt
├── docker-compose.yml          # (Opcional, para produção)
└── LOCAL_SETUP.md             # Este arquivo
```

## 🚨 Troubleshooting

### Erro: "connection refused" na porta 5434

**Solução:** PostgreSQL não está rodando. Execute:
```bash
docker run -d --name postgres-dev -e POSTGRES_USER=inventory -e POSTGRES_PASSWORD=password -e POSTGRES_DB=inventory -p 5434:5432 postgres:15-alpine
```

### Erro: "port already in use"

**Solução:** Encontre o processo:
```bash
# Porta 3000 (Frontend)
lsof -i :3000 | tail -1 | awk '{print $2}' | xargs kill

# Porta 8081 (Go API)
lsof -i :8081 | tail -1 | awk '{print $2}' | xargs kill

# Porta 5001 (OCR)
lsof -i :5001 | tail -1 | awk '{print $2}' | xargs kill
```

### Erro: "API_URL is undefined"

**Solução:** Frontend está tentando chamar API para porta errada. Verifique em `frontend/src/services/api.ts`:
- Deve detectar `localhost` e usar `http://localhost:8081`
- Build deve ser regenerado com: `npm run build`

### Frontend mostra "Nenhum item no estoque"

**Esperado** - Banco começa vazio. Para adicionar dados de teste:
```bash
curl -X POST http://localhost:8081/inventory/in \
  -H "Content-Type: application/json" \
  -d '{
    "part_number": "ABC-123",
    "serial_number": "SN-001",
    "quantity": 5,
    "location": "Rack-A1"
  }'
```

## 📝 Variáveis de Ambiente

### Go Service
```bash
DATABASE_URL="postgres://user:password@host:port/db?sslmode=disable"
PORT=8081  # Porta padrão
LOG_LEVEL="debug"  # Opcional
```

### OCR Service
```bash
PORT=5001  # Porta padrão
DEBUG=false  # Opcional
```

### Frontend (React)
```bash
REACT_APP_API_URL=http://localhost:8081  # Opcional (auto-detecta)
```

## 🔄 Fluxo de Uso

1. Abra http://localhost:3000 no navegador
2. Clique em "Capturar Foto" (mock em ambiente virtual)
3. Sistema envia imagem para OCR service (localhost:5001)
4. OCR retorna texto extraído com confiança
5. Frontend mostra texto para aprovação
6. Após aprovação, dados são salvos no Go API (localhost:8081)
7. Go API salva no PostgreSQL (localhost:5434)
8. Histórico é atualizado com novo item

## 📚 Próximos Passos

- [ ] Substituir OCR mock com PaddleOCR real
- [ ] Adicionar autenticação básica
- [ ] Criar seed de dados de teste
- [ ] Documentar endpoints da API
- [ ] Adicionar testes unitários
- [ ] Configurar CI/CD

## 💡 Tips

- Use `docker logs postgres-dev` para ver logs do banco
- Verifique `echo $DATABASE_URL` para confirmar variável configurada
- Use `psql` para conectar ao banco: `psql -h localhost -p 5434 -U inventory`
- Browser DevTools (F12) para debugar requisições frontend
- Postman/curl para testar endpoints manualmente

---

**Última atualização:** 2024
**Status:** ✅ Sistema completamente funcional localmente
