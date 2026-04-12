# 🚀 GenPOS — Next Implementation Planning

Dokumen ini mencatat rencana pengembangan fitur setelah implementasi Atomic Create Material dan Shared Upload System.

## 📦 Phase 1: Material Lifecycle (CRUD Completion)
Melengkapi fungsionalitas dasar Material agar siap digunakan secara operasional.

- [x] **Update Material**: 
    - Mendukung update metadata (nama, kategori, tipe, is_active).
    - [x] **Patch Image Flow**: Proses upload foto dilakukan di step terakhir atau secara terpisah via `PATCH /materials/:id/image`.
    - [x] Mendukung hapus foto lama secara otomatis jika diganti (interaksi dengan `pkg/fileupload`).
- [x] **Update Material UOMs**: Manajemen penambahan atau perubahan konversi satuan.
- [x] **Delete Material**: Implementasi *Soft Delete* untuk menjaga integritas data pada resep (BOM) yang sudah ada.

## 🧪 Phase 2: Recipe & Production (BOM)
Menghubungkan Material (Bahan Baku) dengan Menu Penjualan.

- [ ] **Recipe Management**:
    - Definisi komposisi (Contoh: 1 Cup Kopi = 18gr Biji Kopi + 150ml Air).
    - Support untuk *Sub-Recipe* (Bahan olahan/Semi-finished).
- [ ] **Auto-Deduction Engine**: Integrasi dengan modul Transaksi untuk memotong stok otomatis saat produk terjual.

## 🍔 Phase 3: Product & Menu Module
Manajemen item yang akan dijual ke pelanggan.

- [ ] **Product Metadata**: Nama, kategori produk, harga jual, pajak.
- [ ] **Variant Support**: Ukuran (Small/Large) atau pilihan (Hot/Ice).
- [ ] **Link to Recipe**: Menghubungkan produk/varian dengan resep yang sudah dibuat di Phase 2.

## 📉 Phase 4: Inventory Advanced
Fitur lanjutan untuk akurasi stok.

- [ ] **Stock Adjustment**: Input manual untuk stok masuk (PO) atau stok rusak/hilang.
- [ ] **Stock Opname**: Fitur verifikasi stok fisik vs sistem secara berkala.
- [ ] **Low Stock Alert**: Notifikasi jika bahan baku mencapai ambang batas minimum.

---

> [!NOTE]
> Seluruh implementasi harus tetap mengikuti **Clean Architecture** dan menggunakan `pkg/fileupload` untuk setiap aset gambar yang diunggah.
