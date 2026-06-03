# API Documentation

## Inventory Service API

Base URL: `http://localhost:8080`

Para desenvolvimento local sem PostgreSQL, o Inventory Service pode rodar com SQLite:

```bash
cd services/inventory
$env:DATABASE_DRIVER='sqlite'
$env:DATABASE_URL='file:./inventory-dev.db?_pragma=foreign_keys(1)'
$env:PORT='8081'
go run .
```

---

## Health Check

### GET /health

Verifica se o serviço está operacional.

**Response:** 200 OK
```json
{
  "status": "healthy",
  "service": "inventory-service",
  "version": "1.0.0"
}
```

---

## Catálogo

### GET /catalog/search

Busca um Part Number no catálogo.

**Query Parameters:**
- `pn` (required): Part Number a buscar

**Example:**
```bash
curl "http://localhost:8080/catalog/search?pn=M393A4K40DB3-CWE"
```

**Response:** 200 OK
```json
{
  "found": true,
  "item": {
    "id": 1,
    "part_number": "M393A4K40DB3-CWE",
    "manufacturer": "Samsung",
    "category": "memory",
    "normalized_description": "DDR4 32GB 3200MHz RDIMM ECC",
    "created_at": "2026-06-01T10:00:00Z",
    "updated_at": "2026-06-01T10:00:00Z"
  }
}
```

**Response:** 404 Not Found
```json
{
  "found": false
}
```

---

## Estoque

### POST /inventory/in

Registra entrada de item em estoque.

**Request Body:**
```json
{
  "part_number": "M393A4K40DB3-CWE",
  "serial_number": "SN12345",
  "quantity": 1,
  "location": "DC-RJ",
  "reason": "Entrada recebida",
  "user_id": 1
}
```

**Response:** 201 Created
```json
{
  "success": true,
  "data": {
    "id": 5,
    "catalog_id": 1,
    "serial_number": "SN12345",
    "quantity": 1,
    "location": "DC-RJ",
    "status": "active",
    "received_at": "2026-06-01T10:30:00Z",
    "last_updated": "2026-06-01T10:30:00Z",
    "catalog": {
      "id": 1,
      "part_number": "M393A4K40DB3-CWE",
      "manufacturer": "Samsung",
      "category": "memory",
      "normalized_description": "DDR4 32GB 3200MHz RDIMM ECC"
    }
  }
}
```

**Error Response:** 500 Internal Server Error
```json
{
  "success": false,
  "error": "Part Number não encontrado"
}
```

---

### GET /inventory/items

Lista itens em estoque com paginação.

**Query Parameters:**
- `limit` (optional, default: 50, max: 500): Número de itens por página
- `offset` (optional, default: 0): Offset para paginação

**Example:**
```bash
curl "http://localhost:8080/inventory/items?limit=10&offset=0"
```

**Response:** 200 OK
```json
{
  "total": 150,
  "items": [
    {
      "id": 5,
      "catalog_id": 1,
      "serial_number": "SN12345",
      "quantity": 1,
      "location": "DC-RJ",
      "status": "active",
      "received_at": "2026-06-01T10:30:00Z",
      "last_updated": "2026-06-01T10:30:00Z",
      "catalog": {
        "id": 1,
        "part_number": "M393A4K40DB3-CWE",
        "manufacturer": "Samsung",
        "category": "memory",
        "normalized_description": "DDR4 32GB 3200MHz RDIMM ECC"
      }
    }
  ]
}
```

---

### GET /inventory/items/{id}

Obtém detalhes de um item específico.

**Path Parameters:**
- `id` (required): ID do item

**Example:**
```bash
curl "http://localhost:8080/inventory/items/5"
```

**Response:** 200 OK
```json
{
  "id": 5,
  "catalog_id": 1,
  "serial_number": "SN12345",
  "quantity": 1,
  "location": "DC-RJ",
  "status": "active",
  "received_at": "2026-06-01T10:30:00Z",
  "last_updated": "2026-06-01T10:30:00Z",
  "catalog": {
    "id": 1,
    "part_number": "M393A4K40DB3-CWE",
    "manufacturer": "Samsung",
    "category": "memory",
    "normalized_description": "DDR4 32GB 3200MHz RDIMM ECC"
  }
}
```

**Response:** 404 Not Found
```json
{
  "error": "Item não encontrado"
}
```

---

## OCR Service API

Base URL: `http://localhost:5001`

### Health Check

#### GET /health

**Response:** 200 OK
```json
{
  "status": "healthy",
  "service": "ocr-service",
  "version": "1.0.0"
}
```

---

### POST /ocr

Extrai texto de uma imagem.

**Request:**
- Content-Type: `multipart/form-data`
- Form Data: `file` (image file)

**Example:**
```bash
curl -X POST http://localhost:5001/ocr \
  -F "file=@/path/to/image.jpg"
```

**Response:** 200 OK
```json
{
  "success": true,
  "text": [
    "M393A4K40DB3-CWE",
    "32GB",
    "PC4-3200AA"
  ],
  "confidence": [0.98, 0.99, 0.97],
  "boxes": [
    [100, 50, 200, 80],
    [100, 90, 200, 110],
    [100, 130, 200, 150]
  ],
  "processing_time_ms": 1234.5
}
```

**Error Response:** 500 Internal Server Error
```json
{
  "success": false,
  "error": "Erro ao processar imagem",
  "text": [],
  "confidence": [],
  "boxes": [],
  "processing_time_ms": 100.2
}
```

---

### POST /ocr/batch

Processa múltiplas imagens em lote.

**Request:**
- Content-Type: `multipart/form-data`
- Form Data: `files` (múltiplos arquivos de imagem)

**Example:**
```bash
curl -X POST http://localhost:5001/ocr/batch \
  -F "files=@image1.jpg" \
  -F "files=@image2.jpg"
```

**Response:** 200 OK
```json
{
  "results": [
    {
      "filename": "image1.jpg",
      "success": true,
      "text": ["M393A4K40DB3-CWE"],
      "confidence": [0.98],
      "boxes": [[100, 50, 200, 80]],
      "processing_time_ms": 1200
    },
    {
      "filename": "image2.jpg",
      "success": true,
      "text": ["HUH721212AL5200"],
      "confidence": [0.99],
      "boxes": [[150, 60, 250, 90]],
      "processing_time_ms": 1150
    }
  ],
  "total": 2
}
```

---

## Códigos de Status HTTP

| Código | Significado |
|--------|------------|
| 200 | OK - Requisição bem-sucedida |
| 201 | Created - Recurso criado |
| 400 | Bad Request - Requisição inválida |
| 404 | Not Found - Recurso não encontrado |
| 500 | Internal Server Error - Erro do servidor |

---

## Rate Limiting

Não implementado em Fase 1. Será adicionado em produções.

---

## Autenticação

Não implementado em Fase 1. Será adicionado em Fase 5.

---

## Versionamento de API

Versão atual: `1.0.0`

Mudanças breaking serão precedidas de nova versão (ex: `/v2/catalog/search`)
