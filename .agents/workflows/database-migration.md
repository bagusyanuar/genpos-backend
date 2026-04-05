---
description: DB Migration via DBML
---

1. **Read DBML**: Cek schema di `docs/databases/*.dbml`.
2. **Create**: `make migrate-create name=...` untuk gen file.
3. **SQL**: Isi `up.sql` (create) & `down.sql` (drop) di `migrations/`.
4. **Sync Entity**: Update `internal/[module]/domain/` (Entity & Tags).
   - Pakai `gorm.Model`.
   - Pastiin tags `json` & `gorm` sinkron.
5. **Verify**: `go build ./...` untuk cek error.
