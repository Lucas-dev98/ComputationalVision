import os
import io
import time
import logging
import numpy as np
import cv2
from PIL import Image

# Workaround para incompatibilidades do runtime Paddle em CPU/Windows.
os.environ.setdefault("FLAGS_use_mkldnn", "0")
os.environ.setdefault("FLAGS_enable_pir_api", "0")

from paddleocr import PaddleOCR
from typing import List, Dict, Any

logger = logging.getLogger(__name__)

class OCRService:
    """Serviço de OCR usando PaddleOCR"""

    @staticmethod
    def _bbox_from_points(points: Any) -> List[float]:
        pts_array = np.array(points)
        if pts_array.size == 0:
            return [0.0, 0.0, 0.0, 0.0]

        if pts_array.ndim == 1:
            return [float(pts_array[0]), float(pts_array[1]), float(pts_array[0]), float(pts_array[1])]

        x_coords = pts_array[:, 0]
        y_coords = pts_array[:, 1]
        return [
            float(np.min(x_coords)),
            float(np.min(y_coords)),
            float(np.max(x_coords)),
            float(np.max(y_coords)),
        ]

    def _extract_flat_items(self, results: Any) -> List[Dict[str, Any]]:
        items: List[Dict[str, Any]] = []
        if not results:
            return items

        # Formato legado do PaddleOCR: [[ [poly, (text, score)], ... ]]
        if isinstance(results, list) and len(results) > 0 and isinstance(results[0], list):
            for line in results:
                if not line:
                    continue
                for word_info in line:
                    if not isinstance(word_info, (list, tuple)) or len(word_info) < 2:
                        continue
                    poly = word_info[0]
                    text_info = word_info[1]
                    if not isinstance(text_info, (list, tuple)) or len(text_info) < 2:
                        continue
                    items.append(
                        {
                            "text": str(text_info[0]),
                            "confidence": float(text_info[1]),
                            "poly": poly,
                        }
                    )
            return items

        # Formato novo (pipeline): lista de dicts com rec_texts/rec_scores/rec_polys
        if isinstance(results, list):
            for result in results:
                if not isinstance(result, dict):
                    continue
                texts = result.get("rec_texts") or []
                scores = result.get("rec_scores") or []
                polys = result.get("rec_polys") or result.get("dt_polys") or []

                for i, text in enumerate(texts):
                    score = scores[i] if i < len(scores) else 0.0
                    poly = polys[i] if i < len(polys) else []
                    items.append(
                        {
                            "text": str(text),
                            "confidence": float(score),
                            "poly": poly,
                        }
                    )
        return items
    
    def __init__(self, lang: List[str] | str = 'pt'):
        """
        Inicializar serviço OCR.
        
        Args:
            lang: Idioma preferencial para OCR
        """
        if isinstance(lang, list):
            candidates = lang + ['en']
        else:
            candidates = [lang, 'en']

        # Remover duplicatas mantendo ordem
        unique_candidates: List[str] = []
        for candidate in candidates:
            if candidate not in unique_candidates:
                unique_candidates.append(candidate)

        last_error: Exception | None = None
        for candidate in unique_candidates:
            try:
                logger.info(f"Inicializando PaddleOCR com idioma: {candidate}")
                try:
                    # API mais nova (v3+): evita pipeline de documentos que falha em alguns ambientes Windows/CPU.
                    self.ocr = PaddleOCR(
                        lang=candidate,
                        use_doc_orientation_classify=False,
                        use_doc_unwarping=False,
                        use_textline_orientation=False,
                        enable_mkldnn=False,
                        cpu_threads=1,
                    )
                except TypeError:
                    # API legada (v2): mantém compatibilidade.
                    self.ocr = PaddleOCR(
                        use_angle_cls=False,
                        lang=candidate,
                        enable_mkldnn=False,
                        cpu_threads=1,
                    )
                logger.info(f"PaddleOCR inicializado com sucesso (lang={candidate})")
                return
            except Exception as error:
                logger.warning(f"Falha ao inicializar PaddleOCR com idioma {candidate}: {error}")
                last_error = error

        raise RuntimeError(f"Não foi possível inicializar o PaddleOCR: {last_error}")
    
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
            results = self.ocr.ocr(image)
            
            # Processar resultados
            text_lines = []
            confidence_scores = []
            boxes = []
            
            items = self._extract_flat_items(results)
            for item in items:
                text_lines.append(item["text"])
                confidence_scores.append(item["confidence"])
                boxes.append(self._bbox_from_points(item["poly"]))
            
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
            
            results = self.ocr.ocr(image)
            
            lines = []
            items = self._extract_flat_items(results)
            for item in items:
                poly = np.array(item["poly"])
                bbox_points = []
                if poly.size > 0 and poly.ndim == 2:
                    bbox_points = [[float(pt[0]), float(pt[1])] for pt in poly]

                lines.append(
                    {
                        "text": item["text"],
                        "confidence": item["confidence"],
                        "words": [
                            {
                                "text": item["text"],
                                "confidence": item["confidence"],
                                "bbox": bbox_points,
                            }
                        ],
                    }
                )
            
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
