# Roadmap do Projeto

## Visão Geral

Este documento registra o histórico de decisões, implementações e progresso do projeto de entrada de estoque por visão computacional.

**Iniciado em:** 01/06/2026

---

## Estrutura de Documentação

```
docs/
├── 00-ROADMAP.md (este arquivo)
├── 01-ARQUITETURA.md
├── 02-TECNOLOGIAS.md
├── FASE-1-MVP.md
├── FASE-2-PARSER.md
├── FASE-3-PESQUISA.md
├── FASE-4-YOLO.md
└── FASE-5-PRODUCAO.md
```

---

## Log de Commits

### [INIT] 01/06/2026 - Estrutura Base

**Commit:** `git init && git add . && git commit -m "feat: estrutura base do projeto"`

**O que foi feito:**
- ✅ Criada estrutura de pastas para 6 microserviços
- ✅ Criados arquivos de configuração base (.gitignore, docker-compose.yml)
- ✅ Iniciada documentação em Obsidian
- ✅ README.md com overview do projeto

**Próximos Passos:**
- [ ] Fase 1: Setup do Frontend React
- [ ] Fase 1: Setup do OCR Service (Python)
- [ ] Fase 1: Setup do PostgreSQL

---

## Arquivos de Referência

- [[01-ARQUITETURA]] - Detalhes técnicos da arquitetura
- [[02-TECNOLOGIAS]] - Justificativa de cada tecnologia
- [[FASE-1-MVP]] - Detalhes da implementação do MVP
