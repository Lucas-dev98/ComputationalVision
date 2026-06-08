# Fase 4 - Contínuo Learning (YOLO + Feedback Humano)

## Objetivo

Criar um ciclo contínuo de aprendizado no projeto:

1. Predição automática (YOLO + OCR + Parser)
2. Correção humana no frontend
3. Persistência de feedback no backend
4. Exportação de amostras difíceis
5. Retreino periódico da YOLO
6. Promoção de modelo por gate de métricas

## O que foi implementado

### Backend (Inventory API)

- `POST /feedback/submit`
- `GET /feedback/active-learning?limit=...&corrections_only=true|false`
- Tabela `feedback_samples` em SQLite e PostgreSQL

Campos salvos por amostra:

- predicted vs final para PN/SN/fabricante/categoria
- flag de correção humana (`correction_applied`)
- confiança
- texto OCR
- imagem (base64 opcional)
- metadados JSON

### Frontend

Ao confirmar entrada de estoque, o frontend envia feedback automaticamente para o endpoint de aprendizado contínuo.

### Pipeline de visão

Scripts criados:

- `scripts/vision/export_active_learning.py`
- `scripts/vision/train_yolo.py`
- `scripts/vision/evaluate_yolo.py`
- `scripts/vision/promote_model.py`

Config:

- `services/vision/config/dataset.example.yaml`
- `services/vision/requirements.txt`

## Fluxo operacional recomendado

1. Rodar app local (OCR, parser, web-research, inventory, frontend)
2. Operar normalmente e coletar feedback humano
3. Exportar feedback para active learning
4. Rotular/organizar dataset de visão
5. Treinar novo candidato YOLO
6. Avaliar métricas
7. Promover apenas se superar baseline

## Comandos

```bash
make vision-export
make vision-train
make vision-eval
make vision-promote
```

## Regras de segurança para promoção

- Nunca promover modelo sem avaliação em holdout
- Exigir ganho mínimo de mAP50 (`--min-map50-gain`)
- Manter modelo ativo anterior para rollback rápido

## Próximos ajustes sugeridos

- CI para validar scripts e endpoint de feedback em pull requests
- Versão de dataset com changelog por treino
- Dashboard simples de taxa de correção por campo (PN/SN)
