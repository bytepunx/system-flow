---
id: ADR-0013
title: design/conventions holds agent norms, one file per topic
status: accepted
date: 2026-09-15
supersedes: []
superseded_by: []
---

# ADR-0013 design/conventions holds agent norms, one file per topic

## Context

Every project run with coding agents starts by re-establishing how the agent should work: how to communicate, when to ask, how to use the board, when a decision needs a record, how to commit. Those norms were partly inline in the baseline `CLAUDE.md`, partly implied by `design/system`, and partly re-taught in conversation. They must ship with the template so a new project has them from the first session, be loaded by every agent session before any change, and be extendable per project without forking the baseline. ADR-0002 fixes three top-level documentation folders and allows new subfolders of `design/` only by ADR. This is that ADR.

## Decision

`design/` gains a fourth documentation type: `design/conventions/`, holding how agents work in this repository. It contains one markdown file per topic area plus a `README.md` index in read order. Files are short, imperative rules with front matter (`title`, `updated`, `audience: agent`, `order`). Each file ends with a marker line, `<!-- system-flow:end-of-baseline -->`, below which a project adds its own rules; `flai upgrade` replaces only the section above the marker. The folder resolves through `layout.design` in `system-flow.yaml`; no new layout key is needed. The template's `CLAUDE.md` instructs every session to read the conventions first. Precedence is: an explicit user instruction, then project additions, then baseline conventions, then the agent's own defaults.

## Consequences

- Norms are written once in the template and inherited by every project; changing a norm is a template release, not a conversation per project.
- The three top-level folders of ADR-0002 stand. `design/` now has `adrs`, `system`, `tech`, and `conventions`, each a documentation type with its own lifecycle.
- `CLAUDE.md` shrinks to a map and a priming instruction. Norms that were inline there move to `design/conventions/`.
- `flai check` validates the folder, `flai prime` prints it in read order, and the dashboard indexes it with the rest of `design/`.
- The conventions are agent-tool agnostic. `CLAUDE.md` is the Claude Code entry point; an `AGENTS.md` or similar pointer can be added for other tools without changing the folder.
- Existing projects need `flai upgrade` (S-020) or a manual copy to gain the folder; until then `flai check` reports it missing.

## Alternatives considered

- A fourth top-level folder `conventions/`: gives the agent audience its own root, but adds a layout key, a manifest migration for existing projects, and a fourth thing to explain, for content that is internal documentation by any definition.
- Everything inline in `CLAUDE.md`: one long file that every agent tool reads differently, no per-topic ownership, no marker-based upgrades per topic.
- A single `design/CONVENTIONS.md`: simpler, but topics evolve at different rates and a project extension of one topic would force a merge of the whole file.
