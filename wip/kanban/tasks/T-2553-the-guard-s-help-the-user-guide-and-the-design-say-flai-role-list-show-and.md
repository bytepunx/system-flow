---
id: T-2553
type: task
nature: improvement
title: The guard's help, the user guide, and the design say flai role list, show, and definition are reads
status: done
parent: S-0568
owner: alex
created: 2026-11-07T08:19:09Z
updated: 2026-11-07T08:35:16Z
transitions:
  - to: ready
    at: 2026-11-07T08:29:49Z
    by: system-flow
  - to: in-progress
    at: 2026-11-07T08:29:50Z
    by: system-flow
  - to: done
    at: 2026-11-07T08:35:16Z
    by: agent-S-0568
stream: S-0568
tags: [flai, guard, docs]
touches: [flai/cmd/guard.go, docs/users/flai-reference.md, docs/users/flai.md, design/system/flai-cli.md, design/system/strategic-agents.md]
after: [T-2551]
usage:
  source: log
  seconds: 318
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 38
      output: 293
      cache_read: 4480347
      cache_write: 29941
      cost: 1.6147
---
# T-2553 The guard's help, the user guide, and the design say flai role list, show, and definition are reads

## Work

Layer 2: this task waits for T-2551, because it states the forms that task settles and the guard's tests fix.

- Add the reads to the list in the `Long` text of `flai guard` in `flai/cmd/guard.go`, then regenerate `docs/users/flai-reference.md` with `make flai-reference`.
- In `docs/users/flai.md`, say under `## Sub-agents` that a sub-agent may read the roles with `flai role list`, `flai role show`, and `flai role definition`, and that `flai role run` and `flai role definition --write` are refused.
- In `design/system/flai-cli.md`, add the reads to the `flai guard` row's list, naming S-0568 and I-0245.
- In `design/system/strategic-agents.md`, say under `## Role files (S-0414)`, in its Commands item, what the guard lets every agent run, as `## Formulas (S-0421)` says for `flai formula list` and `show`.

## Done when

- `flai guard --help` names `role list`, `role show`, and `role definition` without `--write` among the reads, and `docs/users/flai-reference.md` matches what `make flai-reference` generates.
- `docs/users/flai.md`, `design/system/flai-cli.md`, and `design/system/strategic-agents.md` each state the reads and the two refusals, and each file's `updated` date is bumped where it carries one.
- `flai test` on the five files passes.

## Notes
