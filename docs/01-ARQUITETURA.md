# Arquitetura do Sistema

## Visão Geral

```
┌─────────────────┐
│     Câmera      │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Vision Service  │     (Fase 4)
│ OpenCV + YOLO   │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│   OCR Service   │     (Fase 1)
│   PaddleOCR     │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Parsing Service │     (Fase 2)
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Catalog Engine  │     (Fase 2/3)
└────────┬────────┘
         │
    ┌────┴────┐
    ▼         ▼
┌────────┐ ┌──────────┐
│ Local  │ │ Web      │
│ DB     │ │ Research │
└────────┘ └──────────┘
    │         │
    └────┬────┘
         ▼
┌─────────────────┐
│ Classification  │
│ Engine          │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Approval UI     │     (Fase 1)
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Inventory API   │     (Fase 1)
└─────────────────┘
```

## Componentes

### 1. Frontend (React + TypeScript)
- Capture de imagem via webcam
- Preview da imagem capturada
- Formulário de aprovação manual
- Dashboard de histórico

### 2. OCR Service (Python + PaddleOCR)
- Extrai texto das imagens
- Identifica PN, SN e specs técnicas
- Normaliza output

### 3. Parser Service (Go)
- Classifica memórias (DDR3/4/5, capacidade, etc)
- Classifica discos (SATA/SAS/NVMe, interface, velocidade)
- Classifica rede (RJ45/SFP, velocidade)
- Retorna estrutura normalizada

### 4. Catalog Service (Go + PostgreSQL)
- Banco de dados de part numbers
- Integração com catálogo local
- Web scraping para novos PNs

### 5. Inventory API (Go + PostgreSQL)
- CRUD de itens
- Registro de movimentações
- Auditoria

## Fluxo de Dados - Fase 1 (MVP)

```
1. Usuário tira foto do componente
2. Frontend envia para OCR Service
3. OCR extrai PN, SN e especificações
4. Frontend exibe sugestão de componente
5. Usuário aprova ou corrige manualmente
6. Frontend envia para Inventory API
7. Inventory API registra entrada em estoque
8. Sistema retorna confirmação
```

## Fluxo de Dados - Após Fase 4 (Com YOLO)

```
1. Câmera captura imagem em tempo real
2. YOLO detecta tipo de componente e localiza etiqueta
3. OCR extrai PN do local detectado
4. Parser classifica automaticamente
5. Catalog busca especificações (local ou web)
6. Sistema registra entrada automaticamente (com aprovação manual opcional)
```
