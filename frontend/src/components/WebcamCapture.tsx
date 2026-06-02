import React, { useRef, useState } from 'react';
import Webcam from 'react-webcam';
import { Button, Space, Card } from 'antd';
import { CameraOutlined, DownloadOutlined } from '@ant-design/icons';

interface WebcamCaptureProps {
  onCapture: (imageSrc: string, blob: Blob) => void;
  loading?: boolean;
}

export const WebcamCapture: React.FC<WebcamCaptureProps> = ({ onCapture, loading = false }) => {
  const webcamRef = useRef<Webcam>(null);
  const [isCameraReady, setIsCameraReady] = useState(false);

  const handleCapture = () => {
    if (webcamRef.current) {
      const imageSrc = webcamRef.current.getScreenshot();
      if (imageSrc) {
        // Converter base64 para Blob
        fetch(imageSrc)
          .then(res => res.blob())
          .then(blob => {
            onCapture(imageSrc, blob);
          });
      }
    }
  };

  return (
    <Card 
      title="Câmera"
      style={{ marginBottom: '20px' }}
      extra={<CameraOutlined />}
    >
      <div style={{ textAlign: 'center', marginBottom: '20px' }}>
        <Webcam
          ref={webcamRef}
          screenshotFormat="image/jpeg"
          width={500}
          height={400}
          onUserMediaError={(error) => console.error('Erro de câmera:', error)}
          onUserMedia={() => setIsCameraReady(true)}
          style={{
            borderRadius: '8px',
            border: '2px solid #1890ff',
            maxWidth: '100%',
            height: 'auto',
          }}
        />
      </div>
      <Space style={{ width: '100%', justifyContent: 'center' }}>
        <Button
          type="primary"
          size="large"
          onClick={handleCapture}
          disabled={!isCameraReady || loading}
          icon={<DownloadOutlined />}
        >
          Capturar Foto
        </Button>
      </Space>
    </Card>
  );
};

export default WebcamCapture;
