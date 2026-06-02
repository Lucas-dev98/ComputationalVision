#!/usr/bin/env python3
"""
OCR Service Mock - Para testes locais
Retorna dados simulados de OCR
"""
from fastapi import FastAPI, File, UploadFile, HTTPException
from fastapi.middleware.cors import CORSMiddleware
import json
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
        # Validar tipo de arquivo
        if not file.content_type.startswith('image/'):
            raise HTTPException(status_code=400, detail="Arquivo não é uma imagem")
        
        # Ler conteúdo do arquivo
        contents = await file.read()
        
        # Simular processamento
        start_time = time.time()
        
        # Retornar dados simulados
        result = {
            "success": True,
            "text": ["Sistema de Entrada de Estoque", "Visão Computacional"],
            "confidence": [0.95, 0.92],
            "boxes": [
                [[10, 10], [200, 10], [200, 50], [10, 50]],
                [[10, 60], [300, 60], [300, 100], [10, 100]]
            ],
            "processing_time_ms": (time.time() - start_time) * 1000,
            "mock": True
        }
        
        return result
        
    except HTTPException:
        raise
    except Exception as e:
        raise HTTPException(
            status_code=500,
            detail=f"Erro ao processar imagem: {str(e)}"
        )

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=5001, log_level="info")
