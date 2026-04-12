---
trigger: model_decision: "When user asks about available scripts, database seeds, or migration commands."
---

# 🛠️ Make Commands

Use these commands for database management and seeding.

| Command | Description |
|---|---|
| `make migrate-create name=X` | Create new migration files |
| `make migrate-up` | Apply all pending migrations |
| `make migrate-down` | Rollback last migration |
| `make migrate-status` | Show current migration version |
| `make db-seed` | Seed database |
