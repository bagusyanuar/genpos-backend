# Phase 6: Core Transaction & Auto-Inventory Deduction

## Goal
Enable the system to process sales orders and automatically deduct material inventory based on recipes. This is the bridge between the Commerce side (Sales) and the Inventory side (Warehouse).

## High-Level Requirements
- **Order Management**: Create, Read, and Update orders (Waitlist/Hold/Paid).
- **Payment Processing**: Multi-method payments (Cash, QRIS, Bank Transfer).
- **Recipe Deductor Engine**: A background or post-transaction service that calculates material usage based on `order_items`.
- **Transaction Integrity**: Ensuring stock doesn't go below zero if not allowed, and atomic updates.

## Proposed Schema Changes
### `orders`
- `id` (uuid)
- `branch_id` (uuid)
- `customer_name` (string)
- `total_amount` (decimal)
- `payment_status` (enum: UNPAID, PAID, CANCELLED)
- `payment_method` (string)

### `order_items`
- `id` (uuid)
- `order_id` (uuid)
- `product_id` (uuid)
- `variant_id` (uuid)
- `quantity` (decimal)
- `price` (decimal)
- `subtotal` (decimal)

## Execution Steps
- [ ] Create Order & OrderItem domain models.
- [ ] Implement Transaction Repository (SQL Transaction for atomic order creation).
- [ ] **Implement StockDeduction Service**: 
    - For each item SOLD -> Lookup Recipe.
    - Calculate `material_id` and `quantity` to deduct.
    - Call `inventory_repo.UpdateStock` (decrement).
- [ ] Integrate stock deduction into the checkout process.
- [ ] Handle "Out of Stock" scenarios in the order flow.
