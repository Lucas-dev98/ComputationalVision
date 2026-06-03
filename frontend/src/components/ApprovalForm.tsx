import React, { useState, useEffect } from 'react';
import { 
  Form, 
  Input, 
  Button, 
  Card, 
  message, 
  Space, 
  Divider,
  Select,
  Descriptions,
  Tag
} from 'antd';
import { CheckCircleOutlined, FormOutlined } from '@ant-design/icons';

interface ParserResult {
  part_number?: string;
  serial_number?: string;
  manufacturer?: string;
  category?: string;
  normalized_description?: string;
}

interface ApprovalFormProps {
  ocrText?: string[];
  suggestedPN?: string;
  suggestedSN?: string;
  parserResult?: ParserResult;
  onSubmit: (data: {
    part_number: string;
    serial_number?: string;
    quantity: number;
    location: string;
    reason?: string;
  }) => Promise<void>;
  loading?: boolean;
}

export const ApprovalForm: React.FC<ApprovalFormProps> = ({
  ocrText = [],
  suggestedPN = '',
  suggestedSN = '',
  parserResult,
  onSubmit,
  loading = false,
}) => {
  const [form] = Form.useForm();
  const [submitting, setSubmitting] = useState(false);

  useEffect(() => {
    form.setFieldsValue({
      part_number: suggestedPN || parserResult?.part_number || '',
      serial_number: parserResult?.serial_number || suggestedSN || '',
      manufacturer: parserResult?.manufacturer || '',
      category: parserResult?.category || '',
      normalized_description: parserResult?.normalized_description || '',
    });
  }, [suggestedPN, suggestedSN, parserResult, form]);

  const handleSubmit = async (values: any) => {
    setSubmitting(true);
    try {
      await onSubmit({
        part_number: values.part_number,
        serial_number: values.serial_number,
        quantity: values.quantity || 1,
        location: values.location,
        reason: values.reason,
      });
      form.resetFields();
    } catch (error) {
      message.error('Erro ao adicionar item ao estoque');
      console.error(error);
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <Card 
      title="Aprovação de Item"
      extra={<FormOutlined />}
      style={{ marginBottom: '20px' }}
    >
      {parserResult && (
        <Card
          type="inner"
          title="Resumo da extração"
          size="small"
          style={{ marginBottom: '20px' }}
          styles={{ body: { padding: '12px' } }}
        >
          <Descriptions size="small" column={1} bordered>
            <Descriptions.Item label="Part Number">
              {parserResult.part_number || 'Não identificado'}
            </Descriptions.Item>
            <Descriptions.Item label="Serial Number">
              {parserResult.serial_number || 'Não identificado'}
            </Descriptions.Item>
            <Descriptions.Item label="Fabricante">
              {parserResult.manufacturer || 'Não identificado'}
            </Descriptions.Item>
            <Descriptions.Item label="Categoria">
              {parserResult.category ? (
                <Tag color={parserResult.category === 'memory' ? 'blue' : 'geekblue'}>
                  {parserResult.category}
                </Tag>
              ) : (
                'Não identificado'
              )}
            </Descriptions.Item>
            <Descriptions.Item label="Descrição normalizada">
              {parserResult.normalized_description || 'Não gerada'}
            </Descriptions.Item>
          </Descriptions>
        </Card>
      )}

      {/* Exibir texto extraído */}
      {ocrText.length > 0 && (
        <>
          <Card 
            type="inner" 
            title="Texto Extraído (OCR)"
            size="small"
            style={{ marginBottom: '20px' }}
            styles={{ body: { padding: '12px' } }}
          >
            <div style={{ 
              backgroundColor: '#f5f5f5', 
              padding: '12px', 
              borderRadius: '4px',
              maxHeight: '150px',
              overflowY: 'auto',
              fontFamily: 'monospace',
              fontSize: '12px',
            }}>
              {ocrText.map((text, idx) => (
                <div key={idx}>{text}</div>
              ))}
            </div>
          </Card>
          <Divider />
        </>
      )}

      {/* Formulário */}
      <Form
        form={form}
        layout="vertical"
        onFinish={handleSubmit}
      >
        <Form.Item
          label="Part Number (Obrigatório)"
          name="part_number"
          rules={[
            { required: true, message: 'Part Number é obrigatório' },
            { min: 3, message: 'Part Number deve ter pelo menos 3 caracteres' },
          ]}
        >
          <Input 
            placeholder="Ex: M393A4K40DB3-CWE"
            size="large"
          />
        </Form.Item>

        <Form.Item
          label="Serial Number"
          name="serial_number"
        >
          <Input 
            placeholder="Número de série (opcional)"
            size="large"
          />
        </Form.Item>

        <Form.Item
          label="Fabricante"
          name="manufacturer"
        >
          <Input
            placeholder="Fabricante detectado"
            size="large"
            readOnly
          />
        </Form.Item>

        <Form.Item
          label="Categoria"
          name="category"
        >
          <Input
            placeholder="Categoria detectada"
            size="large"
            readOnly
          />
        </Form.Item>

        <Form.Item
          label="Descrição Normalizada"
          name="normalized_description"
        >
          <Input.TextArea
            placeholder="Descrição normalizada pelo parser"
            rows={2}
            readOnly
          />
        </Form.Item>

        <Form.Item
          label="Quantidade"
          name="quantity"
          initialValue={1}
          rules={[
            { required: true, message: 'Quantidade é obrigatória' },
            { pattern: /^[1-9]\d*$/, message: 'Quantidade deve ser um número positivo' },
          ]}
        >
          <Input 
            type="number" 
            min={1}
            placeholder="1"
            size="large"
          />
        </Form.Item>

        <Form.Item
          label="Localização"
          name="location"
          rules={[
            { required: true, message: 'Localização é obrigatória' },
          ]}
        >
          <Select
            placeholder="Selecione a localização"
            size="large"
            options={[
              { label: 'Datacenter RJ', value: 'DC-RJ' },
              { label: 'Armazém SP', value: 'WH-SP' },
              { label: 'Almoxarifado MG', value: 'ST-MG' },
              { label: 'Em trânsito', value: 'TRANSIT' },
            ]}
          />
        </Form.Item>

        <Form.Item
          label="Motivo"
          name="reason"
        >
          <Input.TextArea 
            placeholder="Motivo da entrada (opcional)"
            rows={2}
          />
        </Form.Item>

        <Space style={{ width: '100%', justifyContent: 'center' }}>
          <Button
            type="primary"
            size="large"
            htmlType="submit"
            loading={submitting || loading}
            icon={<CheckCircleOutlined />}
            style={{ minWidth: '200px' }}
          >
            Confirmar Entrada
          </Button>
        </Space>
      </Form>
    </Card>
  );
};

export default ApprovalForm;
