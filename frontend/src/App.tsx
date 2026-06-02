import { useState } from 'react';
import { Layout, Row, Col, message, Alert } from 'antd';
import { CameraOutlined } from '@ant-design/icons';
import { WebcamCapture } from './components/WebcamCapture';
import { ImagePreview } from './components/ImagePreview';
import { ApprovalForm } from './components/ApprovalForm';
import { HistoryTable } from './components/HistoryTable';
import { ocrService, parserService, inventoryService } from './services/api';
import './App.css';

const { Header, Content, Footer } = Layout;

interface OCRResult {
  success: boolean;
  text: string[];
  confidence: number[];
  mock?: boolean;
  detected_part_numbers?: string[];
  detected_serial_numbers?: string[];
  error?: string;
}

interface ParserResult {
  success: boolean;
  part_number?: string;
  serial_number?: string;
  manufacturer?: string;
  category?: string;
  normalized_description?: string;
  confidence?: number;
  signals?: string[];
  error?: string;
}

function App() {
  const [capturedImage, setCapturedImage] = useState<string | undefined>();
  const [ocrResult, setOcrResult] = useState<OCRResult | undefined>();
  const [loadingOCR, setLoadingOCR] = useState(false);
  const [loadingInventory, setLoadingInventory] = useState(false);
  const [suggestedPN, setSuggestedPN] = useState('');
  const [catalogSearchResult, setCatalogSearchResult] = useState<any>(null);
  const [parserResult, setParserResult] = useState<ParserResult | undefined>();
  const [ocrMockActive, setOcrMockActive] = useState(false);
  const [refreshHistory, setRefreshHistory] = useState(0);

  const handleCapture = async (imageSrc: string, blob: Blob) => {
    setCapturedImage(imageSrc);
    setOcrResult(undefined);
    setSuggestedPN('');
    setCatalogSearchResult(null);
    setParserResult(undefined);
    setOcrMockActive(false);

    // Executar OCR automaticamente
    await performOCR(blob);
  };

  const performOCR = async (blob: Blob) => {
    setLoadingOCR(true);
    try {
      const result = await ocrService.extractText(blob);
      
      if (result.success) {
        setOcrResult(result);

        if (result.mock) {
          setOcrMockActive(true);
          setParserResult(undefined);
          setSuggestedPN('');
          setCatalogSearchResult(null);
          message.warning('OCR em modo mock ativo. Inicie o OCR real para extrair os dados corretos da etiqueta.');
          return;
        }

        try {
          const parsed = await parserService.parseOcrText(result.text || []);
          setParserResult(parsed);

          const suggestedPartNumber = parsed.part_number || result.detected_part_numbers?.[0] || result.text[0] || '';
          setSuggestedPN(suggestedPartNumber);

          if (suggestedPartNumber) {
            const searchResult = await inventoryService.searchCatalog(suggestedPartNumber);
            setCatalogSearchResult(searchResult);

            if (searchResult.found) {
              message.success(`Item encontrado no catálogo: ${searchResult.item.part_number}`);
            } else {
              message.info('Item não encontrado no catálogo. Você pode adicionar manualmente.');
            }
          }
        } catch (error) {
          console.error('Erro ao processar no parser:', error);

          const fallbackPartNumber = result.detected_part_numbers?.[0] || result.text[0] || '';
          setSuggestedPN(fallbackPartNumber);
        }
      } else {
        message.error('Erro ao fazer OCR: ' + (result.error || 'Erro desconhecido'));
      }
    } catch (error) {
      message.error('Erro ao fazer OCR');
      console.error(error);
    } finally {
      setLoadingOCR(false);
    }
  };

  const handleSubmit = async (data: {
    part_number: string;
    serial_number?: string;
    quantity: number;
    location: string;
    reason?: string;
  }) => {
    setLoadingInventory(true);
    try {
      const response = await inventoryService.addInventory({
        part_number: data.part_number,
        serial_number: data.serial_number,
        quantity: data.quantity,
        location: data.location,
        reason: data.reason,
        user_id: 1, // TODO: obter do contexto de autenticação
      });

      if (response.success) {
        message.success('Item adicionado ao estoque com sucesso!');
        
        // Limpar formulário
        setCapturedImage(undefined);
        setOcrResult(undefined);
        setSuggestedPN('');
        setCatalogSearchResult(null);
        
        // Atualizar tabela
        setRefreshHistory(prev => prev + 1);
      } else {
        message.error(response.error || 'Erro ao adicionar item');
      }
    } catch (error) {
      message.error('Erro ao adicionar item ao estoque');
      console.error(error);
    } finally {
      setLoadingInventory(false);
    }
  };

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Header
        style={{
          background: '#001529',
          color: 'white',
          padding: '0 50px',
          display: 'flex',
          alignItems: 'center',
          fontSize: '20px',
          fontWeight: 'bold',
        }}
      >
        <CameraOutlined style={{ marginRight: '10px', fontSize: '24px' }} />
        Sistema de Entrada de Estoque por Visão Computacional
      </Header>

      <Content style={{ padding: '24px' }}>
        <div style={{ maxWidth: '1200px', margin: '0 auto' }}>
          {/* Alerta de informações */}
          <Alert
            message="Fase 1 - MVP"
            description="Este é o sistema inicial com captura de webcam e OCR. Novos recursos serão adicionados nas próximas fases."
            type="info"
            showIcon
            closable
            style={{ marginBottom: '24px' }}
          />

          {ocrMockActive && (
            <Alert
              style={{ marginBottom: '24px' }}
              type="warning"
              showIcon
              message="OCR mock ativo"
              description="O serviço atual retorna texto fixo para testes e não lê a foto real. Inicie services/ocr/main.py para identificação real da memória."
            />
          )}

          {ocrResult && catalogSearchResult && (
            <Alert
              style={{ marginBottom: '24px' }}
              type={catalogSearchResult.found ? 'success' : 'warning'}
              showIcon
              message={catalogSearchResult.found ? 'Part Number localizado' : 'Part Number não localizado'}
              description={catalogSearchResult.found && catalogSearchResult.item ? (
                <span>
                  <strong>{catalogSearchResult.item.part_number}</strong> ·{' '}
                  {catalogSearchResult.item.manufacturer} · {catalogSearchResult.item.category}
                </span>
              ) : (
                'Revise o texto detectado antes de confirmar a entrada.'
              )}
            />
          )}

          {parserResult && parserResult.success && (
            <Alert
              style={{ marginBottom: '24px' }}
              type="info"
              showIcon
              message={`Parser: ${parserResult.category || 'unknown'}${parserResult.confidence ? ` · confiança ${Math.round(parserResult.confidence * 100)}%` : ''}`}
              description={
                parserResult.normalized_description
                  ? `Descrição normalizada: ${parserResult.normalized_description}`
                  : 'O parser não gerou descrição normalizada.'
              }
            />
          )}

          {/* Seção de Captura */}
          <Row gutter={24} style={{ marginBottom: '24px' }}>
            <Col xs={24} md={12}>
              <WebcamCapture onCapture={handleCapture} loading={loadingOCR} />
            </Col>
            <Col xs={24} md={12}>
              <ImagePreview 
                imageSrc={capturedImage} 
                loading={loadingOCR}
                title="Preview"
              />
            </Col>
          </Row>

          {/* Seção de Resultados de OCR */}
          {ocrResult && (
            <Row gutter={24} style={{ marginBottom: '24px' }}>
              <Col xs={24}>
                <ApprovalForm
                  ocrText={ocrResult.text}
                  suggestedPN={suggestedPN}
                  parserResult={parserResult}
                  onSubmit={handleSubmit}
                  loading={loadingInventory}
                />
              </Col>
            </Row>
          )}

          {/* Histórico */}
          <Row gutter={24} style={{ marginBottom: '24px' }}>
            <Col xs={24}>
              <HistoryTable key={refreshHistory} />
            </Col>
          </Row>
        </div>
      </Content>

      <Footer style={{ textAlign: 'center', background: '#f0f2f5' }}>
        <div>
          <strong>Sistema de Entrada de Estoque</strong>
          <br />
          <small>Tecnologias: React, Go, Python, PostgreSQL</small>
        </div>
      </Footer>
    </Layout>
  );
}

export default App;
