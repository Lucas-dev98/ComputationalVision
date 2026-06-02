import io
import time
import logging
import numpy as np
import cv2
from PIL import Image
from paddleocr import PaddleOCR
from typing import List, Dict, Any

logger = logging.getLogger(__name__)

class OCRService:
    """Serviço de OCR usando PaddleOCR"""
    
    def __init__(self, lang: List[str] = ['pt', 'en']):
        """
        Inicializar serviço OCR.
        
        Args:
            lang: Lista de idiomas para OCR (português e inglês por padrão)
        """
        logger.info(f"Inicializando PaddleOCR com idiomas: {lang}")
        self.ocr = PaddleOCR(use_angle_cls=True, lang=lang)
        logger.info("PaddleOCR inicializado com sucesso")
    
    def extract_text(self, image_bytes: bytes) -> Dict[str, Any]:
        """
        Extrair texto de uma imagem.
        
        Args:
            image_bytes: Bytes da imagem
            
        Returns:
            Dict com texto, confiança e caixas de detecção
        """
        start_time = time.time()
        
        try:
            # Converter bytes para imagem numpy
            image_array = np.frombuffer(image_bytes, dtype=np.uint8)
            image = cv2.imdecode(image_array, cv2.IMREAD_COLOR)
            
            if image is None:
                raise ValueError("Falha ao decodificar imagem")
            
            # Executar OCR
            logger.info(f"Executando OCR em imagem de tamanho {image.shape}")
            results = self.ocr.ocr(image, cls=True)
            
            # Processar resultados
            text_lines = []
            confidence_scores = []
            boxes = []
            
            if results and len(results) > 0:
                for line in results:
                    if line is None:
                        continue
                    
                    for word_info in line:
                        # word_info é [(pontos), (texto, confiança)]
                        if len(word_info) >= 2:
                            pts = word_info[0]  # Coordenadas
                            text, confidence = word_info[1]  # Texto e confiança
                            
                            text_lines.append(text)
                            confidence_scores.append(float(confidence))
                            
                            # Converter coordenadas para bbox simples
                            pts_array = np.array(pts)
                            x_coords = pts_array[:, 0]
                            y_coords = pts_array[:, 1]
                            
                            bbox = [
                                float(np.min(x_coords)),
                                float(np.min(y_coords)),
                                float(np.max(x_coords)),
                                float(np.max(y_coords))
                            ]
                            boxes.append(bbox)
            
            processing_time = (time.time() - start_time) * 1000  # em ms
            
            logger.info(f"OCR concluído em {processing_time:.2f}ms. "
                       f"Extraído {len(text_lines)} elementos")
            
            return {
                "success": True,
                "text": text_lines,
                "confidence": confidence_scores,
                "boxes": boxes,
                "processing_time_ms": processing_time
            }
            
        except Exception as e:
            logger.error(f"Erro ao processar OCR: {str(e)}", exc_info=True)
            processing_time = (time.time() - start_time) * 1000
            
            return {
                "success": False,
                "error": str(e),
                "text": [],
                "confidence": [],
                "boxes": [],
                "processing_time_ms": processing_time
            }
    
    def extract_structured_text(self, image_bytes: bytes) -> Dict[str, Any]:
        """
        Extrair texto estruturado (linhas agrupadas).
        
        Args:
            image_bytes: Bytes da imagem
            
        Returns:
            Dict com estrutura de linhas e palavras
        """
        start_time = time.time()
        
        try:
            image_array = np.frombuffer(image_bytes, dtype=np.uint8)
            image = cv2.imdecode(image_array, cv2.IMREAD_COLOR)
            
            if image is None:
                raise ValueError("Falha ao decodificar imagem")
            
            results = self.ocr.ocr(image, cls=True)
            
            lines = []
            if results and len(results) > 0:
                for line in results:
                    if line is None:
                        continue
                    
                    line_text = " ".join([word_info[1][0] for word_info in line])
                    line_conf = np.mean([word_info[1][1] for word_info in line])
                    
                    lines.append({
                        "text": line_text,
                        "confidence": float(line_conf),
                        "words": [
                            {
                                "text": word_info[1][0],
                                "confidence": float(word_info[1][1]),
                                "bbox": [
                                    [float(pt[0]), float(pt[1])]
                                    for pt in word_info[0]
                                ]
                            }
                            for word_info in line
                        ]
                    })
            
            processing_time = (time.time() - start_time) * 1000
            
            return {
                "success": True,
                "lines": lines,
                "total_lines": len(lines),
                "processing_time_ms": processing_time
            }
            
        except Exception as e:
            logger.error(f"Erro ao processar OCR estruturado: {str(e)}", exc_info=True)
            processing_time = (time.time() - start_time) * 1000
            
            return {
                "success": False,
                "error": str(e),
                "lines": [],
                "processing_time_ms": processing_time
            }
