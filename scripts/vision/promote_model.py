#!/usr/bin/env python3
"""Promove novo modelo YOLO se métricas superarem baseline."""

from __future__ import annotations

import argparse
import json
from pathlib import Path
import shutil


def read_json(path: Path):
    if not path.exists():
        return {}
    return json.loads(path.read_text(encoding="utf-8"))


def main():
    parser = argparse.ArgumentParser(description="Promote model by metric gate")
    parser.add_argument("--candidate-weights", required=True)
    parser.add_argument("--candidate-metrics", required=True)
    parser.add_argument("--baseline-metrics", default="runs/vision/baseline_eval.json")
    parser.add_argument("--registry-dir", default="services/vision/models")
    parser.add_argument("--min-map50-gain", type=float, default=0.0)
    args = parser.parse_args()

    candidate_metrics = read_json(Path(args.candidate_metrics))
    baseline_metrics = read_json(Path(args.baseline_metrics))

    candidate_map50 = float(candidate_metrics.get("map50", 0.0))
    baseline_map50 = float(baseline_metrics.get("map50", 0.0))

    if candidate_map50 < baseline_map50 + args.min_map50_gain:
        print("Modelo candidato não atingiu gate de promoção.")
        print(f"candidate map50={candidate_map50:.4f}, baseline map50={baseline_map50:.4f}")
        raise SystemExit(2)

    registry = Path(args.registry_dir)
    registry.mkdir(parents=True, exist_ok=True)

    active_weights = registry / "active.pt"
    active_metrics = registry / "active.metrics.json"

    shutil.copyfile(args.candidate_weights, active_weights)
    active_metrics.write_text(json.dumps(candidate_metrics, indent=2), encoding="utf-8")

    print("Modelo promovido com sucesso:")
    print(f"- weights: {active_weights}")
    print(f"- metrics: {active_metrics}")


if __name__ == "__main__":
    main()
