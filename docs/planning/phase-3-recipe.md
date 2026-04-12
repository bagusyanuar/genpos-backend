# Phase 3: Recipe & BOM (Bill of Materials)

## Goal
Connect **Materials** (Raw Ingredients) to **Product Variants** (Sellable Items) to enable automated inventory deduction.

## Requirements
- **Recipe Definition**:
    - 1 Variant can have multiple Materials.
    - Each Material in a recipe has a `quantity` and `unit` (must be convertible to Material's **Base Unit**).
- **Inventory Integration**:
    - Deduction should happen during Sales (Future Phase) or Production (Manual).
    - Support for "Waste" factor (optional but good for senior level).

## Database Schema Proposal
### `recipes`
- `id` (UUID, PK)
- `product_variant_id` (UUID, FK to product_variants)
- `material_id` (UUID, FK to materials)
- `quantity` (Decimal)
- `uom_id` (UUID, FK to units) - for flexible recipe entry.
- `created_at`, `updated_at`

## Implementation Steps
1. **Migration**: Create `recipes` table.
2. **Domain**: Define `Recipe` entity and interfaces.
3. **Usecase**: Logic to validate UOM conversion and prevent recursive recipes.
4. **API**: `GET/POST/PUT/DELETE` for managing recipe of a variant.
