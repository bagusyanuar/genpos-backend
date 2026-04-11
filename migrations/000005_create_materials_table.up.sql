CREATE TABLE materials (
    id UUID PRIMARY KEY,
    sku VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    UNIQUE(sku, deleted_at)
);

COMMENT ON COLUMN materials.sku IS 'Unique Operational Key';
COMMENT ON COLUMN materials.deleted_at IS 'Soft Delete untuk menjaga relasi resep (BOM)';

CREATE INDEX idx_materials_sku ON materials(sku);
CREATE INDEX idx_materials_deleted_at ON materials(deleted_at);
