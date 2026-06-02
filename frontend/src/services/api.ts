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

const apiClient = axios.create({
  baseURL: API_URL,
  timeout: 10000,
});

const ocrClient = axios.create({
  baseURL: OCR_URL,
  timeout: 30000,
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
};

export default {
  apiClient,
  ocrClient,
  ocrService,
  inventoryService,
};
