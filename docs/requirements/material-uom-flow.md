# 🔄 Material & UOM Management Flow

Dokumen ini menjelaskan alur kerja (workflow) untuk manajemen Material dan Satuan (UOM) pada GenPOS, baik dari sisi UI maupun Backend.

## 🏗️ UI Flow (Stepper)

Untuk mempermudah input data bahan baku yang kompleks, proses **Creation** akan menggunakan **Stepper**:

### Step 1: General Information
- Input data profil material.
- Field: `Name`, `SKU`, `Category`, `Material Type` (RAW/SEMI_FINISHED).
- Status: Draft (belum disimpan ke DB).

### Step 2: Units & Conversions (UOM)
- Pengaturan satuan barang.
- **Base Unit**: Wajib menentukan 1 satuan terkecil (misal: Gram, Ml, Pcs).
    - `Multiplier` otomatis 1.
    - `is_default` = true.
- **Conversion Units**: Menambahkan satuan lain (misal: Pack, Box).
    - User input `Multiplier` (misal: 1 Pack = 1000 Gram).
- API Call: `POST /api/v1/materials` (Atomic Create - Simpan Material + UOMs sekaligus).

### Step 3: Upload Image (Optional)
- Proses upload foto barang.
- Dilakukan setelah Material ID didapatkan dari Step 2.
- User bisa skip jika belum ada foto.
- API Call: `PATCH /api/v1/materials/:id/image` (`multipart/form-data`).

### Step 4: Finish & Review
- Konfirmasi akhir dan ringkasan data.

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

## 📊 Stock Display Format

Dalam menampilkan informasi stok kepada User (Frontend), sistem menerapkan metode **Cascading Hierarchical Formatting**. 

### Konsep Cascading
Sistem akan memecah nilai desimal stok pada Base Unit ke dalam satuan-satuan (UOM) turunannya mulai dari yang **terbesar hingga terkecil**.
- **Aturan Multiplier**: Satuan terbesar memiliki multiplier terbesar (e.g., Liter = 1, Mililiter = 0.001).
- **Proses Rendering**:
  1. Unit diurutkan secara **Descending** berdasarkan multiplier (terbesar ke terkecil).
  2. Nilai stok dibagi dengan multiplier unit yang sedang di-loop.
  3. Nilai bulat (Integer) diambil untuk ditampilkan (misal: 2 Liter).
  4. Sisa desimal dikalikan kembali dengan multiplier untuk merekonstruksi sisa nilai mentah (Base Unit), lalu diumpan ke loop UOM berikutnya.
  5. Proses berhenti saat tidak ada sisa, atau sudah mencapai unit terkecil (diakhiri dengan pembulatan 4 angka desimal untuk membuang anomali floating point).

**Contoh Output API (`formatted_stock`)**:
- *Stok: 2.25 Liter* $\rightarrow$ `["2 Liter", "250 Mililiter"]`
- *Stok: 5.5 Kg* $\rightarrow$ `["5 Kg", "500 Gram"]`

Pola ini dirancang pada sisi Backend (DTO) dalam bentuk **Array of Strings** untuk memanjakan UX Frontend. FE bebas merender array ini menjadi list berurut vertikal atau menggabungkannya `.join(' ')` jika ingin ditampilkan sebaris panjang.
