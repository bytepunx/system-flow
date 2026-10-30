---
title: Documentation standard
updated: 2026-10-30
status: active
topics: [all]
---

# Documentation standard

Applies to every markdown file in `design/`, `docs/`, and `wip/`, in this repo and in any project generated from the template.

Every document is markdown with YAML front matter, with two exceptions, both JSON with no front matter:

- A project's design token file is JSON in the W3C Design Tokens Community Group (DTCG) format, at the path the manifest's `design.tokens` key names, `design/system/design-system.tokens.json` by default ([ADR-0146](../adrs/0146-a-project-s-design-tokens-are-a-dtcg-file-the-manifest-names-beside-its.md)).
- A component's examples, named sets of its props, are one file per component under `design/components/`, at `design/components/<sub-project>/<path>`: the sub-project's manifest `name`, then the component's path in it with `.svelte` replaced by `.examples.json` ([ADR-0164](../adrs/0164-a-component-s-examples-are-named-sets-of-props-in-one-json-file-per-component.md), [component-library.md § Examples](component-library.md#examples-s-0443)). The folder's `README.md` is markdown, as every index is.

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

`design/system` and `design/tech` files also carry `topics`, the stories they are for, and conventions and ADRs may ([Topics](conventions.md#topics)); `flai check` warns about a design or tech file without them.

`README.md` files are folder indexes and are exempt. ADRs, work items, and convention files have their own richer schemas; conventions add `audience: agent` and `order`, see [conventions.md](conventions.md). See [work-hierarchy.md](work-hierarchy.md) and the ADR template in `design/adrs/0000-template.md`.

## Naming

- Files and folders are lowercase kebab-case: `repository-layout.md`.
- ADRs are numbered with four digits and a slug: `0003-work-item-hierarchy.md`. `flai adr new` gives the number (one past the highest ADR in the checkout, the main checkout, every story's worktree, the main branch, and every story branch, so parallel stories do not take the same number; gaps are not filled), the slug, the front matter (`id`, `title`, `status`, `date`, `supersedes`, `superseded_by`, and `refines` when given), and the row in `design/adrs/README.md`; the status is `proposed` until `flai adr accept` or `--status accepted`. `flai check` warns with `adr.index` when the files and the index disagree.
- Work items are named by ID and slug: `S-0004-cli-scaffold-and-config.md`. The ID is the stable handle; the slug may change.
- A file name flai makes from a title cuts the title's slug at 80 bytes on a word: back to the last dash when it lies past 40 bytes, else at 80. This holds for work items (made or retitled), issues, threads, conversations, ADRs, and digest entries. So the name and a tool's temporary file beside it, such as Claude Code's `<name>.tmp.<pid>.<hex>`, stay under the 255-byte file name limit ([I-0231](../issues/I-0231-claude-code-s-edit-tool-cannot-write-a-work-item-whose-file-name-is-near-the-255-byte-limit-since-its-temporary-file-name-adds-about-25-bytes.md)). A name made before keeps its slug: `flai check` asks only that a name start with its ID.
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
- `design/adrs` are never edited after acceptance except to set `superseded_by` or `topics` ([ADR-0047](../adrs/0047-an-agent-is-primed-with-what-its-story-s-topics-claim-and-links-select.md)). The tooling holds this: `flai doc` refuses any other change to an accepted, superseded, deprecated, or rejected ADR, and the dashboard's editor shows one read-only; only `flai adr new --supersedes` sets `superseded_by` (S-0060), and `flai adr topics` sets `topics` on an ADR of any status (S-0134). A proposed ADR is a draft and is edited like any document.
- An ADR's status is one of `proposed`, `accepted`, `superseded`, `deprecated`, or `rejected` ([ADR-0223](../adrs/0223-an-adr-s-status-is-one-of-proposed-accepted-superseded-deprecated-or-rejected-a.md)). A rejected ADR is a proposed decision the operator declined in favour of another. It is immutable once set, as an accepted one is. `flai check` reports `adr.status` on an ADR whose status is outside the five; the template, `0000-template.md`, is passed over.
- `wip/kanban` items are edited by agents and by `flai`. Human edits are welcome but must keep front matter valid; `flai check` validates it.
- `wip/agents` narratives are append-only in the log section. The summary sections at the top are rewritten as understanding improves.

## Language

Plain, direct prose. Say what is, not what might be. Prefer tables for lists of parallel facts. Keep sentences short.
