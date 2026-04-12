# 📜 Material Implementation History

Record perkembangan modul Material.

## Phase 1: Material Lifecycle (Completed 2026-04-12)
Melengkapi fungsionalitas dasar Material agar siap digunakan secara operasional.

### Features
- [x] **Update Material**: 
    - Mendukung update metadata (nama, kategori, tipe, is_active).
    - **Optimized Update**: Menggunakan `Select` dan `Updates` untuk keamanan data.
    - **Image Deletion Logic**: Otomatis hapus file lama di storage SETELAH DB update berhasil (Resiliency).
- [x] **Patch Image Flow**: Endpoint `PATCH /materials/:id/image` untuk update foto secara terpisah.
- [x] **Update Material UOMs**: Manajemen penambahan atau perubahan konversi satuan via transactional sync.
- [x] **Delete Material**: Implementasi *Soft Delete* untuk menjaga integritas data pada resep (BOM) yang sudah ada.

### Technical Notes
- Multi-tenancy untuk master data Material dilewati (Global/Shared Master), branch_id hanya ada di tabel Inventory.
- Integrasi dengan `pkg/fileupload` untuk manajemen aset gambar.
