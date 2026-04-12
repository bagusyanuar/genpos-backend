---
trigger: model_decision: "When executing multi-step scaffolding, code optimization, relational upserts, or schema syncing."
---

# ⚔️ GenPOS Agent Skills

This document defines standard operating procedures (SOPs) and macro-behaviors. Reference these skills to execute complex, multi-step workflows predictably.

## 1. Skill: Update-Pattern-Upsert
**Trigger**: When dealing with One-to-Many nested updates in the Usecase layer.
**Execution Steps**:
1. **Diffing**: Query existing nested data.
2. **Upsert Loop**: Update if incoming data has an ID, Create if incoming data has no ID.
3. **Soft Delete**: Delete existing data that is missing from the incoming request.
4. **Transaction**: Wrap everything in a single DB transaction (`Begin()`, `Commit()`, `Rollback()`).

*(Note: Scaffolding and DB Sync procedures have been moved to `.agents/workflows/`, and basic code optimization rules reside in `coding-arch-rules.md`.)*
