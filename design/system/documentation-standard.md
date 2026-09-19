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

`README.md` files are folder indexes and are exempt. ADRs, work items, and convention files have their own richer schemas; conventions add `audience: agent` and `order`, see [conventions.md](conventions.md). See [work-hierarchy.md](work-hierarchy.md) and the ADR template in `design/adrs/0000-template.md`.

## Naming

- Files and folders are lowercase kebab-case: `repository-layout.md`.
- ADRs are numbered with four digits and a slug: `0003-work-item-hierarchy.md`. `flai adr new` gives the number (one more than the highest file present; gaps are not filled), the slug, the front matter (`id`, `title`, `status`, `date`, `supersedes`, `superseded_by`, and `refines` when given), and the row in `design/adrs/README.md`; the status is `proposed` until `flai adr accept` or `--status accepted`. `flai check` warns with `adr.index` when the files and the index disagree.
- Work items are named by ID and slug: `S-0004-cli-scaffold-and-config.md`. The ID is the stable handle; the slug may change.
- Each folder that a reader might land in has a `README.md` that says what the folder is for.

## Content rules

- Headings start at `#` for the title, then `##`. Do not skip levels.
- Diagrams are Mermaid fenced blocks. No binary images inside `design/` or `wip/`.
- Link between documents with relative paths so links work in git hosting and in the dashboard.
- Dates and times are ISO 8601 in UTC. Dates alone are `YYYY-MM-DD`; timestamps are `YYYY-MM-DDTHH:MM:SSZ`.
- Durations are Go duration strings: `2h30m`, `45m`, `3d` is not valid, use `72h`.

## Editing rules for living documents

- `design/conventions/` files are edited above the marker only through the template (and `flai upgrade`); a project edits below the marker. An agent that thinks a baseline rule is wrong proposes the change, it does not make it.
- `design/system` and `design/tech` are edited in place. If a change reverses an earlier decision, write an ADR first, then update the living document and link the ADR.
- `design/adrs` are never edited after acceptance except to set `superseded_by`. The tooling holds this: `flai doc` refuses the body of an accepted, superseded, or deprecated ADR, in the dashboard's editor too, and only `flai adr new --supersedes` sets `superseded_by` (S-0060). A proposed ADR is a draft and is edited like any document.
- `wip/kanban` items are edited by agents and by `flai`. Human edits are welcome but must keep front matter valid; `flai check` validates it.
- `wip/agents` narratives are append-only in the log section. The summary sections at the top are rewritten as understanding improves.

## Language

Plain, direct prose. Say what is, not what might be. Prefer tables for lists of parallel facts. Keep sentences short.
