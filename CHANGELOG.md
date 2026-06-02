# Changelog

Todas as mudanças notáveis neste projeto serão documentadas neste arquivo.

## [1.0.0] - 2026-06-01

### Adicionado - Fase 1 (MVP)

#### Frontend
- Componente WebcamCapture para captura de fotos via webcam
- Componente ImagePreview para visualização de imagens capturadas
- Componente ApprovalForm para aprovação manual de itens
- Componente HistoryTable com histórico de estoque
- Serviço de API com integração com backend
- Interface com Ant Design e TypeScript
- Responsividade para desktop

#### OCR Service
- Endpoint POST /ocr para extração de texto
- Endpoint POST /ocr/batch para processamento em lote
- Suporte a português e inglês
- Retorno estruturado (texto, confiança, caixas)
- PaddleOCR integration

#### Inventory Service
- Endpoint GET /health para health check
- Endpoint GET /catalog/search para busca de Part Number
- Endpoint POST /inventory/in para entrada de estoque
- Endpoint GET /inventory/items para listar estoque
- Endpoint GET /inventory/items/{id} para obter item
- PostgreSQL integration
- CORS habilitado

#### Database
- Tabela catalog (Part Numbers e especificações)
- Tabela inventory (Itens em estoque)
- Tabela movements (Histórico de movimentações)
- Tabela audit_log (Logs de auditoria)
- Índices otimizados
- Script de inicialização com dados de teste

#### Infraestrutura
- Docker Compose com orquestração completa
- Dockerfile para Frontend, OCR Service, Inventory Service
- PostgreSQL container com inicialização automática
- Redis container
- Network interno compartilhado
- Health checks

#### Documentação
- README.md com visão geral
- QUICKSTART.md com guia de início rápido
- Documentação arquitetura em `/docs/01-ARQUITETURA.md`
- Justificativa de tecnologias em `/docs/02-TECNOLOGIAS.md`
- Roadmap de fases em `/docs/00-ROADMAP.md`
- Detalhes de Fase 1 em `/docs/FASE-1-MVP.md`
- Histórico de commits em `/docs/COMMITS.md`
- Makefile com comandos úteis
- .env.example para configuração
- CONTRIBUTING.md com guia de contribuição

### Próximas Fases

- **Fase 2:** Parser + Classificação automática
- **Fase 3:** Pesquisa web automática
- **Fase 4:** YOLO + Detecção visual
- **Fase 5:** Deploy em produção

---

## Versionamento

Este projeto segue [Semantic Versioning](https://semver.org/).

### Formato
- MAJOR - mudanças incompatíveis
- MINOR - novas funcionalidades compatíveis
- PATCH - correções de bugs

Exemplo: v1.2.3

---

## Como Contribuir

Ver [CONTRIBUTING.md](CONTRIBUTING.md) para instruções detalhadas.
