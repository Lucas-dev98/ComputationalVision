#!/usr/bin/env python3
"""Treino YOLO para fase de visão (etiqueta/componente)."""

from __future__ import annotations

import argparse
from pathlib import Path
from ultralytics import YOLO


def main():
    parser = argparse.ArgumentParser(description="Train YOLO model")
    parser.add_argument("--model", default="yolov8n.pt")
    parser.add_argument("--data", default="services/vision/config/dataset.yaml")
    parser.add_argument("--epochs", type=int, default=100)
    parser.add_argument("--imgsz", type=int, default=640)
    parser.add_argument("--batch", type=int, default=8)
    parser.add_argument("--project", default="runs/vision")
    parser.add_argument("--name", default="yolo_continual")
    args = parser.parse_args()

    if not Path(args.data).exists():
        raise SystemExit(f"Arquivo de dataset não encontrado: {args.data}")

    model = YOLO(args.model)
    model.train(
        data=args.data,
        epochs=args.epochs,
        imgsz=args.imgsz,
        batch=args.batch,
        project=args.project,
        name=args.name,
    )


if __name__ == "__main__":
    main()
