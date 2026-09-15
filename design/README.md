# design

Internal documentation for the system-flow monorepo. Nothing here is written for end users; that lives in `docs/`.

| Folder | Purpose | Lifecycle |
|--------|---------|-----------|
| `adrs/` | Architecture Decision Records. One decision, at one point in time, per file. | Immutable once accepted. Superseded by newer ADRs, never edited in place. |
| `system/` | The living design of the system as it is right now. | Edited whenever a conversation resolves new details, direction, standards, or decisions. |
| `tech/` | Every active technology choice: version, where it is used, and why. | Edited whenever a dependency is added, upgraded, or removed. |

Rules that apply to everything under `design/`:

- Markdown only. Diagrams are Mermaid fenced blocks so they render in the dashboard.
- Every file carries YAML front matter with at least `title` and `updated`. See [system/documentation-standard.md](system/documentation-standard.md).
- When an ADR changes a decision, update the affected `system/` and `tech/` documents in the same change.
