# Guia de Início Rápido

## Pré-requisitos

- Docker e Docker Compose instalados
- Git
- (Opcional) Go 1.21+, Python 3.11+, Node 20+

## Iniciar o Projeto

### 1. Clone/Acesse o repositório

```bash
cd /home/lucasbastos/ComputationalVision
```

### 2. Inicie todos os serviços

```bash
docker-compose -f infra/docker/docker-compose.yml up -d
```

### 3. Aguarde a inicialização

```bash
# Verificar status
docker-compose -f infra/docker/docker-compose.yml ps
```

### 4. Acesse a aplicação

- **Frontend:** http://localhost:3000
- **Inventory API:** http://localhost:8080
- **OCR API:** http://localhost:5001
- **PostgreSQL:** localhost:5432 (credenciais: inventory/inventory_dev)
- **Redis:** localhost:6379

## Endpoints Principais

### Inventory API

#### Buscar Part Number
```bash
curl http://localhost:8080/catalog/search?pn=M393A4K40DB3-CWE
```

#### Adicionar ao Estoque
```bash
curl -X POST http://localhost:8080/inventory/in \
  -H "Content-Type: application/json" \
  -d '{
    "part_number": "M393A4K40DB3-CWE",
    "serial_number": "SN12345",
    "quantity": 1,
    "location": "DC-RJ",
    "reason": "Entrada recebida",
    "user_id": 1
  }'
```

#### Listar Estoque
```bash
curl http://localhost:8080/inventory/items?limit=50&offset=0
```

### OCR Service

#### Extrair Texto de Imagem
```bash
curl -X POST http://localhost:5001/ocr \
  -F "file=@/path/to/image.jpg"
```

## Desenvolvimento Local (sem Docker)

### OCR Service
```bash
cd services/ocr
pip install -r requirements.txt
python main.py
```

### Inventory Service
```bash
cd services/inventory
go run .
```

### Frontend
```bash
cd frontend
npm install
npm start
```

## Estrutura de Pastas

```
ComputationalVision/
├── frontend/               # React + TypeScript
│   ├── src/
│   ├── public/
│   ├── package.json
│   └── Dockerfile
├── services/
│   ├── ocr/               # OCR Service (Python)
│   ├── inventory/         # Inventory Service (Go)
│   ├── vision/            # Vision Service (Python) - Fase 4
│   ├── parser/            # Parser Service (Go) - Fase 2
│   ├── catalog/           # Catalog Service (Go) - Fase 2
│   └── web-research/      # Web Research Service (Go) - Fase 3
├── infra/
│   ├── docker/
│   │   ├── docker-compose.yml
│   │   └── postgres/
│   │       └── init.sql
│   └── kubernetes/
├── docs/                  # Documentação
│   ├── 00-ROADMAP.md
│   ├── 01-ARQUITETURA.md
│   ├── 02-TECNOLOGIAS.md
│   ├── FASE-1-MVP.md
│   └── COMMITS.md
├── .gitignore
└── README.md
```

## Troubleshooting

### Erro ao conectar ao PostgreSQL
```bash
# Verificar se PostgreSQL está rodando
docker-compose -f infra/docker/docker-compose.yml logs postgres

# Aguarde alguns segundos e tente novamente
```

### Porta já em uso
```bash
# Encontrar processo usando a porta
lsof -i :3000  # para frontend
lsof -i :8080  # para API
lsof -i :5001  # para OCR
```

### Limpar tudo e recomeçar
```bash
docker-compose -f infra/docker/docker-compose.yml down -v
docker-compose -f infra/docker/docker-compose.yml up -d
```

## Próximos Passos

1. **Fase 1 Conclusão:** Testes e integração completa
2. **Fase 2:** Parser e classificação automática
3. **Fase 3:** Pesquisa web automática
4. **Fase 4:** YOLO e detecção visual
5. **Fase 5:** Deploy em produção

## Documentação Completa

Ver pasta `/docs` para documentação detalhada de cada componente.
