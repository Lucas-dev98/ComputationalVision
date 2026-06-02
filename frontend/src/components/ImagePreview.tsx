import React from 'react';
import { Card, Spin, Empty } from 'antd';
import { LoadingOutlined } from '@ant-design/icons';

interface ImagePreviewProps {
  imageSrc?: string;
  loading?: boolean;
  title?: string;
}

export const ImagePreview: React.FC<ImagePreviewProps> = ({ 
  imageSrc, 
  loading = false, 
  title = 'Preview da Imagem' 
}) => {
  if (loading) {
    return (
      <Card title={title} style={{ marginBottom: '20px' }}>
        <div style={{ textAlign: 'center', padding: '40px 0' }}>
          <Spin 
            indicator={<LoadingOutlined style={{ fontSize: 48 }} spin />} 
            tip="Processando imagem..."
          />
        </div>
      </Card>
    );
  }

  if (!imageSrc) {
    return (
      <Card title={title} style={{ marginBottom: '20px' }}>
        <Empty description="Nenhuma imagem capturada" />
      </Card>
    );
  }

  return (
    <Card title={title} style={{ marginBottom: '20px' }}>
      <div style={{ textAlign: 'center' }}>
        <img
          src={imageSrc}
          alt="Preview"
          style={{
            maxWidth: '100%',
            maxHeight: '400px',
            borderRadius: '8px',
            border: '1px solid #d9d9d9',
          }}
        />
      </div>
    </Card>
  );
};

export default ImagePreview;
