---
trigger: pattern: .*\.(sql|dbml)$
---

# 🗄️ Database Rules

- **Naming**: `snake_case`, tables **plural**.
- **Soft Deletes**: Always use `deleted_at` (GORM style).
- **Indexing**: Wajib index pada `branch_id`, `deleted_at`, & filter columns.
- **Integrity**: Strict Foreign Keys & Unique constraints enforced.
- **Migrations**: Always code-based via `golang-migrate`. No manual `ALTER TABLE` in production.
- **Schema Source of Truth**: `docs/databases/*.dbml`.


