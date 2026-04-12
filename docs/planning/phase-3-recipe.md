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
- **Live COGS Estimation**:
    - The API must dynamically calculate Estimated COGS based on `(Recipe Qty in Base Unit * Material base_cost) + Variant overhead_cost`.
    - **Override Pattern**: Support manual `subtotal_cost` input per recipe line. If provided (greater than 0), the manual value is used instead of the systemic calculation.
    - This provides realtime Gross Margin feedback to the manager without affecting historical transaction COGS.
- **Data Integrity (Crucial)**:
    - **Restricted Delete**: Materials and UOMs CANNOT be deleted if they are used in active recipes.
    - **Cascaded Delete**: If a Product Variant is soft-deleted, all associated Recipes must also be soft-deleted automatically.
    - **Ubiquitous Language**: Quantities in recipes are entered in a specific `UOM` but must be logically convertible to the Material's `Base Unit`.

## Database Schema Proposal
### `recipes`
- `id` (UUID, PK)
- `product_variant_id` (UUID, FK to product_variants)
- `material_id` (UUID, FK to materials)
- `quantity` (Decimal)
- `uom_id` (UUID, FK to units) - for flexible recipe entry.
- `created_at`, `updated_at`

## Implementation Steps
- [x] **Schema Adjustment**: Add `base_cost` to materials and `overhead_cost` to variants. Setup `recipes` schema.
- [x] **Domain**: Define `Recipe` entity and interfaces in `internal/recipe/domain`.
- [x] **Repository**: Implement DB operations for Recipes.
- [x] **Usecase**: Logic to validate UOM conversion, calculate Live COGS, and prevent recursive recipes.
- [ ] **API**: `GET/POST/PUT/DELETE` for managing recipe of a variant.
