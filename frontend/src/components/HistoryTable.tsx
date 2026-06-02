import React, { useCallback, useEffect, useState } from 'react';
import { Table, Card, Empty, Button, Tag } from 'antd';
import { ReloadOutlined } from '@ant-design/icons';
import { inventoryService } from '../services/api';

interface InventoryItem {
  id: number;
  catalog_id: number;
  serial_number: string;
  quantity: number;
  location: string;
  status: string;
  received_at: string;
  catalog?: {
    part_number: string;
    manufacturer: string;
    category: string;
  };
}

export const HistoryTable: React.FC = () => {
  const [items, setItems] = useState<InventoryItem[]>([]);
  const [loading, setLoading] = useState(false);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const pageSize = 10;

  const loadItems = useCallback(async () => {
    setLoading(true);
    try {
      const offset = (page - 1) * pageSize;
      const response = await inventoryService.listInventory(pageSize, offset);
      setItems(response.items || []);
      setTotal(response.total || 0);
    } catch (error) {
      console.error('Erro ao carregar histórico:', error);
    } finally {
      setLoading(false);
    }
  }, [page]);

  useEffect(() => {
    loadItems();
  }, [loadItems]);

  const columns = [
    {
      title: 'PN',
      dataIndex: ['catalog', 'part_number'],
      key: 'part_number',
      width: 200,
      ellipsis: true,
    },
    {
      title: 'Serial Number',
      dataIndex: 'serial_number',
      key: 'serial_number',
      width: 150,
      ellipsis: true,
    },
    {
      title: 'Quantidade',
      dataIndex: 'quantity',
      key: 'quantity',
      width: 100,
      align: 'center' as const,
    },
    {
      title: 'Localização',
      dataIndex: 'location',
      key: 'location',
      width: 120,
    },
    {
      title: 'Status',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => {
        const colors: { [key: string]: string } = {
          active: 'green',
          inactive: 'red',
          transferred: 'blue',
        };
        return <Tag color={colors[status] || 'default'}>{status}</Tag>;
      },
    },
    {
      title: 'Data',
      dataIndex: 'received_at',
      key: 'received_at',
      width: 180,
      render: (date: string) => new Date(date).toLocaleString('pt-BR'),
    },
  ];

  if (items.length === 0 && !loading) {
    return (
      <Card title="Histórico">
        <Empty description="Nenhum item no estoque" />
      </Card>
    );
  }

  return (
    <Card 
      title="Histórico de Estoque"
      extra={
        <Button 
          type="text" 
          icon={<ReloadOutlined />}
          onClick={loadItems}
          loading={loading}
        >
          Atualizar
        </Button>
      }
    >
      <Table
        columns={columns}
        dataSource={items.map(item => ({ ...item, key: item.id }))}
        loading={loading}
        pagination={{
          current: page,
          pageSize,
          total,
          onChange: setPage,
        }}
        size="small"
      />
    </Card>
  );
};

export default HistoryTable;
