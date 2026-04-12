# 🚀 GenPOS — Next Implementation Planning

Dokumen ini mencatat rencana pengembangan fitur setelah implementasi Atomic Create Material dan Shared Upload System.

## 📋 Roadmaps
Progress yang sudah selesai dipindahkan ke [docs/history/](file:///d:/project-go/genpos-backend/docs/history/) untuk menjaga kebersihan roadmap.

## 🍔 Phase 2: Product & Menu Module
Manajemen item yang akan dijual ke pelanggan.

- [ ] **Product Metadata**: 
    - Nama, kategori produk, deskripsi, image.
- [ ] **Variant Support**: 
    - Implementasi tabel `product_variants`.
    - Tiap produk minimal 1 varian (Default/Regular).
    - Support harga berbeda per varian (misal: S/M/L).
- [ ] **Branch Mapping**: 
    - Implementasi Many-to-Many via `branch_products`.
    - Filter ketersediaan menu per cabang.

## 🧪 Phase 3: Recipe & Production (BOM)
Menghubungkan Material (Bahan Baku) dengan Menu Penjualan.

- [ ] **Recipe Management**:
    - Definisi komposisi per **Varian Produk** (Contoh: 1 Cup Kopi Ice = 18gr Biji Kopi + 150ml Air).
    - Support untuk *Sub-Recipe* (Bahan olahan/Semi-finished).
- [ ] **Auto-Deduction Engine**: Integrasi dengan modul Transaksi untuk memotong stok otomatis saat produk terjual.

## 📉 Phase 4: Inventory Advanced
Fitur lanjutan untuk akurasi stok.

- [ ] **Stock Adjustment**: Input manual untuk stok masuk (PO) atau stok rusak/hilang.
- [ ] **Stock Opname**: Fitur verifikasi stok fisik vs sistem secara berkala.
- [ ] **Low Stock Alert**: Notifikasi jika bahan baku mencapai ambang batas minimum.

---

> [!NOTE]
> Seluruh implementasi harus tetap mengikuti **Clean Architecture** dan menggunakan `pkg/fileupload` untuk setiap aset gambar yang diunggah.
