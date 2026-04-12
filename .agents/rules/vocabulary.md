---
trigger: model_decision: "When discussing domain logic, designing database schemas, or writing architecture specs."
---

# 📖 GenPOS Vocabulary & Business Definitions

Always strictly adhere to these definitions when discussing features, writing code, or naming variables to ensure ubiquitous language across the project.

## 📦 Inventory & Materials
*   **Material (Bahan Baku)**: Raw ingredients purchased from suppliers. Materials are **NOT** sold directly to customers. Examples: "Kopi Biji Roasted", "Gula Pasir", "Susu Segar".
*   **UOM (Unit of Measurement)**: The unit used to measure a material (e.g., Kg, Gram, Liter, Cup).
*   **Base Unit**: The smallest/primary unit used for deducting stock in recipes. All other UOMs must convert to this Base Unit. (e.g., Stock is 1000 Grams, not 1 Kg).
*   **Inventory**: The actual stock tracking table, separated by branch.

## 🍔 Products & Menus
*   **Product (Barang Jualan / Menu)**: The final item displayed on the POS and sold to customers. Examples: "Kopi Susu Gula Aren", "Nasi Goreng Spesial". 
*   **Variant (Varian)**: Options for a specific product, usually affecting price or SKU (e.g., "Regular", "Large", "Hot", "Ice"). Every product **MUST** have at least 1 variant.
*   **Category**: Grouping for Products (e.g., "Main Course", "Beverages").

## 🍳 Production & Recipes
*   **Recipe / BOM (Bill of Materials)**: The exact formula of Materials needed to produce 1 Product Variant. When a Variant is sold, the inventory is deducted based on this Recipe using the Base Unit.

## 🏢 Multi-Tenancy & Structure
*   **Tenant / Branch (Toko/Cabang)**: The operational unit. Almost all transactions, inventory, and product availability are filtered per Branch (`branch_id`).
*   **Branch_Products**: A pivot table determining which Products are available to be sold in which Branch.

## 💳 Transactions (Future)
*   **Order (Pesanan)**: A completed or ongoing transaction on the POS.
*   **Void**: Cancellation of an order or item, which must trigger an audit log and inventory rollback.
