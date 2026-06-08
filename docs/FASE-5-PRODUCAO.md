# Fase 5 - Produção

## Objetivo

Adicionar protecoes operacionais minimas para preparar os servicos para uso fora do ambiente estritamente local.

## Entrega Atual

### 1. Rate Limiting Basico por IP

Os servicos Go abaixo agora possuem middleware de rate limiting configuravel:

- `services/inventory`
- `services/parser`
- `services/web-research`

Variaveis de ambiente suportadas:

```bash
RATE_LIMIT_ENABLED=true
RATE_LIMIT_REQUESTS=120
RATE_LIMIT_WINDOW_SECONDS=60
```

Comportamento:

- Limita requisicoes por IP dentro de uma janela fixa
- Retorna `429 Too Many Requests` quando excede o limite
- Expõe `Retry-After`, `X-RateLimit-Limit` e `X-RateLimit-Remaining`
- Nao limita `/health` nem `OPTIONS`

## Validacao Executada

- `go test ./...` em `services/parser`
- `go test ./...` em `services/inventory`
- `go test ./...` em `services/web-research`
- Smoke test em runtime no parser com `RATE_LIMIT_REQUESTS=1`, validando `200` na primeira chamada e `429` na segunda

## Proximos Incrementos Naturais

1. Adicionar autenticacao para rotas mutaveis (`inventory/in`, feedback, futuras rotas administrativas)
2. Expor metricas Prometheus (`/metrics`) para latencia, erros e OCR timeout
3. Adicionar persistencia/cleanup da estrutura de contadores para cenarios de longa execucao
4. Aplicar a mesma protecao ao OCR service FastAPI