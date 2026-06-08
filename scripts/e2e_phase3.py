#!/usr/bin/env python3
"""E2E local para validar fluxo da Fase 3 sem dependencias externas.

Valida:
- health de OCR, parser, web-research e inventory
- OCR endpoint com imagem PNG embutida
- parser com texto conhecido (PN/SN)
- pesquisa no catalogo
- pesquisa web automatica
- entrada em estoque (quando PN existe no catalogo)
"""

from __future__ import annotations

import argparse
import base64
import json
import sys
import time
import urllib.error
import urllib.parse
import urllib.request


ONE_PIXEL_PNG_BASE64 = (
    "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR4nGNgYAAAAAMA"
    "ASsJTYQAAAAASUVORK5CYII="
)


def http_json(method: str, url: str, body: dict | None = None, timeout: int = 20):
    data = None
    headers = {"Accept": "application/json"}
    if body is not None:
        data = json.dumps(body).encode("utf-8")
        headers["Content-Type"] = "application/json"

    req = urllib.request.Request(url=url, method=method.upper(), data=data, headers=headers)
    with urllib.request.urlopen(req, timeout=timeout) as resp:
        payload = resp.read().decode("utf-8")
        return resp.status, json.loads(payload) if payload else {}


def http_multipart_file(url: str, field_name: str, filename: str, content: bytes, timeout: int = 90):
    boundary = "----CVBoundary7MA4YWxkTrZu0gW"
    body = (
        f"--{boundary}\r\n"
        f"Content-Disposition: form-data; name=\"{field_name}\"; filename=\"{filename}\"\r\n"
        "Content-Type: image/png\r\n\r\n"
    ).encode("utf-8") + content + f"\r\n--{boundary}--\r\n".encode("utf-8")

    req = urllib.request.Request(
        url=url,
        method="POST",
        data=body,
        headers={
            "Content-Type": f"multipart/form-data; boundary={boundary}",
            "Accept": "application/json",
        },
    )

    with urllib.request.urlopen(req, timeout=timeout) as resp:
        payload = resp.read().decode("utf-8")
        return resp.status, json.loads(payload) if payload else {}


def expect(name: str, condition: bool, details: str = ""):
    if condition:
        print(f"[OK] {name}")
        return
    suffix = f" -> {details}" if details else ""
    print(f"[FAIL] {name}{suffix}")
    raise AssertionError(f"Falha: {name}{suffix}")


def main():
    parser = argparse.ArgumentParser(description="E2E Fase 3 (local)")
    parser.add_argument("--ocr", default="http://localhost:5001", help="Base URL do OCR service")
    parser.add_argument("--parser", dest="parser_url", default="http://localhost:8082", help="Base URL do parser service")
    parser.add_argument("--web", default="http://localhost:8083", help="Base URL do web-research service")
    parser.add_argument("--inventory", default="http://localhost:8081", help="Base URL do inventory service")
    args = parser.parse_args()

    # 1) Health checks
    for name, base in [
        ("ocr", args.ocr),
        ("parser", args.parser_url),
        ("web-research", args.web),
        ("inventory", args.inventory),
    ]:
        status, payload = http_json("GET", f"{base}/health", timeout=15)
        expect(f"health {name}", status == 200 and payload.get("status") == "healthy", str(payload))

    # 2) OCR request (imagem minima embutida)
    png_bytes = base64.b64decode(ONE_PIXEL_PNG_BASE64)
    ocr_status, ocr_payload = http_multipart_file(f"{args.ocr}/ocr", "file", "pixel.png", png_bytes, timeout=120)
    expect("ocr status", ocr_status == 200, str(ocr_payload))
    expect("ocr formato", isinstance(ocr_payload.get("success"), bool), str(ocr_payload))

    # 3) Parser com texto de referencia conhecido
    sample_lines = [
        "M393A4K40DB3-CWE",
        "Samsung DDR4 32GB 3200MHz RDIMM ECC",
        "SN: SN-1234567890",
    ]
    p_status, p_payload = http_json("POST", f"{args.parser_url}/parse", {"text": sample_lines}, timeout=20)
    expect("parser status", p_status == 200, str(p_payload))
    expect("parser success", p_payload.get("success") is True, str(p_payload))
    expect("parser pn", p_payload.get("part_number") == "M393A4K40DB3-CWE", str(p_payload))
    expect("parser sn", p_payload.get("serial_number") == "SN-1234567890", str(p_payload))

    pn = p_payload["part_number"]

    # 4) Catalog lookup
    inv_query = urllib.parse.quote(pn)
    i_status, i_payload = http_json("GET", f"{args.inventory}/catalog/search?pn={inv_query}", timeout=15)
    expect("catalog lookup status", i_status in (200, 404), str(i_payload))

    # 5) Web research fallback/enrichment
    w_status, w_payload = http_json(
        "POST",
        f"{args.web}/research",
        {
            "part_number": pn,
            "manufacturer": p_payload.get("manufacturer", ""),
            "category": p_payload.get("category", ""),
            "normalized_description": p_payload.get("normalized_description", ""),
            "tokens": sample_lines,
        },
        timeout=30,
    )
    expect("web research status", w_status == 200, str(w_payload))
    expect("web research success", w_payload.get("success") is True, str(w_payload))

    # 6) Inventory IN quando item existe no catalogo
    if i_status == 200 and i_payload.get("found") is True:
        serial = f"SN-E2E-{int(time.time())}"
        in_status, in_payload = http_json(
            "POST",
            f"{args.inventory}/inventory/in",
            {
                "part_number": pn,
                "serial_number": serial,
                "quantity": 1,
                "location": "E2E-LOCAL",
                "reason": "Teste E2E Fase 3",
                "user_id": 1,
            },
            timeout=20,
        )
        expect("inventory in status", in_status == 201, str(in_payload))
        expect("inventory in success", in_payload.get("success") is True, str(in_payload))
    else:
        print("[WARN] PN nao encontrado no catalogo local; pulando etapa inventory/in")

    print("\nE2E Fase 3 concluido com sucesso.")


if __name__ == "__main__":
    try:
        main()
    except urllib.error.HTTPError as exc:
        details = exc.read().decode("utf-8", errors="ignore")
        print(f"[HTTP ERROR] {exc.code} {exc.reason}: {details}")
        sys.exit(1)
    except Exception as exc:
        print(f"[ERROR] {exc}")
        sys.exit(1)
