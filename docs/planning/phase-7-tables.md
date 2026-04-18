# Phase 7: Table & Floor Management

## Goal
Manage physical sitting arrangements and order assignments for F&B operations. This allows the POS to track "Open Bills" per table.

## High-Level Requirements
- **Floor Plan**: Define areas (Indoor, Outdoor, VIP).
- **Table Management**: Define tables per area with capacity.
- **Table Status tracking**: Available, Occupied, Reserved, Cleaning.
- **Order Linking**: Every order in a dine-in scenario must be linked to a `table_id`.

## Proposed Schema Changes
### `tables`
- `id` (uuid)
- `branch_id` (uuid)
- `name` (string)
- `area` (string)
- `status` (enum)
- `current_order_id` (uuid, nullable)

## Execution Steps
- [ ] Create Table & Area domain models.
- [ ] Implement Table Management API (CRUD).
- [ ] Implement Table Status switching logic.
- [ ] Update Order flow to support `dine_in` with `table_id` reference.
- [ ] Implement "Move Table" and "Split/Merge Table" logic.
