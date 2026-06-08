# Vision Service (YOLO) - Contínuo Learning

Este diretório concentra o pipeline de detecção visual para Fase 4.

## Estrutura

- `config/dataset.example.yaml`: template de dataset YOLO.
- `requirements.txt`: dependências para treino/avaliação.
- `models/`: registro local de modelo ativo (`active.pt`).

## Pipeline

1. Exportar feedback humano da API:

```bash
python scripts/vision/export_active_learning.py --inventory-url http://localhost:8081 --limit 500 --corrections-only
```

2. Organizar dataset em formato YOLO.
3. Treinar candidato:

```bash
python scripts/vision/train_yolo.py --model yolov8n.pt --data services/vision/config/dataset.yaml
```

4. Avaliar candidato:

```bash
python scripts/vision/evaluate_yolo.py --weights runs/vision/yolo_continual/weights/best.pt --data services/vision/config/dataset.yaml
```

5. Promover modelo por gate de métricas:

```bash
python scripts/vision/promote_model.py --candidate-weights runs/vision/yolo_continual/weights/best.pt --candidate-metrics runs/vision/latest_eval.json
```

## Observação

A promoção só deve ocorrer quando o candidato superar baseline em mAP50 sem regressão nas classes críticas.
