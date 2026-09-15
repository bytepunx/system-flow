---
title: Documentation standard
updated: 2026-09-15
status: active
---

# Documentation standard

Applies to every markdown file in `design/`, `docs/`, and `wip/`, in this repo and in any project generated from the template.

## Front matter

All files begin with a YAML front matter block delimited by `---`. The dashboard indexes front matter, so a file without it is rendered but not searchable by metadata.

Minimum keys for `design/` and `docs/`:

```yaml
---
title: Human readable title
updated: 2026-09-15        # ISO date, bumped on every meaningful edit
status: active             # active | draft | deprecated
---
```

`README.md` files are folder indexes and are exempt. ADRs and work items have their own richer schemas. See [work-hierarchy.md](work-hierarchy.md) and the ADR template in `design/adrs/0000-template.md`.

## Naming

- Files and folders are lowercase kebab-case: `repository-layout.md`.
- ADRs are numbered with four digits and a slug: `0003-work-item-hierarchy.md`.
- Work items are named by ID and slug: `S-004-cli-scaffold-and-config.md`. The ID is the stable handle; the slug may change.
- Each folder that a reader might land in has a `README.md` that says what the folder is for.

## Content rules

- Headings start at `#` for the title, then `##`. Do not skip levels.
- Diagrams are Mermaid fenced blocks. No binary images inside `design/` or `wip/`.
- Link between documents with relative paths so links work in git hosting and in the dashboard.
- Dates and times are ISO 8601 in UTC. Dates alone are `YYYY-MM-DD`; timestamps are `YYYY-MM-DDTHH:MM:SSZ`.
- Durations are Go duration strings: `2h30m`, `45m`, `3d` is not valid, use `72h`.

## Editing rules for living documents

- `design/system` and `design/tech` are edited in place. If a change reverses an earlier decision, write an ADR first, then update the living document and link the ADR.
- `design/adrs` are never edited after acceptance except to set `superseded_by`.
- `wip/kanban` items are edited by agents and by `flai`. Human edits are welcome but must keep front matter valid; `flai check` validates it.
- `wip/agents` narratives are append-only in the log section. The summary sections at the top are rewritten as understanding improves.

## Language

Plain, direct prose. Say what is, not what might be. Prefer tables for lists of parallel facts. Keep sentences short.
