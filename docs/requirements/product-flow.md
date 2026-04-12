# 🍕 Product & Variant Management Flow

Dokumen ini menjelaskan alur kerja (workflow) untuk manajemen Produk dan Variannya di GenPOS. Sistem didesain untuk mendukung industri F&B yang memiliki banyak opsi dalam satu item menu.

## 🏗️ UI Flow (Creation)

Berbeda dengan Material, Produk memiliki struktur berjenjang (Product -> Variants).

### Step 1: Base Information
- Input data utama produk.
- Field: `Name`, `Category` (must be type PRODUCT), `Description`, `Image`.
- API: `POST /api/v1/products` (Simpan data header).

### Step 2: Variant & Pricing
- Menentukan pilihan produk. 
- Setiap produk minimal memiliki **1 Variant Default** (misal: "Regular").
- Field per Variant:
    - `Variant Name`: (misal: "Small", "Large", "Hot", "Ice").
    - `Price`: Harga jual.
    - `SKU`: Manual input atau auto-generate (Unique PK).
    - `is_active`: Status aktif per varian.
- API: `POST /api/v1/products/:id/variants` (Bulk Create).

### Step 3: Link to Recipe (Future Phase)
- Setiap varian dapat dihubungkan ke **Recipe (BOM)** secara spesifik.
- Contoh: Kopi Susu "Small" dan "Large" memiliki resep (pemakaian bahan baku) yang berbeda.

---

## 🛠️ Data Structure (Preview)

### Table: `products`
- `id` (UUID)
- `category_id` (UUID)
- `name` (String)
- `description` (Text)
- `image_url` (String)
- `is_active` (Bool)

### Table: `product_variants`
- `id` (UUID)
- `product_id` (UUID)
- `name` (String) - Contoh: "S", "M", "L"
- `sku` (String, Unique)
- `price` (Decimal)
- `is_active` (Bool)

### Table: `branch_products` (Many-to-Many)
- `branch_id` (UUID)
- `product_id` (UUID)
- `is_active` (Bool) - Default true. Cabang bisa menonaktifkan produk ini secara mandiri.
*(Catatan: Mengatur ketersediaan produk per cabang. Jika record tidak ada, berarti cabang tersebut tidak menjual produk ini).*

---

## 🔒 Business Rules (Guardrails)

1. **Unique SKU**: SKU varian harus unik di seluruh sistem (global).
2. **Minimal 1 Variant**: Setiap produk wajib memiliki minimal satu varian agar bisa muncul di Menu Penjualan.
3. **Category Lock**: Hanya bisa memilih kategori yang bertipe `PRODUCT`.
4. **Branch Availability**: Ketersediaan produk di POS Kasir difilter berdasarkan tabel `branch_products`. Master (Pusat) menentukan produk apa saja yang di-assign ke cabang mana.
5. **Soft Delete**: Menghapus produk atau varian menggunakan `deleted_at`. Jika produk dihapus, seluruh variannya otomatis ikut ditandai terhapus.
6. **Recipe Integrity**: Varian yang sudah memiliki record transaksi atau resep aktif tidak boleh dihapus (hanya bisa di-nonaktifkan).

---

## 🔄 Lifecycle Example
**Item: Kopi Susu Gula Aren**
- Variant 1: `Ice` | Price: `20.000` | SKU: `Kopi-Ice-01`
- Variant 2: `Hot` | Price: `18.000` | SKU: `Kopi-Hot-01`

Jika stok Bahan Baku di Inventory menipis, sistem akan mengecek Recipe yang terhubung ke masing-masing Variant ini untuk memberikan alert.
