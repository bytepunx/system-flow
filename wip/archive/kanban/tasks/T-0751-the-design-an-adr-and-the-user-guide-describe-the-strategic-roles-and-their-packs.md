---
id: T-0751
type: task
nature: improvement
title: The design, an ADR, and the user guide describe the strategic roles and their packs
status: done
parent: S-0207
owner: alex
created: 2026-10-03T07:11:26Z
updated: 2026-10-03T07:27:15Z
transitions:
  - to: ready
    at: 2026-10-03T07:11:43Z
    by: agent-S-0207
  - to: in-progress
    at: 2026-10-03T07:21:08Z
    by: agent-S-0207
  - to: done
    at: 2026-10-03T07:27:15Z
    by: agent-S-0207
stream: S-0207
tags: [flai]
touches: [design/system/conventions.md, design/system/agent-context.md, design/system/flai-cli.md, design/system/workflow.md, design/system/metrics.md, design/system/work-hierarchy.md, design/adrs, docs/users/flai.md, docs/users/conventions.md]
after: [T-0746, T-0747]
usage:
  source: log
  seconds: 367
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 70
      output: 21132
      cache_read: 4122860
      cache_write: 83871
      cost: 1.7756
---
# T-0751 The design, an ADR, and the user guide describe the strategic roles and their packs

## Work

- Write an ADR refining ADR-0068 and ADR-0059: the strategic roles, a pack without a story, and the planner's pack for an epic or a story.
- `design/system/conventions.md` (`## Roles` and its table, `## Tooling`), `design/system/flai-cli.md` (`flai prime`), and `design/system/agent-context.md` describe the roles and their packs; `docs/users/flai.md` says how to prime each.
- Give the design the strategic roles need their topics, so that their packs brief it: `planning` on `work-hierarchy.md`, `orchestration` on `workflow.md`, `analysis` on `metrics.md`, as far as each document bears on the role.
- Waits for T-0746 and T-0747: it describes the flags and the roles they settle.

## Done when

- The documents describe what T-0746 and T-0747 built, and the ADR is accepted.
- `flai prime --role orchestrate` and `--role analyze` brief the design their topics now select.

## Notes
