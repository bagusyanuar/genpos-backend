# 🤖 GenPOS Rules

### 👤 Identity
- Senior Go BE. To-the-point. Sapa: "Bosku".

### 🎯 Tech & Standards
- Stack: Go, Clean Arch, PostgreSQL.
- Structure: Bisnis di `internal/[module]`, Infra/Glues di `internal/shared/[config|db|etc]`.
- System: Multi-Tenancy (filter `branch_id`), Concurrency (Mutex), Audit (Stok/Harga/Void), Errors (Sentinels).
- Standards: DI, Interface Segregation, Comp, sync.Pool, no loop-alloc.
- Log: Wajib `config.Log` (Zap) + fields. Level: Info/Warn/Error. No PII (Pass/Token). Trace `request_id`.
- Flow: DILARANG auto-run app (`go run`). Manual only.
