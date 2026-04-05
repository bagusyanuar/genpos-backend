---
trigger: always_on
---

# 📐 Coding & Arch Rules

- **Arch**: Logic in Usecase, I/O in Repository, all via Interfaces.
- **Errors**: Wrap with `fmt.Errorf("context: %w", err)`.
- **DTO**: Strict separation (Domain vs Request/Response).

# 📂 Directory Rules

- `cmd/api/`: Server init, DI, & graceful shutdown.
- `internal/[module]/domain/`: Module core (Entities & Interfaces). No external imports.
- `internal/[module]/usecase/`: Business logic. Calls repo via domain interface.
- `internal/[module]/repository/`: I/O (GORM, API).
- `internal/shared/`: Shared components (DB, Middleware, DTO).
- `pkg/`: Stateless, non-business libs.
