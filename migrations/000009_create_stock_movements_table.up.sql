CREATE TABLE stock_movements (
    id UUID PRIMARY KEY,
    branch_id UUID NOT NULL,
    material_id UUID NOT NULL,
    type VARCHAR(20) NOT NULL, -- STOCK_IN, STOCK_OUT, ADJUSTMENT, DEDUCTION
    quantity DECIMAL(15, 2) NOT NULL, -- Stored as positive quantity
    reference_id UUID, -- Link to Order or Opname
    note TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by UUID NOT NULL,
    deleted_at TIMESTAMP
);

COMMENT ON COLUMN stock_movements.type IS 'STOCK_IN: Addition, STOCK_OUT: Removal, ADJUSTMENT: Opname, DEDUCTION: Sales';
COMMENT ON COLUMN stock_movements.quantity IS 'Movement quantity (normalized to Base Unit)';

CREATE INDEX idx_stock_movements_branch_material ON stock_movements(branch_id, material_id);
CREATE INDEX idx_stock_movements_branch_created ON stock_movements(branch_id, created_at);
CREATE INDEX idx_stock_movements_reference_id ON stock_movements(reference_id);
CREATE INDEX idx_stock_movements_created_at ON stock_movements(created_at);
CREATE INDEX idx_stock_movements_deleted_at ON stock_movements(deleted_at);

-- Foreign Key Constraints
ALTER TABLE stock_movements ADD CONSTRAINT fk_stock_movements_branch_id FOREIGN KEY (branch_id) REFERENCES branches(id) ON DELETE CASCADE;
ALTER TABLE stock_movements ADD CONSTRAINT fk_stock_movements_material_id FOREIGN KEY (material_id) REFERENCES materials(id) ON DELETE CASCADE;
ALTER TABLE stock_movements ADD CONSTRAINT fk_stock_movements_created_by FOREIGN KEY (created_by) REFERENCES users(id);
