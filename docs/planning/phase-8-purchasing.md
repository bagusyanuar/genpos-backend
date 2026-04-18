# Phase 8: Purchasing & Supplier Management

## Goal
Formalize the stock-in process via Purchase Orders (PO) and Goods Receiving (GR). This provides a more accurate audit trail of material sourcing and costs.

## High-Level Requirements
- **Supplier Directory**: Manage vendors and their contact/payment info.
- **Purchase Orders**: Create POs to suppliers.
- **Goods Receiving**: Partial or full receiving of PO items with automatic inventory increments.
- **Cost Tracking**: Ensure `base_cost` (Moving Average) reflects the latest purchase prices.

## Proposed Schema Changes
### `suppliers`
- `id` (uuid)
- `name` (string)
- `email/phone` (string)
- `address` (text)

### `purchase_orders`
- `id` (uuid)
- `branch_id` (uuid)
- `supplier_id` (uuid)
- `status` (DRAFT, SENT, PARTIAL, COMPLETED)
- `total_price` (decimal)

### `purchase_order_items`
- `id` (uuid)
- `po_id` (uuid)
- `material_id` (uuid)
- `quantity` (decimal)
- `unit_price` (decimal)

## Execution Steps
- [ ] Create Supplier management module.
- [ ] Implement Purchase Order workflow (Approval flow).
- [ ] Implement Goods Receiving logic (GRN - Goods Received Note).
- [ ] **Cost Engine**: Calculate New Moving Average Cost on items receipt.
- [ ] Automatically create `STOCK_IN` movements on receipt.
