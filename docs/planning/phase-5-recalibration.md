# Phase 5: Advanced System Recalibration

## Goal
Provide an Enterprise-level feature to safely change a Material's "Default/Base UOM" dynamically, even after transactions and recipes have been built against it, without breaking the integrity of historical COGS and stock quantities.

## The Challenge
Changing the Base Unit implies that all existing stock quantities and moving average costs are technically in the "wrong" unit. A simple flip of the `is_default` boolean will corrupt the entire warehouse valuation and recipe deductions.

## The Recalibration Engine (Mathematical Approach)
This operation must be executed within a **single, strict database transaction**.

### Example Scenario
- **Old Base Unit**: Gram (`multiplier = 1`)
- **Other UOM**: Kg (`multiplier = 1000`)
- **Current State**:
    - Stock Quantity = 5000 (grams)
    - Base Cost = Rp 15 / gram
- **User Action**: Make `Kg` the new Base Unit.

### Execution Steps:
- [ ] **Find Conversion Factor (CF)**: Get the current `multiplier` of the target UOM. 
- [ ] **Re-normalize Multipliers (`material_uoms`)**: Divide all existing multipliers belonging to this material by the CF.
- [ ] **Convert Remaining Stock (`inventories`)**: Update `New Stock = Old Stock / CF`.
- [ ] **Convert Moving Average Cost (`materials.base_cost`)**: Update `New Base Cost = Old Base Cost * CF`.
- [ ] **Verify Impact on Recipes (`recipes`)**: No Action Needed (relies on `uom_id`).
- [ ] **Verify Impact on Past Sales (`order_items`)**: No Action Needed (relies on hardcoded COGS).

## Requirements
- This feature should be heavily restricted (Superadmin/Manager only).
- An Audit Log must be created mentioning: `Recalibrated UOM from [X] to [Y] with CF=[Z]`.
