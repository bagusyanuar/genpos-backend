---
trigger: pattern: .*\.go$
---

# 📐 Coding & Arch Rules

- **Arch**: Logic in Usecase, I/O in Repo, via Interfaces.
- **Errors**: Wrap with `fmt.Errorf("context: %w", err)`.
- **DTO**: Strict separation (Domain vs Req/Res).
- **Optimization**: Always use `make([]T, 0, cap)` for slices. Use non-destructive Upsert for relational data.

## 📂 Directory Rules

```
cmd/api/              → Entry point (config + bootstrap only)
internal/
  [module]/
    domain/           → Entity, Repository interface, Usecase interface. No external package imports.
    usecase/          → Business logic (calls repo via interface).
    repository/       → GORM / I/O implementation.
    delivery/http/    → Handler + DTO (Req/Res).
  shared/
    bootstrap/        → Server setup, middleware, routes.
    container/        → DI wiring (1 file per module).
    config/           → App config (Viper).
    database/         → DB connection.
    middleware/        → Auth, logging, etc.
pkg/
  jwt/                → JWT helpers.
  request/            → PaginationParam, shared request types.
  response/           → Unified API response wrapper.
  validator/          → go-playground/validator helpers.
migrations/           → SQL up/down files (golang-migrate).
docs/databases/       → DBML schema files.
```

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
return c.Status(fiber.StatusOK).JSON(response.Success(data, "message"))

// Error
return c.Status(fiber.StatusInternalServerError).JSON(response.Error(err.Error()))
```
