import os
import logging
from fastapi import FastAPI, File, UploadFile, HTTPException
from fastapi.responses import JSONResponse
from fastapi.middleware.cors import CORSMiddleware
from dotenv import load_dotenv
from ocr_service import OCRService

# Carregar variáveis de ambiente
load_dotenv()

# Configurar logging
logging.basicConfig(
    level=os.getenv('LOG_LEVEL', 'INFO'),
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

# Inicializar FastAPI
app = FastAPI(
    title="OCR Service",
    description="Serviço de OCR para extração de texto de componentes de TI",
    version="1.0.0"
)

# Adicionar CORS
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Inicializar serviço OCR na startup
ocr_service = None

@app.on_event("startup")
async def startup_event():
    global ocr_service
    logger.info("Inicializando OCR Service...")
    ocr_service = OCRService()
    logger.info("OCR Service inicializado com sucesso")

@app.get("/health")
async def health():
    return {
        "status": "healthy",
        "service": "ocr-service",
        "version": "1.0.0"
    }

@app.post("/ocr")
async def extract_text(file: UploadFile = File(...)):
    """
    Endpoint para extrair texto de uma imagem.
    
    Retorna:
    {
        "success": bool,
        "text": [list of strings],
        "confidence": [list of floats],
        "boxes": [[x1, y1, x2, y2, ...], ...],
        "processing_time_ms": float
    }
    """
    try:
        logger.info(f"Processando arquivo: {file.filename}")
        
        # Validar tipo de arquivo
        if not file.content_type.startswith('image/'):
            raise HTTPException(status_code=400, detail="Arquivo não é uma imagem")
        
        # Ler conteúdo do arquivo
        contents = await file.read()
        
        # Processar com OCR
        result = ocr_service.extract_text(contents)
        
        logger.info(f"OCR concluído: {len(result['text'])} linhas extraídas")
        
        return JSONResponse(
            status_code=200,
            content=result
        )
        
    except HTTPException:
        raise
    except Exception as e:
        logger.error(f"Erro ao processar OCR: {str(e)}", exc_info=True)
        raise HTTPException(
            status_code=500,
            detail=f"Erro ao processar imagem: {str(e)}"
        )

@app.post("/ocr/batch")
async def extract_text_batch(files: list[UploadFile] = File(...)):
    """
    Endpoint para extrair texto de múltiplas imagens.
    """
    try:
        results = []
        for file in files:
            if not file.content_type.startswith('image/'):
                results.append({
                    "filename": file.filename,
                    "error": "Arquivo não é uma imagem"
                })
                continue
            
            contents = await file.read()
            result = ocr_service.extract_text(contents)
            result["filename"] = file.filename
            results.append(result)
        
        return JSONResponse(
            status_code=200,
            content={"results": results, "total": len(results)}
        )
        
    except Exception as e:
        logger.error(f"Erro ao processar batch OCR: {str(e)}", exc_info=True)
        raise HTTPException(
            status_code=500,
            detail=f"Erro ao processar imagens: {str(e)}"
        )

if __name__ == "__main__":
    import uvicorn
    
    host = os.getenv('HOST', '0.0.0.0')
    port = int(os.getenv('PORT', 5001))
    
    logger.info(f"Iniciando servidor em {host}:{port}")
    
    uvicorn.run(
        "main:app",
        host=host,
        port=port,
        reload=os.getenv('ENV', 'production') == 'development',
        log_level=os.getenv('LOG_LEVEL', 'info').lower()
    )
