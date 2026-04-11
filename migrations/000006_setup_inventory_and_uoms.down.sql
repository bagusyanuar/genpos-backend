ALTER TABLE inventories DROP CONSTRAINT IF EXISTS fk_inventories_branch_id;
ALTER TABLE inventories DROP CONSTRAINT IF EXISTS fk_inventories_material_id;
ALTER TABLE material_uoms DROP CONSTRAINT IF EXISTS fk_material_uoms_unit_id;
ALTER TABLE material_uoms DROP CONSTRAINT IF EXISTS fk_material_uoms_material_id;

DROP TABLE IF EXISTS inventories;
DROP TABLE IF EXISTS material_uoms;
