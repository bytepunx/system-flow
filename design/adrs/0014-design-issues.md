---
id: ADR-0014
title: design/issues records recurring friction with counts and cost
status: accepted
date: 2026-09-15
supersedes: []
superseded_by: []
---

# ADR-0014 design/issues records recurring friction with counts and cost

## Context

Friction repeats: a tool at the wrong version, a fixture that git ignores, a folder that vanishes on clone, a hand edit that drifts. Each occurrence costs wall-clock time and is usually solved in place and forgotten. The continuous-improvement convention asks agents to record such items so their frequency and cost become visible and remediation can be prioritised. That needs a home under `design/`, which ADR-0002 allows only by ADR.

## Decision

`design/` gains `issues/`: one markdown file per issue named `I-nnn-slug.md`, plus `summary.md`, a generated table of open issues, and a `README.md`. Issue front matter carries `id`, `title`, `class` (`defect`, `blocker`, `efficiency`, `impression`), `status` (`open`, `closed`), `count`, `cost` (average wall-clock per occurrence as a duration), `first_reported`, `last_reported`, and `updated`. The body has a description and one dated instance per occurrence. Agents record an occurrence when it happens, increment existing issues rather than duplicating, and include the summary in a story's review report when it changed. At epic completion the summary is presented and the operator decides on remediation. `flai issue` commands (S-027) own the front matter and regenerate the summary.

## Consequences

- `design/` now has five documentation types: `adrs`, `system`, `tech`, `conventions`, `issues`.
- Issue files are the only `design/` documents with a counter that changes often; they are still committed like everything else so history shows when friction started and stopped.
- `flai check` validates the folder once S-027 lands; until then the schema is followed by hand.
- Remediation work is ordinary work: a story under the relevant epic, with nature `remediation` or `improvement`, that closes the issue when done.

## Alternatives considered

- A single `design/issues.md` list: cheap, but no per-issue history and a merge conflict magnet once two agents record at once.
- Work items with nature `remediation`: those are the fix, not the record; an issue may be observed many times before anyone decides to fix it.
- A tracker outside the repository: invisible to agents without integrations and not versioned with the code.
