CREATE TABLE material_audits (
    id UUID PRIMARY KEY,
    material_id UUID NOT NULL REFERENCES materials(id) ON DELETE CASCADE,
    action VARCHAR(50) NOT NULL,
    note TEXT,
    created_by UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_material_audits_material_id ON material_audits(material_id);
CREATE INDEX idx_material_audits_action ON material_audits(action);
