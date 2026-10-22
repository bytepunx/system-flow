---
id: T-2034
type: task
nature: improvement
title: The guard's help, the generated reference, the user docs, and the design say a sub-agent may list and serve components
status: done
parent: S-0477
owner: alex
created: 2026-10-22T20:58:42Z
updated: 2026-10-22T21:13:59Z
transitions:
  - to: ready
    at: 2026-10-22T21:13:39Z
    by: agent-S-0477
  - to: in-progress
    at: 2026-10-22T21:13:40Z
    by: agent-S-0477
  - to: done
    at: 2026-10-22T21:13:59Z
    by: agent-S-0477
stream: S-0477
tags: [flai, guard, components, docs]
touches: [flai/cmd/guard.go, docs/users/flai-reference.md, docs/users/flai.md, design/system/agent-adapters.md, design/system/component-library.md]
after: [T-2032]
usage:
  source: log
  seconds: 19
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 5
      output: 1341
      cache_read: 316012
      cache_write: 9158
      cost: 0.1545
---
# T-2034 The guard's help, the generated reference, the user docs, and the design say a sub-agent may list and serve components

## Work

Waits for T-2032, since the help text describes the read list it changes and the reference is generated from the help.

- `flai/cmd/guard.go`: in the `Long` help, add `components list and serve` to the list of flai commands a sub-agent may run, in its alphabetical place after `cod`.
- `docs/users/flai-reference.md`: regenerate with `make flai-reference` (`scripts/flai-reference.sh`), which also rewrites the flag index in `docs/operators/settings.md`; commit only what changed.
- `docs/users/flai.md` § Sub-agents: in the paragraph that says what a sub-agent may read and run, name `flai components list` and `flai components serve`, saying `serve` runs the sub-project's own dev server on 127.0.0.1 and `config` stays refused (S-0477, I-0165).
- `design/system/agent-adapters.md` § The guard: in the sentence that lists what the guard refuses a sub-agent, note that `components list` and `serve` are reads since S-0477, and why.
- `design/system/component-library.md` § The dev server: add that a story's sub-agent may run `flai components serve`, so a headless pass of the harness need not fall to the story's agent (I-0165).

Bump `updated` in each design document. Run `flai test` on the files changed for the markdown lint.

## Done when

- `flai guard --help` names `components list and serve` among the reads, and `docs/users/flai-reference.md` matches the help, so `TestReference` in `flai/cmd/reference_test.go` passes
- `docs/users/flai.md`, `design/system/agent-adapters.md`, and `design/system/component-library.md` say a sub-agent may list and serve components and that `config` is refused
- `flai test` on the changed files passes, the markdown lint included

## Notes

`design/system/flai-cli.md` is left alone: its `flai components` row does not mention the guard, and the guard has no row of its own there.
