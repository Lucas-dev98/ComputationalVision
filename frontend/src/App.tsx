import { useState } from 'react';
import { Layout, Row, Col, message, Alert } from 'antd';
import { CameraOutlined } from '@ant-design/icons';
import { WebcamCapture } from './components/WebcamCapture';
import { ImagePreview } from './components/ImagePreview';
import { ApprovalForm } from './components/ApprovalForm';
import { HistoryTable } from './components/HistoryTable';
import { ocrService, parserService, inventoryService, webResearchService } from './services/api';
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
  tokens?: string[];
  error?: string;
}

interface WebResearchResult {
  success: boolean;
  part_number?: string;
  found?: boolean;
  manufacturer?: string;
  category?: string;
  normalized_description?: string;
  confidence?: number;
  signals?: string[];
  error?: string;
}

const PART_NUMBER_REGEX = /\b[A-Z0-9][A-Z0-9._/-]{5,}\b/g;
const SERIAL_LABELED_REGEX = /\b(?:SN|S\/N|S\.N\.|SER\.?\s*NO\.?|SERIAL|SERIAL\s*NO|SERIALNUMBER)\s*[:#-]?\s*([A-Z0-9][A-Z0-9-]{4,})\b/i;
const SERIAL_CANDIDATE_REGEX = /\b[A-Z0-9][A-Z0-9-]{7,}\b/g;

function isLikelySpecToken(token: string): boolean {
  const upperToken = token.toUpperCase();
  if (["DDR3", "DDR4", "DDR5", "RDIMM", "UDIMM", "DIMM", "ECC", "NVME", "SSD", "HDD"].includes(upperToken)) {
    return true;
  }
  if (/^\d+(GB|TB|MHZ|RPM|MT\/S)$/i.test(upperToken)) {
    return true;
  }
  if (/^PC[345]-?\d+/i.test(upperToken)) {
    return true;
  }
  return false;
}

function candidateScore(token: string): number {
  let score = 0;
  if (/\d/.test(token)) {
    score += 2;
  }
  if (/[A-Z]/.test(token)) {
    score += 2;
  }
  if (/[-_/.]/.test(token)) {
    score += 2;
  }
  if (token.length >= 10) {
    score += 2;
  }
  if (/^SN[-_]?/i.test(token)) {
    score -= 3;
  }
  return score;
}

function extractPartNumberFromOCRText(lines: string[]): string {
  const candidates: string[] = [];

  for (const rawLine of lines || []) {
    const upperLine = (rawLine || "").toUpperCase();
    const matches = upperLine.match(PART_NUMBER_REGEX) || [];

    for (const match of matches) {
      const token = match.trim();
      if (token.length < 6) {
        continue;
      }
      if (!/\d/.test(token)) {
        continue;
      }
      if (isLikelySpecToken(token)) {
        continue;
      }
      candidates.push(token);
    }
  }

  if (candidates.length === 0) {
    return "";
  }

  const sorted = [...new Set(candidates)].sort((a, b) => {
    const scoreDiff = candidateScore(b) - candidateScore(a);
    if (scoreDiff !== 0) {
      return scoreDiff;
    }
    return b.length - a.length;
  });

  return sorted[0] || "";
}

function extractSerialNumberFromOCRText(lines: string[], partNumber = ""): string {
  const normalizedPartNumber = partNumber.toUpperCase().trim();

  const isLikelyWWN = (value: string, line: string) => {
    const normalizedValue = value.toUpperCase().replace(/[^A-Z0-9]/g, '');
    const upperLine = line.toUpperCase();

    if (upperLine.includes('WWN') || upperLine.includes('WORLD WIDE NAME') || upperLine.includes('NAA')) {
      return true;
    }

    // WWN costuma ser hexadecimal longo (ex.: 5001B448B6351C20)
    if (/^[0-9A-F]{16,32}$/.test(normalizedValue)) {
      return true;
    }

    return false;
  };

  for (const rawLine of lines || []) {
    const upperLine = (rawLine || "").toUpperCase();
    const labeled = upperLine.match(SERIAL_LABELED_REGEX);
    if (labeled?.[1]) {
      const serial = labeled[1].trim();
      if (serial !== normalizedPartNumber && !isLikelySpecToken(serial) && !isLikelyWWN(serial, upperLine)) {
        return serial;
      }
    }
  }

  const candidates: string[] = [];
  for (const rawLine of lines || []) {
    const upperLine = (rawLine || "").toUpperCase();
    const matches = upperLine.match(SERIAL_CANDIDATE_REGEX) || [];

    for (const token of matches) {
      const value = token.trim();
      if (value === normalizedPartNumber) {
        continue;
      }
      if (!/[A-Z]/.test(value) || !/\d/.test(value)) {
        continue;
      }
      if (isLikelySpecToken(value)) {
        continue;
      }
      if (isLikelyWWN(value, upperLine)) {
        continue;
      }
      candidates.push(value);
    }
  }

  if (candidates.length === 0) {
    return "";
  }

  const sorted = [...new Set(candidates)].sort((a, b) => {
    const hasSNA = /^SN[-_]?/.test(a) ? 1 : 0;
    const hasSNB = /^SN[-_]?/.test(b) ? 1 : 0;
    if (hasSNA !== hasSNB) {
      return hasSNB - hasSNA;
    }
    return b.length - a.length;
  });

  return sorted[0] || "";
}

function extractLabeledSerialNumberFromOCRText(lines: string[], partNumber = ""): string {
  const normalizedPartNumber = partNumber.toUpperCase().trim();

  const isLikelyWWN = (value: string, line: string) => {
    const normalizedValue = value.toUpperCase().replace(/[^A-Z0-9]/g, '');
    const upperLine = line.toUpperCase();

    if (upperLine.includes('WWN') || upperLine.includes('WORLD WIDE NAME') || upperLine.includes('NAA')) {
      return true;
    }

    if (/^[0-9A-F]{16,32}$/.test(normalizedValue)) {
      return true;
    }

    return false;
  };

  for (const rawLine of lines || []) {
    const upperLine = (rawLine || '').toUpperCase();
    const labeled = upperLine.match(SERIAL_LABELED_REGEX);
    if (!labeled?.[1]) {
      continue;
    }

    const serial = labeled[1].trim();
    if (serial === normalizedPartNumber) {
      continue;
    }
    if (isLikelySpecToken(serial)) {
      continue;
    }
    if (isLikelyWWN(serial, upperLine)) {
      continue;
    }

    return serial;
  }

  return '';
}

function App() {
  const [capturedImage, setCapturedImage] = useState<string | undefined>();
  const [ocrResult, setOcrResult] = useState<OCRResult | undefined>();
  const [loadingOCR, setLoadingOCR] = useState(false);
  const [loadingInventory, setLoadingInventory] = useState(false);
  const [suggestedPN, setSuggestedPN] = useState('');
  const [suggestedSN, setSuggestedSN] = useState('');
  const [catalogSearchResult, setCatalogSearchResult] = useState<any>(null);
  const [parserResult, setParserResult] = useState<ParserResult | undefined>();
  const [webResearchResult, setWebResearchResult] = useState<WebResearchResult | undefined>();
  const [ocrMockActive, setOcrMockActive] = useState(false);
  const [refreshHistory, setRefreshHistory] = useState(0);

  const handleCapture = async (imageSrc: string, blob: Blob) => {
    setCapturedImage(imageSrc);
    setOcrResult(undefined);
    setSuggestedPN('');
    setSuggestedSN('');
    setCatalogSearchResult(null);
    setParserResult(undefined);
    setWebResearchResult(undefined);
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
          const ocrFallbackPartNumber = extractPartNumberFromOCRText(result.text || []);
          const suggestedPartNumber = parsed.part_number || result.detected_part_numbers?.[0] || ocrFallbackPartNumber || result.text[0] || '';
          const labeledSerial = extractLabeledSerialNumberFromOCRText(result.text || [], suggestedPartNumber);
          const ocrFallbackSerial = extractSerialNumberFromOCRText(result.text || [], suggestedPartNumber);
          const suggestedSerialNumber = labeledSerial || parsed.serial_number || result.detected_serial_numbers?.[0] || ocrFallbackSerial || '';

          const parsedWithFallback = {
            ...parsed,
            serial_number: suggestedSerialNumber,
          };

          setParserResult(parsedWithFallback);
          setSuggestedPN(suggestedPartNumber);
          setSuggestedSN(suggestedSerialNumber);

          if (!parsed.part_number && ocrFallbackPartNumber) {
            message.success(`PN identificado via OCR: ${ocrFallbackPartNumber}`);
          }
          if (!parsed.serial_number && suggestedSerialNumber) {
            message.success(`Serial Number identificado via OCR: ${suggestedSerialNumber}`);
          }

          if (suggestedPartNumber) {
            let shouldRunWebResearch = true;

            try {
              const searchResult = await inventoryService.searchCatalog(suggestedPartNumber);
              setCatalogSearchResult(searchResult);

              if (searchResult.found) {
                shouldRunWebResearch = false;
                message.success(`Item encontrado no catálogo: ${searchResult.item.part_number}`);
                setWebResearchResult(undefined);
              } else {
                message.info('Item não encontrado no catálogo. Você pode adicionar manualmente.');
              }
            } catch (catalogError) {
              setCatalogSearchResult({ found: false, error: 'Inventory API indisponível no momento.' });
              message.warning('Inventory API indisponível. Continuando com pesquisa web automática.');
              console.warn('Falha ao consultar catálogo:', catalogError);
            }

            if (shouldRunWebResearch) {
              try {
                const research = await webResearchService.researchPartNumber({
                  part_number: suggestedPartNumber,
                  manufacturer: parsedWithFallback.manufacturer,
                  category: parsedWithFallback.category,
                  normalized_description: parsedWithFallback.normalized_description,
                  tokens: parsedWithFallback.tokens || result.text || [],
                });

                setWebResearchResult(research);

                if (research.success && research.confidence && research.confidence >= 0.5) {
                  setParserResult({
                    ...parsedWithFallback,
                    manufacturer: parsedWithFallback.manufacturer || research.manufacturer,
                    category: parsedWithFallback.category && parsedWithFallback.category !== 'unknown' ? parsedWithFallback.category : research.category,
                    normalized_description: parsedWithFallback.normalized_description || research.normalized_description,
                  });
                  message.success('Pesquisa web automática sugeriu enriquecimento para o item.');
                }
              } catch (researchError) {
                console.warn('Pesquisa web indisponível neste momento:', researchError);
              }
            }
          }
        } catch (error) {
          console.error('Erro ao processar no parser:', error);

          const ocrFallbackPartNumber = extractPartNumberFromOCRText(result.text || []);
          const fallbackPartNumber = result.detected_part_numbers?.[0] || ocrFallbackPartNumber || result.text[0] || '';
          const labeledSerial = extractLabeledSerialNumberFromOCRText(result.text || [], fallbackPartNumber);
          const fallbackSerialNumber = labeledSerial || result.detected_serial_numbers?.[0] || extractSerialNumberFromOCRText(result.text || [], fallbackPartNumber) || '';
          setSuggestedPN(fallbackPartNumber);
          setSuggestedSN(fallbackSerialNumber);

          if (fallbackPartNumber) {
            try {
              const research = await webResearchService.researchPartNumber({
                part_number: fallbackPartNumber,
                tokens: result.text || [],
              });
              setWebResearchResult(research);
            } catch (researchError) {
              console.warn('Pesquisa web indisponível durante fallback:', researchError);
            }
          }
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
        const correctionApplied =
          (suggestedPN || '') !== (data.part_number || '') ||
          (suggestedSN || '') !== (data.serial_number || '');

        try {
          await inventoryService.submitFeedback({
            part_number_predicted: suggestedPN || parserResult?.part_number || '',
            part_number_final: data.part_number || '',
            serial_number_predicted: suggestedSN || parserResult?.serial_number || '',
            serial_number_final: data.serial_number || '',
            manufacturer_predicted: parserResult?.manufacturer || '',
            manufacturer_final: parserResult?.manufacturer || '',
            category_predicted: parserResult?.category || '',
            category_final: parserResult?.category || '',
            correction_applied: correctionApplied,
            confidence: parserResult?.confidence || 0,
            image_data: capturedImage || '',
            ocr_text: ocrResult?.text || [],
            meta_json: JSON.stringify({
              web_research: webResearchResult,
              catalog_search: catalogSearchResult,
            }),
          });
        } catch (feedbackError) {
          console.warn('Não foi possível registrar feedback de aprendizado:', feedbackError);
        }

        message.success('Item adicionado ao estoque com sucesso!');
        
        // Limpar formulário
        setCapturedImage(undefined);
        setOcrResult(undefined);
        setSuggestedPN('');
        setSuggestedSN('');
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

          {webResearchResult && webResearchResult.success && (
            <Alert
              style={{ marginBottom: '24px' }}
              type={webResearchResult.found ? 'success' : 'info'}
              showIcon
              message={webResearchResult.found ? 'Pesquisa web automática concluída' : 'Pesquisa web executada sem resultados fortes'}
              description={
                webResearchResult.manufacturer || webResearchResult.category || webResearchResult.normalized_description
                  ? `Sugestão: ${webResearchResult.manufacturer || 'Fabricante N/D'} · ${webResearchResult.category || 'Categoria N/D'} · ${webResearchResult.normalized_description || 'Sem descrição normalizada'}`
                  : 'Nenhum enriquecimento adicional encontrado para este PN.'
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
                  suggestedSN={suggestedSN}
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
