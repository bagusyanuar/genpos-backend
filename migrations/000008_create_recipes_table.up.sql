CREATE TABLE recipes (
    id UUID PRIMARY KEY,
    product_variant_id UUID NOT NULL,
    material_id UUID NOT NULL,
    uom_id UUID NOT NULL,
    quantity DECIMAL(15, 4) NOT NULL DEFAULT 0,
    subtotal_cost DECIMAL(15, 4) NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    CONSTRAINT fk_recipes_product_variant_id FOREIGN KEY (product_variant_id) REFERENCES product_variants(id) ON DELETE CASCADE,
    CONSTRAINT fk_recipes_material_id FOREIGN KEY (material_id) REFERENCES materials(id) ON DELETE RESTRICT,
    CONSTRAINT fk_recipes_uom_id FOREIGN KEY (uom_id) REFERENCES units(id) ON DELETE RESTRICT,
    UNIQUE(product_variant_id, material_id, deleted_at)
);

CREATE INDEX idx_recipes_product_variant_id ON recipes(product_variant_id);
CREATE INDEX idx_recipes_material_id ON recipes(material_id);
CREATE INDEX idx_recipes_deleted_at ON recipes(deleted_at);

COMMENT ON TABLE recipes IS 'Bill of Materials (BOM) linking product variants to raw materials';
COMMENT ON COLUMN recipes.uom_id IS 'Unit used in this recipe, must be convertible to Material''s base unit';
