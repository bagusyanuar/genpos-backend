# Phase 2 History: Product & Menu Module

## Achieved
- **Database Schema**:
    - Created `products`, `product_variants`, and `branch_products` (pivot) tables.
    - Integrated with `categories` and `branches`.
    - Soft delete support and composite indexing (`branch_id`, `deleted_at`, `sku`).
- **Domain & Usecase**:
    - Atomic Creation: Create Product + Variants + Branch Assignment in a single transaction.
    - **Optimization**: Implemented **Non-destructive Upsert** for variants to maintain relational integrity (crucial for future Recipes).
- **Delivery (API)**:
    - **Step 1-3 (Atomic)**: `POST /api/v1/products` handles Info, Variants, and Branch assignments.
    - **Step 4 (Image)**: `PATCH /api/v1/products/:id/image` for separate image upload (Stepper flow).
    - Full CRUD support (`GET`, `PUT`, `DELETE`).
    - Smart filtering by `branch_id` for Cashier visibility.

## Technical Standards Applied
- **Pre-allocated Slices**: Optimized memory usage in loops.
- **Transactional Consistency**: Used `Begin/Commit/Rollback` for all multi-table operations.
- **DTO Separation**: Strict separation between Request/Response and Domain Entities.
