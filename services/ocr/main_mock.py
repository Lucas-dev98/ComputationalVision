#!/usr/bin/env python3
"""
OCR Service Mock - Para testes locais
Retorna dados simulados de OCR
"""
from fastapi import FastAPI, File, UploadFile, HTTPException
from fastapi.middleware.cors import CORSMiddleware
import re
import time

app = FastAPI(
    title="OCR Service",
    description="Serviço de OCR para extração de texto",
    version="1.0.0-mock"
)

# Adicionar CORS
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

MOCK_TEXT_LINES = [
    "M393A4K40DB3-CWE",
    "SN: SN-1234567890",
    "Samsung DDR4 32GB 3200MHz RDIMM ECC",
]


def _build_boxes(count):
    boxes = []
    for index in range(count):
        top = 10 + (index * 50)
        boxes.append(
            [
                [10, top],
                [320, top],
                [320, top + 30],
                [10, top + 30],
            ]
        )
    return boxes


def _find_part_numbers(text_lines):
    candidates = []
    pattern = re.compile(r"\b[A-Z0-9][A-Z0-9._/-]{5,}\b")

    for line in text_lines:
        for match in pattern.findall(line.upper()):
            if len(match) < 8:
                continue
            if not any(char.isdigit() for char in match):
                continue
            if match.startswith("SN") and "-" in match:
                continue
            if match not in candidates:
                candidates.append(match)

    return candidates


def _find_serial_numbers(text_lines):
    candidates = []
    serial_pattern = re.compile(r"\bSN[:\- ]*([A-Z0-9][A-Z0-9\-]{5,})\b")

    for line in text_lines:
        for match in serial_pattern.findall(line.upper()):
            if any(char.isdigit() for char in match) and match not in candidates:
                candidates.append(match)

    return candidates


def _build_mock_result():
    boxes = _build_boxes(len(MOCK_TEXT_LINES))
    detected_part_numbers = _find_part_numbers(MOCK_TEXT_LINES)
    detected_serial_numbers = _find_serial_numbers(MOCK_TEXT_LINES)

    return {
        "success": True,
        "text": MOCK_TEXT_LINES,
        "confidence": [0.98, 0.94, 0.92],
        "boxes": boxes,
        "processing_time_ms": 12.5,
        "mock": True,
        "detected_part_numbers": detected_part_numbers,
        "detected_serial_numbers": detected_serial_numbers,
        "structured": {
            "lines": [
                {
                    "text": text,
                    "confidence": confidence,
                    "bbox": bbox,
                }
                for text, confidence, bbox in zip(MOCK_TEXT_LINES, [0.98, 0.94, 0.92], boxes)
            ],
            "total_lines": len(MOCK_TEXT_LINES),
        },
    }


def _validate_image_upload(file):
    if not file.content_type or not file.content_type.startswith("image/"):
        raise HTTPException(status_code=400, detail="Arquivo não é uma imagem")

@app.get("/health")
async def health():
    return {
        "status": "healthy",
        "service": "ocr-service",
        "version": "1.0.0-mock"
    }

@app.post("/ocr")
async def extract_text(file: UploadFile = File(...)):
    """
    Endpoint simulado para extrair texto de uma imagem.
    Retorna dados mock para testes.
    """
    try:
        _validate_image_upload(file)
        await file.read()

        start_time = time.time()
        result = _build_mock_result()
        result["processing_time_ms"] = (time.time() - start_time) * 1000

        return result
        
    except HTTPException:
        raise
    except Exception as e:
        raise HTTPException(
            status_code=500,
            detail=f"Erro ao processar imagem: {str(e)}"
        )


@app.post("/ocr/batch")
async def extract_text_batch(files: list[UploadFile] = File(...)):
    """
    Endpoint simulado para OCR em lote.
    """
    try:
        results = []

        for file in files:
            if not file.content_type or not file.content_type.startswith("image/"):
                results.append({"filename": file.filename, "error": "Arquivo não é uma imagem"})
                continue

            await file.read()
            result = _build_mock_result()
            result["filename"] = file.filename
            results.append(result)

        return {"results": results, "total": len(results)}

    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Erro ao processar imagens: {str(e)}")


@app.post("/ocr/structured")
async def extract_structured_text(file: UploadFile = File(...)):
    """
    Endpoint simulado com saída estruturada por linhas.
    """
    try:
        _validate_image_upload(file)
        await file.read()

        result = _build_mock_result()
        return {
            "success": True,
            "lines": result["structured"]["lines"],
            "total_lines": result["structured"]["total_lines"],
            "detected_part_numbers": result["detected_part_numbers"],
            "detected_serial_numbers": result["detected_serial_numbers"],
            "processing_time_ms": result["processing_time_ms"],
            "mock": True,
        }

    except HTTPException:
        raise
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Erro ao processar imagem: {str(e)}")

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=5001, log_level="info")
