# Escolha de Tecnologias

## Frontend

**Tecnologia:** React + TypeScript

**Motivos:**
- Excelente para captura de webcam
- Hot reload durante desenvolvimento
- Comunidade grande e bibliotecas maduras
- TypeScript para type safety
- Bom desempenho em tempo real

**Bibliotecas:**
- `react-webcam` - Captura de câmera
- `axios` - Requisições HTTP
- `antd` ou `shadcn/ui` - UI components

---

## Serviço de Estoque (Inventory Service)

**Tecnologia:** Go + PostgreSQL

**Motivos de Go:**
- Baixo consumo de memória (importante em produção)
- Alta concorrência nativa (goroutines)
- Excelente desempenho para APIs REST
- Binary único, fácil de deployar
- Rápido para operações de I/O (BD, web scraping)

**Responsável por:**
- CRUD de estoque
- Gestão de usuários
- Auditoria de movimentações
- Catálogo de Part Numbers
- Logs de todas as operações

---

## Serviço de Visão Computacional (Vision Service)

**Tecnologia:** Python + OpenCV + YOLO

**Motivos:**
- Praticamente todo ecossistema de IA/ML usa Python
- OpenCV é mature e bem documentado
- YOLO é state-of-the-art para object detection
- Comunidade imensa
- Performance adequada com GPU (CUDA/cuDNN)

**Bibliotecas:**
- `opencv-python` - Processamento de imagem
- `ultralytics` - YOLO v8
- `numpy` - Operações numéricas

---

## Serviço OCR

**Tecnologia:** Python + PaddleOCR

**Motivos:**
- PaddleOCR é mais rápido que Tesseract
- Suporta múltiplos idiomas bem
- Menos overhead que OCR via IA
- Python para consistência com Vision Service

**Bibliotecas:**
- `paddleocr` - OCR
- `opencv-python` - Pré/pós processamento

---

## Serviço de Parsing (Parser Service)

**Tecnologia:** Go

**Motivos:**
- Processamento textual simples e extremamente rápido
- Regex é eficiente em Go
- Bom para CPU-bound tasks
- Pode rodar como microserviço lightweight

**Responsável por:**
```
Entrada:  "32GB PC4-3200AA DDR4 UDIMM"
Saída:    {
  "type": "memory",
  "capacity": "32GB",
  "frequency": "3200MHz",
  "standard": "DDR4",
  "form_factor": "UDIMM"
}
```

---

## Serviço de Pesquisa Web (Web Research Service)

**Tecnologia:** Go

**Motivos:**
- Concorrência: muitas conexões simultâneas
- HTTP client nativo eficiente
- Scraping com `colly` é muito rápido
- Cache com Redis reduz requisições

**Fluxo:**
```
PN não encontrado
    ↓
Pesquisa internet
    ↓
Obtém especificações
    ↓
Normaliza
    ↓
Salva em catálogo
    ↓
Cache em Redis
```

---

## Banco de Dados

**Tecnologia:** PostgreSQL

**Motivos:**
- Relacional e bem estruturado para este caso
- ACID guarantees
- Suporta JSON (para dados semi-estruturados)
- Excellent para auditoria com timestamps
- Performance para queries complexas

**Tabelas principais:**
- `catalog` - Part Numbers e especificações
- `inventory` - Estoque atual
- `movements` - Histórico de movimentações
- `users` - Auditoria
- `audit_log` - Log de todas as operações

---

## Cache

**Tecnologia:** Redis

**Motivos:**
- Cache em memória para Part Numbers já pesquisados
- Session storage rápido
- Pub/Sub para notificações em tempo real
- TTL automático para dados expirados

---

## Containerização

**Tecnologia:** Docker + Docker Compose (dev), Kubernetes (prod)

**Motivos:**
- Consistência entre dev e produção
- Fácil onboarding de novos devs
- Cada microserviço em container isolado
- Kubernetes para scaling automático

---

## Observabilidade

**Tecnologia:** Prometheus + Grafana + Loki

**Motivos:**
- Prometheus para métricas
- Grafana para visualização
- Loki para logs centralizados
- Stack maduro e well-supported

---

## Resumo de Stack

| Componente | Tecnologia | Justificativa |
|-----------|-----------|---------------|
| Frontend | React + TS | Webcam, tempo real, produtividade |
| Inventory | Go | Performance, concorrência, APIs |
| Vision | Python | Ecossistema ML, OpenCV, YOLO |
| OCR | Python | PaddleOCR, consistência |
| Parser | Go | Velocidade, regex eficiente |
| Web Research | Go | Concorrência, scraping, performance |
| Banco | PostgreSQL | Relacional, ACID, JSON, auditoria |
| Cache | Redis | In-memory, TTL, pub/sub |
| Container | Docker | Consistência, isolation |
| Orquestração | K8s (opcional) | Scaling, HA |
