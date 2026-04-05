---
trigger: always_on
---

# 📐 Coding & Arch Rules

- **Arch**: Logic in Usecase, I/O in Repo, via Interfaces.
- **Errors**: Wrap with `fmt.Errorf("context: %w", err)`.
- **DTO**: Strict separation (Domain vs Req/Res).

# 📂 Directory Rules

- `cmd/api/`: App entry point. Only loads config & calls bootstrap.
- `internal/[module]/domain/`: Module core (Entities & Interfaces). No external imports.
- `internal/[module]/usecase/`: Business logic. Calls repo via domain interface.
- `internal/[module]/repository/`: I/O (GORM, API).
- `internal/shared/container/`: Centralized DI wiring for all modules.
- `internal/shared/bootstrap/`: Server setup, middleware, & route registry.
- `pkg/`: Stateless helpers.
