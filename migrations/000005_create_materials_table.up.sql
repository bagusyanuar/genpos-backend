CREATE TABLE materials (
    id UUID PRIMARY KEY,
    category_id UUID,
    sku VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    material_type VARCHAR(50),
    image_url TEXT,
    base_cost DECIMAL(15,2) NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    UNIQUE(sku, deleted_at)
);

COMMENT ON COLUMN materials.sku IS 'Unique Operational Key';
COMMENT ON COLUMN materials.deleted_at IS 'Soft Delete untuk menjaga relasi resep (BOM)';

CREATE INDEX idx_materials_sku ON materials(sku);
CREATE INDEX idx_materials_category_id ON materials(category_id);
CREATE INDEX idx_materials_is_active ON materials(is_active);
CREATE INDEX idx_materials_deleted_at ON materials(deleted_at);

-- Foreign Key Constraints
ALTER TABLE materials ADD CONSTRAINT fk_materials_category_id FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE SET NULL;
