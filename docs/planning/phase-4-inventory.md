# Phase 4: Inventory & Stock Movement

## Goal
Manage physical stock levels across branches through manual adjustments, stock takes (opname), and tracking all movements for auditing.

## Requirements
- **Stock Movement Types**:
    - `STOCK_IN`: Adding stock (e.g., from purchase or manual addition).
    - `STOCK_OUT`: Removing stock (e.g., waste, expired, or manual removal).
    - `ADJUSTMENT`: Correction from Stock Opname.
    - `DEDUCTION`: Automated deduction from Sales (via Recipe/BOM).
- **Stock Opname**:
    - Periodic counting of physical stock to sync with system stock.
    - Captures "System Stock", "Actual Stock", and "Difference".
- **Audit Log**:
    - Every change in `inventory.quantity` must record who, when, why, and the delta.

## Planned Database Schema
### `stock_movements`
- `id` (UUID)
- `branch_id`, `material_id`
- `type` (ENUM: IN, OUT, ADJUST, DEDUCTION)
- `quantity` (Decimal) - normalized to Base Unit.
- `reference_id` (UUID, optional) - Link to Order or Opname.
- `note` (Text)
- `created_at`, `created_by` (UUID)

## Implementation Steps
1. **Migration**: Create `stock_movements` table.
2. **Domain**: Define `StockMovement` and `InventoryAdjustment` entities.
3. **Usecase**: Logic for atomic stock updates (Inventory + Movement Log).
4. **API**: Endpoints for manual stock in/out and Opname recording.
