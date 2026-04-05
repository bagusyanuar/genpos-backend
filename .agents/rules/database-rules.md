---
trigger: always_on
---

# 🗄️ Database Rules

- Naming: Snake_case, tables plural.
- Indexing: Wajib index pada `branch_id`, `deleted_at`, & filter columns jika exist.
- Soft Deletes: Gunakan `deleted_at` (GORM style).
- Migrations: No manual `ALTER TABLE` in prod. Always code-based.
- Integrity: Strict Foreign Keys & Unique constraints.
