-- Inicialização do banco de dados PostgreSQL

-- Criar tabela de Catálogo
CREATE TABLE IF NOT EXISTS catalog (
    id BIGSERIAL PRIMARY KEY,
    part_number VARCHAR(200) UNIQUE NOT NULL,
    serial_pattern VARCHAR(200),
    manufacturer VARCHAR(100),
    category VARCHAR(100),
    normalized_description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Criar tabela de Estoque
CREATE TABLE IF NOT EXISTS inventory (
    id BIGSERIAL PRIMARY KEY,
    catalog_id BIGINT NOT NULL REFERENCES catalog(id) ON DELETE CASCADE,
    serial_number VARCHAR(200),
    quantity INTEGER DEFAULT 1,
    location VARCHAR(100),
    status VARCHAR(50) DEFAULT 'active',
    received_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Criar tabela de Movimentações
CREATE TABLE IF NOT EXISTS movements (
    id BIGSERIAL PRIMARY KEY,
    inventory_id BIGINT NOT NULL REFERENCES inventory(id) ON DELETE CASCADE,
    quantity INTEGER NOT NULL,
    operation VARCHAR(50) NOT NULL, -- 'IN', 'OUT', 'TRANSFER', 'ADJUSTMENT'
    reason TEXT,
    user_id BIGINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Criar tabela de Auditoria
CREATE TABLE IF NOT EXISTS audit_log (
    id BIGSERIAL PRIMARY KEY,
    action VARCHAR(100) NOT NULL,
    entity_type VARCHAR(50),
    entity_id BIGINT,
    old_values JSONB,
    new_values JSONB,
    user_id BIGINT,
    ip_address INET,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Criar tabela de Feedback para Aprendizado Contínuo
CREATE TABLE IF NOT EXISTS feedback_samples (
    id BIGSERIAL PRIMARY KEY,
    part_number_predicted VARCHAR(200),
    part_number_final VARCHAR(200),
    serial_number_predicted VARCHAR(200),
    serial_number_final VARCHAR(200),
    manufacturer_predicted VARCHAR(100),
    manufacturer_final VARCHAR(100),
    category_predicted VARCHAR(100),
    category_final VARCHAR(100),
    correction_applied BOOLEAN NOT NULL DEFAULT FALSE,
    confidence NUMERIC(6,4),
    image_data TEXT,
    ocr_text JSONB,
    meta_json JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Criar índices para performance
CREATE INDEX IF NOT EXISTS idx_catalog_pn ON catalog(part_number);
CREATE INDEX IF NOT EXISTS idx_catalog_manufacturer ON catalog(manufacturer);
CREATE INDEX IF NOT EXISTS idx_catalog_category ON catalog(category);
CREATE INDEX IF NOT EXISTS idx_inventory_catalog ON inventory(catalog_id);
CREATE INDEX IF NOT EXISTS idx_inventory_serial ON inventory(serial_number);
CREATE INDEX IF NOT EXISTS idx_inventory_status ON inventory(status);
CREATE INDEX IF NOT EXISTS idx_movements_inventory ON movements(inventory_id);
CREATE INDEX IF NOT EXISTS idx_movements_created ON movements(created_at);
CREATE INDEX IF NOT EXISTS idx_audit_created ON audit_log(created_at);
CREATE INDEX IF NOT EXISTS idx_feedback_created ON feedback_samples(created_at);
CREATE INDEX IF NOT EXISTS idx_feedback_correction ON feedback_samples(correction_applied);

-- Inserir alguns Part Numbers de teste
INSERT INTO catalog (part_number, manufacturer, category, normalized_description)
VALUES
  ('M393A4K40DB3-CWE', 'Samsung', 'memory', 'DDR4 32GB 3200MHz RDIMM ECC'),
  ('HUH721212AL5200', 'HGST', 'disk', 'HDD 12TB SAS 7200RPM 256MB Cache'),
  ('INTEL_SSDPE2MX450G7', 'Intel', 'ssd', 'NVMe SSD 450GB 3D TLC PCI-E 3.0'),
  ('E82968-001', 'HP', 'network', 'NIC 10Gb Ethernet Adapter Dual Port'),
  ('SL230SG8-4U', 'SuperMicro', 'memory', 'DDR4 Memory Slot 1'),
  ('SEAG7000NM000F', 'Seagate', 'disk', 'HDD 7200RPM SAS 2TB'),
  ('INTEL_750_SSDPEDMW400G4X9', 'Intel', 'ssd', 'SSD NVMe 400GB U.2'),
  ('Q26142-L21', 'HP', 'network', 'NIC 10Gb Ethernet 2-port 546FLR-SFP28 Adapter'),
  ('AD3-00057896', 'ADI', 'raid', 'RAID Controller 12Gb/s SAS'),
  ('UCSC-PCIE-QSFP28', 'Cisco', 'network', 'QSFP28 Module Adapter')
ON CONFLICT (part_number) DO NOTHING;

-- Inserir alguns itens de estoque de teste
INSERT INTO inventory (catalog_id, serial_number, quantity, location, status)
SELECT id, CONCAT('SN-', SUBSTRING(MD5(part_number), 1, 12)), 1, 'DC-RJ', 'active'
FROM catalog
LIMIT 5;
