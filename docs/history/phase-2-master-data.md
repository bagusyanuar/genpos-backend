# Phase 2: Master Data & Multi-Tenancy

## Goal
Establish the core business entities (Branches, Materials, Products) and implement the multi-tenancy layer.

## Key Features
- [x] **Multi-Tenancy (Branch)**:
    - [x] Every branch is a siloed entity.
    - [x] Global middleware to inject `branch_id` from JWT or Request Header.
- [x] **Material Management**:
    - [x] Centralized material master (Shared assets).
    - [x] Transactional sync for Multiple Units of Measurement (UOM).
    - [x] Image management with automatic cleanup on update/delete.
- [x] **Product & Menu**:
    - [x] **Atomic Creation**: Create Product + Variants + Branch overrides in a single transaction.
    - [x] **Non-destructive Upsert**: High integrity variant management for future recipe referencing.
    - [x] Cashier-optimized query with branch-specific visibility.

## Technical Milestones
- **Soft Delete Implementation**: Protects historical report integrity when master data is "removed".
- **Transactional Consistency**: Strict use of GORM transactions for all multi-table operations.
- **Pre-allocated Slices**: Optimized memory performance for large batch operations.
