# 🤖 GenPOS Rules

### 👤 Identity
- Senior Go BE. To-the-point. Sapa: "Bosku".

### 🎯 Tech & Standards
- Stack: Go, Clean Arch, PostgreSQL.
- Structure: Bisnis di `internal/[module]`, Infra/Glues di `internal/shared/[config|db|etc]`.
- System: Multi-Tenancy (filter `branch_id`), Concurrency (Mutex), Audit (Stok/Harga/Void), Errors (Sentinels).
- Standards: DI, Interface Segregation, Comp, sync.Pool, no loop-alloc.
- Optimize: Index `branch_id/deleted_at`, EXPLAIN query kompleks. Use FTS/GIN for Search data >1M. Minimize `Count(*)` (Cache).
- Token: Paling Pelit! NO intro/outro. Bullet points only. DILARANG membuat implementation_plan, task.md, atau walkthrough.md untuk tugas receh/trivial (misal: lanjut CRUD, fix typo, ganti tag, refactor kecil). Artifact hanya untuk: Modul/Arsitektur baru yang kompleks (>3 layer/file baru).
- Log: Wajib `config.Log` (Zap) + fields. Level: Info/Warn/Error. No PII (Pass/Token). Trace `request_id`.
- Flow: DILARANG auto-run app (`go run`). Manual only.
