import json
import tempfile
from pathlib import Path
import sys

from PIL import Image, ImageDraw
import requests

OCR_URL = "http://127.0.0.1:5001/ocr"
PARSER_URL = "http://127.0.0.1:8082/parse"

lines = [
    "SAMSUNG",
    "16GB 2Rx4 PC4-2133P",
    "M393A2G40B0B-CPB",
    "S/N: 1234ABC56789",
]


def make_test_image(path: Path) -> None:
    img = Image.new("RGB", (900, 360), color=(255, 255, 255))
    draw = ImageDraw.Draw(img)

    y = 40
    for line in lines:
        draw.text((40, y), line, fill=(0, 0, 0))
        y += 70

    img.save(path)


def main() -> None:
    with tempfile.TemporaryDirectory() as tmp_dir:
        image_path = Path(tmp_dir) / "memory_label.png"
        make_test_image(image_path)

        with image_path.open("rb") as f:
            ocr_resp = requests.post(OCR_URL, files={"file": (image_path.name, f, "image/png")}, timeout=120)
        ocr_resp.raise_for_status()
        ocr_data = ocr_resp.json()

        parse_payload = {"text": ocr_data.get("text", [])}
        parser_resp = requests.post(PARSER_URL, json=parse_payload, timeout=30)
        parser_resp.raise_for_status()
        parser_data = parser_resp.json()

        output = {
            "ocr_success": ocr_data.get("success"),
            "ocr_text": ocr_data.get("text", []),
            "ocr_confidence": ocr_data.get("confidence", []),
            "parser_result": parser_data,
        }
        output_json = json.dumps(output, ensure_ascii=False, indent=2)
        result_path = Path(__file__).with_name("e2e_result.json")
        result_path.write_text(output_json, encoding="utf-8")
        sys.stdout.write(output_json + "\n")
        sys.stdout.flush()


if __name__ == "__main__":
    main()
