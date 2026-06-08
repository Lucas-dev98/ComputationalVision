import axios from 'axios';

// Detectar porta da API (desenvolvimento ou produção)
const API_URL = (() => {
  if (process.env.REACT_APP_API_URL) {
    return process.env.REACT_APP_API_URL;
  }
  // Fallback para ambiente de desenvolvimento
  return window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1'
    ? 'http://localhost:8081'
    : 'http://localhost:8080';
})();

const OCR_URL = process.env.REACT_APP_OCR_URL || 'http://localhost:5001';
const PARSER_URL = process.env.REACT_APP_PARSER_URL || 'http://localhost:8082';
const WEB_RESEARCH_URL = process.env.REACT_APP_WEB_RESEARCH_URL || 'http://localhost:8083';
const OCR_TIMEOUT_MS = Number(process.env.REACT_APP_OCR_TIMEOUT_MS || '90000');

const apiClient = axios.create({
  baseURL: API_URL,
  timeout: 10000,
});

const ocrClient = axios.create({
  baseURL: OCR_URL,
  timeout: OCR_TIMEOUT_MS,
});

const parserClient = axios.create({
  baseURL: PARSER_URL,
  timeout: 15000,
});

const webResearchClient = axios.create({
  baseURL: WEB_RESEARCH_URL,
  timeout: 12000,
});

// Serviço de OCR
export const ocrService = {
  async extractText(imageFile: File | Blob) {
    const formData = new FormData();
    formData.append('file', imageFile);
    
    try {
      const response = await ocrClient.post('/ocr', formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      });
      return response.data;
    } catch (error) {
      const axiosError = error as { code?: string };

      // O OCR real pode levar mais tempo no primeiro request (cold start dos modelos).
      if (axiosError?.code === 'ECONNABORTED') {
        console.warn('Timeout no OCR, tentando novamente com timeout estendido...');
        const retryResponse = await ocrClient.post('/ocr', formData, {
          headers: {
            'Content-Type': 'multipart/form-data',
          },
          timeout: Math.max(OCR_TIMEOUT_MS, 120000),
        });
        return retryResponse.data;
      }

      console.error('Erro ao fazer OCR:', error);
      throw error;
    }
  },
};

// Serviço de Inventário
export const inventoryService = {
  async searchCatalog(partNumber: string) {
    try {
      const response = await apiClient.get(`/catalog/search?pn=${partNumber}`);
      return response.data;
    } catch (error) {
      if (axios.isAxiosError(error) && error.response?.status === 404) {
        return { found: false };
      }
      console.error('Erro ao buscar no catálogo:', error);
      throw error;
    }
  },

  async addInventory(data: {
    part_number: string;
    serial_number?: string;
    quantity?: number;
    location?: string;
    reason?: string;
    user_id?: number;
  }) {
    try {
      const response = await apiClient.post('/inventory/in', data);
      return response.data;
    } catch (error) {
      console.error('Erro ao adicionar ao estoque:', error);
      throw error;
    }
  },

  async listInventory(limit = 50, offset = 0) {
    try {
      const response = await apiClient.get(`/inventory/items?limit=${limit}&offset=${offset}`);
      return response.data;
    } catch (error) {
      console.error('Erro ao listar estoque:', error);
      throw error;
    }
  },

  async getInventoryItem(id: number) {
    try {
      const response = await apiClient.get(`/inventory/items/${id}`);
      return response.data;
    } catch (error) {
      console.error('Erro ao obter item:', error);
      throw error;
    }
  },

  async submitFeedback(data: {
    part_number_predicted?: string;
    part_number_final?: string;
    serial_number_predicted?: string;
    serial_number_final?: string;
    manufacturer_predicted?: string;
    manufacturer_final?: string;
    category_predicted?: string;
    category_final?: string;
    correction_applied: boolean;
    confidence?: number;
    image_data?: string;
    ocr_text?: string[];
    meta_json?: string;
  }) {
    try {
      const response = await apiClient.post('/feedback/submit', data);
      return response.data;
    } catch (error) {
      console.error('Erro ao enviar feedback de aprendizado:', error);
      throw error;
    }
  },

  async listActiveLearningFeedback(limit = 200, correctionsOnly = true) {
    try {
      const response = await apiClient.get(`/feedback/active-learning?limit=${limit}&corrections_only=${correctionsOnly}`);
      return response.data;
    } catch (error) {
      console.error('Erro ao listar feedbacks de aprendizado:', error);
      throw error;
    }
  },
};

export const parserService = {
  async parseOcrText(text: string[]) {
    try {
      const response = await parserClient.post('/parse', { text });
      return response.data;
    } catch (error) {
      console.error('Erro ao processar texto no parser:', error);
      throw error;
    }
  },
};

export const webResearchService = {
  async researchPartNumber(data: {
    part_number: string;
    manufacturer?: string;
    category?: string;
    normalized_description?: string;
    tokens?: string[];
  }) {
    try {
      const response = await webResearchClient.post('/research', data);
      return response.data;
    } catch (error) {
      console.error('Erro na pesquisa web automática:', error);
      throw error;
    }
  },
};

const services = {
  apiClient,
  ocrClient,
  parserClient,
  webResearchClient,
  ocrService,
  parserService,
  webResearchService,
  inventoryService,
};

export default services;
