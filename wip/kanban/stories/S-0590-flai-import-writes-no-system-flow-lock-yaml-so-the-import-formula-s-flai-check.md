---
id: S-0590
type: story
nature: remediation
title: flai import writes no system-flow.lock.yaml, so the import formula's flai check --strict step fails on layout.lock
status: in-progress
owner: alex
created: 2026-11-07T09:44:23Z
updated: 2026-11-07T10:04:09Z
transitions:
  - to: ready
    at: 2026-11-07T10:02:40Z
    by: alex
  - to: in-progress
    at: 2026-11-07T10:03:02Z
    by: agent-S-0590
tags: [flai, import, template]
topics: [template, import]
touches: [flai/cmd/import.go, flai/cmd/import_commit.go, flai/cmd/import_test.go, flai/cmd/import_commit_test.go, docs/users/flai-reference.md, docs/users/flai.md, design/system/flai-cli.md, design/system/template.md, design/issues/I-0302-flai-import-writes-no-system-flow-lock-yaml-so-the-import-formula-s-flai-check.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 60
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 20
          output: 149
          cache_read: 513653
          cache_write: 116084
          cost: 0.1577
cost_of_delay:
  inputs:
    time_lost_per_cycle: 10m
    by: flai
    at: 2026-11-07T09:44:23Z
  value: 25
  by: planner-S-0590
  at: 2026-11-07T10:04:09Z
forecast:
  duration: 45m
  delivery: 2026-11-07T10:49:00Z
  basis: flai forecast's 24m (117 s per unit over 40 done medium remediation stories, size 12) raised to 45m for the integration tier a flai/cmd change runs and for the lock-hash choice T-2570 makes and tests.
  by: planner-S-0590
  at: 2026-11-07T10:04:09Z
finalized:
  by: alex
  at: 2026-11-07T10:02:09Z
---
# S-0590 flai import writes no system-flow.lock.yaml, so the import formula's flai check --strict step fails on layout.lock

## Goal

This story remediates [I-0302](../../../design/issues/I-0302-flai-import-writes-no-system-flow-lock-yaml-so-the-import-formula-s-flai-check.md), "flai import writes no system-flow.lock.yaml, so the import formula's flai check --strict step fails on layout.lock". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0302 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0302 is closed with `flai issue close I-0302 --reason` saying what fixed it

## Tasks
- T-2570 flai import writes system-flow.lock.yaml as flai new does and commits it with the import, so the import formula's check finds no layout.lock
- T-2571 The user guide and the design say flai import writes the lock, and a kept file counts as the project's own
- T-2572 Close I-0302 once the import tests and the integration tier pass with the lock written

## Notes

Cost of delay inputs set by flai from I-0302. time_lost_per_cycle 10m: 10m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-07T08:50:28Z, 0 days before this story; under one cycle counts as one).

### Planning

Proposed fix: `runImport` in `flai/cmd/import.go` renders the template and never saves a lock, while `flai new` does (`flai/cmd/new.go`). The import saves the lock from the render's hashes, as `flai new` does, and `commitImportPaths` commits it. The lock records the rendered hash of a file the import kept. So `flai upgrade` treats that file as the project's own, a conflict, and never overwrites it as unchanged baseline. No ADR records "import writes no lock". Only `design/system/template.md` and `docs/users/flai.md` describe it, so this story changes those documents and writes no ADR.

Touches, all files, no folders:

| Touch | Source |
|-------|--------|
| `flai/cmd/import.go` | layout: `runImport`'s render step, where the lock is missing |
| `flai/cmd/import_commit.go` | layout: `commitImportPaths`, the paths `--commit` and `--formula` stage |
| `flai/cmd/import_test.go` | co-change 27% with the two files above; its comment says an import writes no lock |
| `flai/cmd/import_commit_test.go` | co-change 18%; the reproducing test for `--formula` |
| `docs/users/flai-reference.md` | co-change 36%; generated from the import help, which changes |
| `docs/users/flai.md` | co-change 27%; design: says "`import` writes no lock" |
| `design/system/flai-cli.md` | co-change 27%; design: § Import flow |
| `design/system/template.md` | design: says an imported project records its values only at its first upgrade |
| `design/issues/I-0302-…md`, `design/issues/summary.md` | criterion 2, `flai issue close` |

Co-change also listed `flai/cmd/new.go` (45%). Left out, because the fix reuses its `lockVars` without changing it. `flai/cmd/upgrade.go`, `flai/cmd/serve_import.go`, and the config files were left out too: the fix does not change them.

Forecast: flai gave 24m, 117 s per unit of size over 40 done medium-band remediation stories, times size 12. Raised to 45m. The change under `flai/cmd` makes the close-out run the integration and smoke tiers. T-2570 also makes and tests a choice of hash. The size counts neither. Delivery is 45m from the start at 10:04Z.

Cost of delay: 25.00 USD a week, as `flai cod` works it out from the 10m input, kept as flai gave it. I-0302 has one instance, but every import from the board meets it, so the true rate may be higher. The input is the operator's, so I left it alone.

No `after`: the story reads no other story's output.
