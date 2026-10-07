---
id: T-1159
type: task
nature: feature
title: metrics.md defines the turn classes with an ADR, and the CLI design and the user guide describe the report
status: done
parent: S-0293
owner: alex
created: 2026-10-07T09:28:29Z
updated: 2026-10-07T09:48:02Z
transitions:
  - to: ready
    at: 2026-10-07T09:28:45Z
    by: agent-S-0293
  - to: in-progress
    at: 2026-10-07T09:44:27Z
    by: agent-S-0293
  - to: done
    at: 2026-10-07T09:48:02Z
    by: agent-S-0293
stream: S-0293
tags: []
touches: [design/system/metrics.md, design/adrs, design/system/work-hierarchy.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md]
after: [T-1158]
usage:
  source: log
  seconds: 215
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 55
      output: 25949
      cache_read: 3681573
      cache_write: 109221
      cost: 1.9664
---
# T-1159 metrics.md defines the turn classes with an ADR, and the CLI design and the user guide describe the report

## Work

Define the turn classes in `design/system/metrics.md`, the contract with the dashboard, with a new ADR (`flai adr new` from the worktree), and `usage.turns` in `design/system/work-hierarchy.md`. Describe the report in `design/system/flai-cli.md` and `docs/users/flai.md`, and regenerate `docs/users/flai-reference.md` with `make flai-reference` for the changed `flai stats` help. Waits for T-1158's report, which it describes.

## Done when

- `metrics.md` defines a turn and each class, `usage.turns`, and the `turns` keys of `--json` and the text section, linking the ADR.
- `flai-cli.md`, `docs/users/flai.md`, and `work-hierarchy.md` describe it, and `flai-reference.md` is regenerated.

## Notes
