#!/usr/bin/env python3
"""Avaliação de modelo YOLO e geração de resumo JSON."""

from __future__ import annotations

import argparse
import json
from pathlib import Path
from ultralytics import YOLO


def main():
    parser = argparse.ArgumentParser(description="Evaluate YOLO model")
    parser.add_argument("--weights", required=True)
    parser.add_argument("--data", default="services/vision/config/dataset.yaml")
    parser.add_argument("--imgsz", type=int, default=640)
    parser.add_argument("--out", default="runs/vision/latest_eval.json")
    args = parser.parse_args()

    model = YOLO(args.weights)
    metrics = model.val(data=args.data, imgsz=args.imgsz)

    summary = {
        "weights": args.weights,
        "map50": float(getattr(metrics.box, "map50", 0.0)),
        "map50_95": float(getattr(metrics.box, "map", 0.0)),
        "precision": float(getattr(metrics.box, "mp", 0.0)),
        "recall": float(getattr(metrics.box, "mr", 0.0)),
    }

    out = Path(args.out)
    out.parent.mkdir(parents=True, exist_ok=True)
    out.write_text(json.dumps(summary, indent=2), encoding="utf-8")
    print(json.dumps(summary, indent=2))


if __name__ == "__main__":
    main()
