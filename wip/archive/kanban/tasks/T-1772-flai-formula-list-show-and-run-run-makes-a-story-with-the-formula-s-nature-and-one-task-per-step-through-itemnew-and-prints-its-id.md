---
id: T-1772
type: task
nature: feature
title: "flai formula list, show, and run: run makes a story with the formula's nature and one task per step, through itemnew, and prints its ID"
status: done
parent: S-0421
owner: alex
created: 2026-10-16T05:47:27Z
updated: 2026-11-02T02:40:22Z
transitions:
  - to: ready
    at: 2026-11-02T02:20:52Z
    by: system-flow
  - to: in-progress
    at: 2026-11-02T02:20:53Z
    by: system-flow
  - to: done
    at: 2026-11-02T02:40:22Z
    by: agent-S-0421
stream: S-0421
tags: [cli]
touches: [flai/cmd/formula.go, flai/cmd/formula_test.go, flai/cmd/root.go, docs/users/flai-reference.md, docs/operators/settings.md, flai/internal/workitem/create.go, flai/internal/guard/guard.go, flai/internal/guard/formula_test.go]
after: [T-1769, T-1770]
usage:
  source: log
  seconds: 1169
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 109
      output: 40680
      cache_read: 6115201
      cache_write: 191539
      cost: 3.1798
---
# T-1772 flai formula list, show, and run: run makes a story with the formula's nature and one task per step, through itemnew, and prints its ID

## Work

- Add `flai formula list [--json]` and `flai formula show <name> [--json]`.
- Add `flai formula run <name> [--input key=value]... [--epic E-nnnn] [--json]`. It refuses a missing or undeclared input before writing anything. Through `itemnew`, it creates:
  - the story, with the formula's nature, `formula: <name>`, a goal from the formula's description, and one criterion per step;
  - one task per step, with `## Work` naming what the step runs with its inputs filled, `## Done when` its exit, and the steps' `after` as the tasks' `after`.
- It prints the story's ID.
- Regenerate `docs/users/flai-reference.md`.

Waits for the loader and the front matter field.

## Done when

`go test ./cmd/...` passes. A fixture formula's instance passes `flai check --strict`.

## Notes
