# 🏗️ GenPOS Backend — Agent Config

## Project
Go backend for a Point-of-Sale system. Module: `github.com/bagusyanuar/genpos-backend`.

**Stack**: Go · Fiber v2 · GORM · PostgreSQL · Zap · Viper · JWT · golang-migrate

---

## 📂 Structure (Clean Architecture)

```
cmd/api/              → Entry point (config + bootstrap only)
internal/
  [module]/
    domain/           → Entity, Repository interface, Usecase interface
    usecase/          → Business logic (calls repo via interface)
    repository/       → GORM / I/O implementation
    delivery/http/    → Handler + DTO (Req/Res)
  shared/
    bootstrap/        → Server setup, middleware, routes
    container/        → DI wiring (1 file per module)
    config/           → App config (Viper)
    database/         → DB connection
    middleware/        → Auth, logging, etc.
pkg/
  jwt/                → JWT helpers
  request/            → PaginationParam, shared request types
  response/           → Unified API response wrapper
  validator/          → go-playground/validator helpers
migrations/           → SQL up/down files (golang-migrate)
docs/databases/       → DBML schema files
```

---

## 📐 Architecture Rules

- **Domain** (`domain/[m].go`): Entity + `[M]Repository` interface + `[M]Usecase` interface. No external package imports.
- **Usecase** (`usecase/[m]_usecase.go`): Business logic only. Calls repo via interface. Error wrap: `fmt.Errorf("[m]_uc.[fn]: %w", err)`. Log with `config.Log.Error(...)`.
- **Repository** (`repository/[m]_repository.go`): All GORM/DB calls. Error wrap: `fmt.Errorf("[m]_repo.[fn]: %w", err)`.
- **Handler** (`delivery/http/handler.go`): Parse request → call usecase → return response. DTOs in `delivery/http/dto.go`.
- **Container** (`shared/container/[m]_module.go`): Wire repo → usecase → handler for each module.
- **DTO**: Strict separation. Domain Entity ≠ Request/Response struct.

---

## 🗄️ Database Rules

- Naming: `snake_case`, tables **plural**.
- Always use `deleted_at` (soft delete, GORM style).
- Index: `branch_id`, `deleted_at`, all filter columns.
- Migrations: code-based via `golang-migrate`. No manual `ALTER TABLE`.
- FK + Unique constraints enforced.
- Schema source of truth: `docs/databases/*.dbml`.

---

## ⚡ Workflows

### `/create-usecase` — New Module / Feature
1. **Domain** `internal/[m]/domain/[m].go`: Entity + `[M]Repository` + `[M]Usecase` interfaces.
2. **Repository** `internal/[m]/repository/[m]_repository.go`: GORM impl of Repository interface.
3. **Usecase** `internal/[m]/usecase/[m]_usecase.go`: Impl of Usecase interface.
4. **DTO** `internal/[m]/delivery/http/dto.go`: Req/Res structs.
5. **Handler** `internal/[m]/delivery/http/handler.go`: Parse → call UC → respond.
6. **Container** `internal/shared/container/[m]_module.go`: Wire all layers.
7. **Verify**: `go build ./...`

### `/database-migration` — DB Schema Change
1. Read schema in `docs/databases/*.dbml`.
2. `make migrate-create name=<name>` → gen files in `migrations/`.
3. Write `up.sql` (CREATE/ALTER) and `down.sql` (DROP/REVERT).
4. Sync entity in `internal/[m]/domain/[m].go` (tags: `json` + GORM).
5. `go build ./...` for verification.

---

## 🔑 Key Patterns

### Entity (domain layer)
```go
type Foo struct {
    ID        uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
    Name      string         `gorm:"type:varchar(100);not null" json:"name"`
    BranchID  uuid.UUID      `gorm:"type:uuid;index;not null" json:"branch_id"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (f *Foo) BeforeCreate(tx *gorm.DB) (err error) {
    if f.ID == uuid.Nil { f.ID = uuid.New() }
    return
}
```

### Filter + Pagination
```go
type FooFilter struct {
    Search string
    request.PaginationParam  // embedded from pkg/request
}
```

### Error wrapping
```go
// usecase
return nil, fmt.Errorf("foo_uc.FindByID: %w", err)

// repository
return nil, fmt.Errorf("foo_repo.FindByID: %w", err)
```

### Unified Response (pkg/response)
```go
// Success
return response.Success(c, data)

// Error
return response.Error(c, fiber.StatusBadRequest, err)
```

---

## 🛠️ Make Commands

| Command | Description |
|---|---|
| `make migrate-create name=X` | Create new migration files |
| `make migrate-up` | Apply all pending migrations |
| `make migrate-down` | Rollback last migration |
| `make migrate-status` | Show current migration version |
| `make db-seed` | Seed database |

---

## 🪙 Token-Saving Tips

- `@` only files being edited or used as direct reference.
- Focus per layer per session: Domain → Repo → UC → Handler.
- Use `/create-usecase` and `/database-migration` slash commands.
- Skip verbose summaries — read results directly from files/artifacts.
- Use **Plan mode** for multi-layer or ambiguous tasks.
- Never attach `.log`, `.env`, `go.sum` to context.
- Gaya "telegraphic" OK: `"Implement FindAll di branch UC, gas."`.
