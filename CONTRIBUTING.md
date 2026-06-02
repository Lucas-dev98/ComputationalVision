# Guia de Contribuição

## Bem-vindo!

Este projeto está em desenvolvimento ativo. Aqui está como você pode contribuir.

## Como Começar

1. **Clone o repositório**
   ```bash
   git clone https://github.com/lucas-dev98/ComputationalVision.git
   cd ComputationalVision
   ```

2. **Configure o ambiente**
   ```bash
   cp .env.example .env
   make install
   ```

3. **Inicie os serviços**
   ```bash
   make up
   ```

## Workflow de Desenvolvimento

### Branch Naming
- `feature/` - Novas funcionalidades
- `fix/` - Correções de bugs
- `docs/` - Documentação
- `refactor/` - Refatoração
- `test/` - Testes

Exemplo:
```bash
git checkout -b feature/parser-service
```

### Commit Messages
Seguimos [Conventional Commits](https://www.conventionalcommits.org/)

```
feat: adicionar suporte a DDR5
fix: corrigir conexão com PostgreSQL
docs: atualizar guia de instalação
refactor: simplificar lógica de OCR
test: adicionar testes do parser
```

### Antes de Fazer um Commit

1. **Teste localmente**
   ```bash
   make test
   ```

2. **Formatar código**
   ```bash
   make format
   ```

3. **Verificar linter**
   ```bash
   make lint
   ```

## Estrutura de Código

### Frontend (React/TypeScript)
- Use componentes funcionais com hooks
- Tipagem obrigatória
- Componentes no `src/components/`
- Serviços no `src/services/`

### Backend (Go)
- Use idiomas Go
- `gofmt` para formatação
- Handlers em `handlers.go`
- Models em `models.go`

### Python (OCR, Vision)
- Use `black` para formatação
- Type hints obrigatórios
- Docstrings em todas as funções

## Pull Requests

1. **Descrição clara**
   - O que foi mudado?
   - Por quê?
   - Como testar?

2. **Commits pequenos e lógicos**
   - Um commit por funcionalidade
   - Histórico claro

3. **Testes inclusos**
   - Testes unitários
   - Testes de integração
   - Testes E2E (se aplicável)

## Relatando Bugs

Crie uma issue com:
- Descrição clara do problema
- Passos para reproduzir
- Resultado esperado vs atual
- Environment (OS, versões, etc)

## Sugestões de Features

Abra uma discussão descrevendo:
- O problema que resolve
- Casos de uso
- Impacto na arquitetura

## Documentação

- Atualize README se mudar comportamento
- Adicione comentários no código complexo
- Mantenha `/docs` atualizado
- Atualize CHANGELOG.md

## Fases do Projeto

- **Fase 1 (MVP):** ✅ Completa - Webcam + OCR + Estoque
- **Fase 2:** Parser + Classificação automática
- **Fase 3:** Pesquisa web automática
- **Fase 4:** YOLO + Detecção visual
- **Fase 5:** Deploy em produção

Priorize trabalho em fases sequenciais.

## Comunicação

- Issues para bugs e features
- Discussões para perguntas
- Pull request reviews são construtivas

## Código de Conduta

- Respeito mútuo
- Inclusão
- Feedback construtivo
- Foco no projeto

## Dúvidas?

1. Verifique `/docs`
2. Veja issues relacionadas
3. Abra uma discussão

---

**Obrigado por contribuir! 🎉**
