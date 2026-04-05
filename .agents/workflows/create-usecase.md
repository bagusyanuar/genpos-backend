---
description: Create Usecase
---

1. **Domain**: `.../domain/[m].go` -> Interface `Usecase` & `Repository`.
2. **Logic**: `.../usecase/[m]_usecase.go` -> Implement. 
   - Wrap: `fmt.Errorf("[m]_uc.[fn]: %w", err)`
   - Log: `config.Log.Error("...", zap.Error(err))`
3. **DTO**: `.../delivery/http/dto.go` -> Req/Res models.
4. **Handler**: `.../delivery/http/handler.go` -> Logic & Route.