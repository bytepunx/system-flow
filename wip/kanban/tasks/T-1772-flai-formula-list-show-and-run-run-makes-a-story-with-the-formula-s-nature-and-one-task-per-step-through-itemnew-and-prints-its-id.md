---
id: T-1772
type: task
nature: feature
title: "flai formula list, show, and run: run makes a story with the formula's nature and one task per step, through itemnew, and prints its ID"
status: backlog
parent: S-0421
owner: alex
created: 2026-10-16T05:47:27Z
updated: 2026-10-21T12:38:26Z
transitions: []
stream: S-0421
tags: [cli]
touches: [flai/cmd/formula.go, flai/cmd/formula_test.go, flai/cmd/root.go, docs/users/flai-reference.md, docs/operators/settings.md]
after: [T-1769, T-1770]
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
