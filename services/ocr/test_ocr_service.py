import unittest
import io
from unittest.mock import patch, MagicMock
from services.ocr.ocr_service import OCRService


class TestOCRService(unittest.TestCase):
    """Testes para OCR Service"""

    @classmethod
    def setUpClass(cls):
        """Inicializar OCR Service uma vez para todos os testes"""
        print("Inicializando OCRService para testes...")
        cls.ocr_service = OCRService()

    def test_ocr_service_init(self):
        """Teste de inicialização do OCR Service"""
        self.assertIsNotNone(self.ocr_service.ocr)

    def test_extract_text_invalid_image(self):
        """Teste com imagem inválida"""
        invalid_image = b"not an image"
        result = self.ocr_service.extract_text(invalid_image)
        
        self.assertFalse(result["success"])
        self.assertIn("error", result)
        self.assertEqual(len(result["text"]), 0)

    def test_extract_text_structure(self):
        """Teste que o resultado tem a estrutura esperada"""
        # Criar uma imagem fake (branca)
        from PIL import Image
        import io
        
        img = Image.new('RGB', (100, 100), color='white')
        img_bytes = io.BytesIO()
        img.save(img_bytes, format='PNG')
        img_bytes.seek(0)
        
        result = self.ocr_service.extract_text(img_bytes.read())
        
        # Verificar estrutura do resultado
        self.assertIn("success", result)
        self.assertIn("text", result)
        self.assertIn("confidence", result)
        self.assertIn("boxes", result)
        self.assertIn("processing_time_ms", result)
        
        # Verificar tipos
        self.assertIsInstance(result["text"], list)
        self.assertIsInstance(result["confidence"], list)
        self.assertIsInstance(result["boxes"], list)
        self.assertIsInstance(result["processing_time_ms"], float)

    def test_extract_structured_text(self):
        """Teste de extração de texto estruturado"""
        from PIL import Image
        import io
        
        img = Image.new('RGB', (100, 100), color='white')
        img_bytes = io.BytesIO()
        img.save(img_bytes, format='PNG')
        img_bytes.seek(0)
        
        result = self.ocr_service.extract_structured_text(img_bytes.read())
        
        self.assertIn("success", result)
        self.assertIn("lines", result)
        self.assertIn("total_lines", result)
        self.assertIn("processing_time_ms", result)

    def test_processing_time_positive(self):
        """Teste que tempo de processamento é positivo"""
        from PIL import Image
        import io
        
        img = Image.new('RGB', (100, 100), color='white')
        img_bytes = io.BytesIO()
        img.save(img_bytes, format='PNG')
        img_bytes.seek(0)
        
        result = self.ocr_service.extract_text(img_bytes.read())
        self.assertGreater(result["processing_time_ms"], 0)


if __name__ == '__main__':
    unittest.main(verbosity=2)
