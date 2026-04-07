---
description: Create Usecase
---

1. **Domain**: `.../domain/[m].go` (Interfaces ONLY, NO entity).
2. **Logic**: `.../usecase/[m]_usecase.go`. Wrap: `fmt.Errorf("[m]_uc.[fn]: %w", err)`. Log: `config.Log.Error`.
3. **DTO**: `.../delivery/http/dto.go`.
4. **Handler**: `.../delivery/http/handler.go`.
5. **Module**: `.../shared/container/[m]_module.go` (Wire) & `container.go` (Inject).
6. **Boot**: `.../shared/bootstrap/app.go` (Register).