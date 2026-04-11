CREATE TABLE material_uoms (
    id UUID PRIMARY KEY,
    material_id UUID NOT NULL,
    unit_id UUID NOT NULL,
    multiplier DECIMAL(15, 4) NOT NULL DEFAULT 1,
    is_default BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);

COMMENT ON COLUMN material_uoms.multiplier IS 'Conversion factor. 1 this unit = multiplier * base unit';
COMMENT ON COLUMN material_uoms.is_default IS 'True if this is the Base Unit (Smallest unit for stock)';

CREATE UNIQUE INDEX idx_material_uoms_unique ON material_uoms(material_id, unit_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_material_uoms_material_id ON material_uoms(material_id);
CREATE UNIQUE INDEX idx_material_uoms_only_one_default ON material_uoms (material_id) WHERE is_default = true AND deleted_at IS NULL;
CREATE INDEX idx_material_uoms_deleted_at ON material_uoms(deleted_at);

CREATE TABLE inventories (
    id UUID PRIMARY KEY,
    material_id UUID NOT NULL,
    branch_id UUID NOT NULL,
    stock DECIMAL(15, 2) NOT NULL DEFAULT 0,
    min_stock DECIMAL(15, 2) NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);

COMMENT ON COLUMN inventories.stock IS 'Current stock level (Strictly in Base Unit)';
COMMENT ON COLUMN inventories.min_stock IS 'Reorder point (Strictly in Base Unit)';

CREATE UNIQUE INDEX idx_inventories_material_branch ON inventories(material_id, branch_id);
CREATE INDEX idx_inventories_branch_id ON inventories(branch_id);
CREATE INDEX idx_inventories_material_id ON inventories(material_id);
CREATE INDEX idx_inventories_deleted_at ON inventories(deleted_at);

-- Foreign Key Constraints
ALTER TABLE material_uoms ADD CONSTRAINT fk_material_uoms_material_id FOREIGN KEY (material_id) REFERENCES materials(id) ON DELETE CASCADE;
ALTER TABLE material_uoms ADD CONSTRAINT fk_material_uoms_unit_id FOREIGN KEY (unit_id) REFERENCES units(id) ON DELETE CASCADE;

ALTER TABLE inventories ADD CONSTRAINT fk_inventories_material_id FOREIGN KEY (material_id) REFERENCES materials(id) ON DELETE CASCADE;
ALTER TABLE inventories ADD CONSTRAINT fk_inventories_branch_id FOREIGN KEY (branch_id) REFERENCES branches(id) ON DELETE CASCADE;
