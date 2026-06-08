#!/usr/bin/env python3
"""Exporta feedback de aprendizado ativo do Inventory API para dataset de visão."""

from __future__ import annotations

import argparse
import json
from pathlib import Path
import urllib.request


def fetch_feedback(base_url: str, limit: int, corrections_only: bool):
    query = f"{base_url}/feedback/active-learning?limit={limit}&corrections_only={'true' if corrections_only else 'false'}"
    with urllib.request.urlopen(query, timeout=30) as resp:
        return json.loads(resp.read().decode("utf-8"))


def main():
    parser = argparse.ArgumentParser(description="Export active-learning samples")
    parser.add_argument("--inventory-url", default="http://localhost:8081")
    parser.add_argument("--limit", type=int, default=500)
    parser.add_argument("--corrections-only", action="store_true", default=True)
    parser.add_argument("--out-dir", default="datasets/vision/active_learning")
    args = parser.parse_args()

    payload = fetch_feedback(args.inventory_url, args.limit, args.corrections_only)
    items = payload.get("items", [])

    out_dir = Path(args.out_dir)
    out_dir.mkdir(parents=True, exist_ok=True)

    with (out_dir / "feedback_samples.jsonl").open("w", encoding="utf-8") as f:
        for item in items:
            f.write(json.dumps(item, ensure_ascii=False) + "\n")

    # Exporta imagem base64 para arquivo quando disponível.
    images_dir = out_dir / "images"
    images_dir.mkdir(exist_ok=True)
    exported_images = 0
    for item in items:
        image_data = item.get("image_data") or ""
        if not image_data.startswith("data:image"):
            continue
        try:
            header, b64 = image_data.split(",", 1)
            ext = "jpg"
            if "png" in header:
                ext = "png"
            image_path = images_dir / f"sample_{item.get('id')}.{ext}"
            image_path.write_bytes(__import__("base64").b64decode(b64))
            exported_images += 1
        except Exception:
            continue

    print(f"Amostras exportadas: {len(items)}")
    print(f"Imagens exportadas: {exported_images}")
    print(f"Saída: {out_dir}")


if __name__ == "__main__":
    main()
