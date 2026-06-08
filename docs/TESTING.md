# Guia de Testes

## Executar Todos os Testes

```bash
make test
```

## Testes do OCR Service

### Requisitos
- Python 3.11+
- Dependências instaladas: `pip install -r services/ocr/requirements.txt`

### Executar Testes OCR

```bash
cd services/ocr
python -m pytest test_ocr_service.py -v
```

Ou com unittest:
```bash
python -m unittest test_ocr_service.py -v
```

### Testes Inclusos
- ✓ Inicialização do OCR Service
- ✓ Extração de texto com imagem válida
- ✓ Validação de estrutura de resultado
- ✓ Extração de texto estruturado
- ✓ Tempo de processamento

---

## Testes do Inventory Service

### Requisitos
- Go 1.21+
- PostgreSQL rodando
- DATABASE_URL configurada

### Executar Testes

```bash
cd services/inventory
go test -v
```

Com coverage:
```bash
go test -v -cover
```

### Testes Inclusos
- ✓ Conexão com banco de dados
- ✓ Busca em catálogo (existente)
- ✓ Busca em catálogo (não encontrado)
- [ ] Adição de item ao estoque
- [ ] Listagem de estoque
- [ ] Movimentações

---

## Testes de Integração

### E2E Fase 3 (local, sem Docker)

Com os serviços locais em execução (OCR, parser, web-research, inventory e frontend):

```bash
make test-e2e
```

Ou diretamente:

```bash
python scripts/e2e_phase3.py
```

O script valida:
- health checks dos serviços
- chamada OCR
- parse de PN/SN
- busca no catálogo
- pesquisa web
- entrada em estoque (quando PN existe no catálogo)

### Frontend + Backend

```bash
# Terminal 1: Iniciar backend
make api-dev

# Terminal 2: Iniciar OCR
make ocr-dev

# Terminal 3: Iniciar frontend
make frontend-dev

# Terminal 4: Rodar testes E2E
cd frontend
npm run test:e2e
```

---

## Testes Manual

### 1. Testar OCR Service

```bash
# Criar imagem de teste
curl -X POST http://localhost:5001/ocr \
  -F "file=@test-image.jpg"
```

### 2. Testar Inventory Service

```bash
# Health check
curl http://localhost:8080/health

# Buscar Part Number
curl http://localhost:8080/catalog/search?pn=M393A4K40DB3-CWE

# Adicionar ao estoque
curl -X POST http://localhost:8080/inventory/in \
  -H "Content-Type: application/json" \
  -d '{
    "part_number": "M393A4K40DB3-CWE",
    "serial_number": "SN12345",
    "quantity": 1,
    "location": "DC-RJ"
  }'

# Listar estoque
curl http://localhost:8080/inventory/items
```

### 3. Testar Frontend

1. Abra http://localhost:3000
2. Clique em "Capturar Foto"
3. Aprove o item
4. Verifique em "Histórico"

---

## Cobertura de Testes Esperada

| Componente | Cobertura Esperada |
|-----------|-------------------|
| OCR Service | 80%+ |
| Inventory Service | 75%+ |
| Frontend Components | 70%+ |
| API Integration | 85%+ |

---

## GitHub Actions CI/CD

Testes rodam automaticamente em:
- Push para master
- Pull requests
- Agendado (diariamente)

Ver `.github/workflows/ci.yml`

---

## Debugging

### OCR Service
```bash
# Logs com debug
LOG_LEVEL=DEBUG python main.py

# Verificar PaddleOCR
python -c "from paddleocr import PaddleOCR; ocr = PaddleOCR(); print(ocr)"
```

### Inventory Service
```bash
# Rodar com debug
LOG_LEVEL=debug go run .

# Verificar conexão com BD
psql postgresql://inventory:inventory_dev@localhost:5432/inventory_db
```

### Frontend
```bash
# DevTools do navegador (F12)
# Redux DevTools para state management
# Network tab para ver requisições
```

---

## Performance Testing

### Load Testing
```bash
# Instalar Apache Bench
brew install httpd  # macOS
apt-get install apache2-utils  # Linux

# Testar endpoint
ab -n 1000 -c 10 http://localhost:8080/health
```

### Memory Profiling (Go)
```bash
go test -memprofile=mem.prof
go tool pprof mem.prof
```

---

## Relatórios de Teste

Gerar relatório de coverage:
```bash
go test -coverprofile=coverage.out
go tool cover -html=coverage.out
```

---

## Próximas Melhorias

- [ ] Aumentar cobertura OCR (+ 15%)
- [ ] Adicionar testes de carga
- [ ] Tests E2E com Cypress/Playwright
- [ ] Mutation testing
- [ ] Performance benchmarking
