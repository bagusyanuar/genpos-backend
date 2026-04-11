# 🔄 Material & UOM Management Flow

Dokumen ini menjelaskan alur kerja (workflow) untuk manajemen Material dan Satuan (UOM) pada GenPOS, baik dari sisi UI maupun Backend.

## 🏗️ UI Flow (Stepper)

Untuk mempermudah input data bahan baku yang kompleks, proses **Creation** akan menggunakan **Stepper**:

### Step 1: General Information
- Input data profil material.
- Field: `Name`, `SKU`, `Category`, `Material Type` (RAW/SEMI_FINISHED), `Image`.
- Status: Draft (belum disimpan ke DB).

### Step 2: Units & Conversions (UOM)
- Pengaturan satuan barang.
- **Base Unit**: Wajib menentukan 1 satuan terkecil (misal: Gram, Ml, Pcs).
    - `Multiplier` otomatis 1.
    - `is_default` = true.
- **Conversion Units**: Menambahkan satuan lain (misal: Pack, Box).
    - User input `Multiplier` (misal: 1 Pack = 1000 Gram).
- API Call: `POST /api/v1/materials` (Atomic Create - Simpan Material + UOMs sekaligus).

### Step 3: Review & Summary
- Konfirmasi data sebelum finalisasi.

---

## 🛠️ Management Flow (Post-Creation)

Setelah material dibuat, perubahan satuan dilakukan melalui tab khusus:

### Tab: "Satuan & Konversi"
- Menampilkan list UOM yang ada.
- User bisa:
    1. **Tambah Satuan Baru**: Menambahkan konversi baru.
    2. **Edit Satuan**: Mengubah multiplier (Hanya jika stok satuan tersebut 0).
    3. **Hapus Satuan**: 
        - **Soft Delete**: Data tidak hilang dari DB (untuk keperluan histori), tapi tidak muncul di pilihan transaksi.
        - **Guard**: Dilarang hapus jika sedang digunakan di Recipe aktif atau memiliki Inventory record.
- API Call: `PUT /api/v1/materials/:id/uoms` (**Bulk Sync**).

---

## 🔒 Backend Safety Rules

1. **Transactional Sync**: Proses update UOM dilakukan dalam satu transaksi database.
2. **Soft Delete**: Menggunakan `deleted_at` pada tabel `material_uoms`.
3. **Multiplier Guard**: Jika `multiplier` diubah, sistem harus memvalidasi dampaknya terhadap nilai stok yang sudah ada.
4. **Unique Unit**: Dalam satu material, tidak boleh ada `unit_id` yang duplikat.
