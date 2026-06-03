# Fase 3 - Pesquisa Web Automatica

## Objetivo

Adicionar enriquecimento automatico de metadados quando o item nao for encontrado no catalogo local.

## Escopo Implementado

- Novo microservico `services/web-research` em Go.
- Endpoint `POST /research` para enriquecimento com base em Part Number.
- Endpoint `GET /research?pn=...` para testes rapidos.
- Integracao no frontend para chamada automatica quando `catalog/search` retorna nao encontrado.
- Exibicao de sugestao de fabricante, categoria e descricao normalizada na UI.

## Estrategia de Pesquisa

1. Executa busca web por `part number + datasheet + specifications + manufacturer`.
2. Extrai titulos, links e snippets dos principais resultados.
3. Aplica heuristicas de NLP simples para inferir:
   - fabricante
   - categoria
   - descricao normalizada
4. Calcula score de confianca e retorna sinais para auditoria.
5. Se a busca externa falhar, aplica fallback por heuristica local (prefixos e tokens conhecidos).

## API

### POST /research

Request:

```json
{
  "part_number": "M393A4K40DB3-CWE",
  "manufacturer": "",
  "category": "unknown",
  "normalized_description": "",
  "tokens": ["DDR4", "32GB", "RDIMM"]
}
```

Response (exemplo):

```json
{
  "success": true,
  "part_number": "M393A4K40DB3-CWE",
  "found": true,
  "manufacturer": "Samsung",
  "category": "memory",
  "normalized_description": "SAMSUNG MEMORY 32GB PC4-3200",
  "confidence": 0.82,
  "sources": [
    {
      "title": "Samsung M393A4K40DB3-CWE Datasheet",
      "url": "https://...",
      "snippet": "..."
    }
  ],
  "signals": ["manufacturer:web:Samsung", "category:web:memory", "description:web"]
}
```

## Execucao Local

```bash
cd services/web-research
$env:PORT='8083'
go run .
```

Healthcheck:

```bash
curl http://localhost:8083/health
```

## Testes

```bash
cd services/web-research
go test ./...
```
