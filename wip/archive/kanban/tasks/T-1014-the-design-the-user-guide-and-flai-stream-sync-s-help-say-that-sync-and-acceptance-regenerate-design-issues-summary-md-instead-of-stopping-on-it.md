---
id: T-1014
type: task
nature: improvement
title: The design, the user guide, and flai stream sync's help say that sync and acceptance regenerate design/issues/summary.md instead of stopping on it
status: done
parent: S-0278
owner: alex
created: 2026-10-06T11:32:30Z
updated: 2026-10-06T20:03:26Z
transitions:
  - to: ready
    at: 2026-10-06T19:59:54Z
    by: agent-S-0278
  - to: in-progress
    at: 2026-10-06T19:59:54Z
    by: agent-S-0278
  - to: done
    at: 2026-10-06T20:03:26Z
    by: agent-S-0278
stream: S-0278
tags: [flai, docs]
touches: [flai/cmd/stream.go, design/system/flai-cli.md, design/system/continuous-improvement.md, docs/users/flai.md, docs/users/flai-reference.md]
after: [T-1012, T-1013]
usage:
  source: log
  seconds: 212
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 59
      output: 317
      cache_read: 1736971
      cache_write: 69946
      cost: 0.7854
---
# T-1014 The design, the user guide, and flai stream sync's help say that sync and acceptance regenerate design/issues/summary.md instead of stopping on it

## Work

Describe the behaviour T-1012 and T-1013 built, and link T-1011's ADR:

- `design/system/flai-cli.md`: the `flai stream sync` row. When a rebase stops only on generated files, it regenerates them and continues. The trial merge leaves them out of a pair's conflicts. Also say where `flai accept` is described, if it is.
- `design/system/continuous-improvement.md`, `## Summary`: the file is generated, a sync or acceptance regenerates it rather than merging it, and the ADR is linked.
- `docs/users/flai.md`: the section on story branches and `flai stream sync`. Change the conflict example and the resolution steps where they still tell the agent to resolve `summary.md` by hand.
- `flai/cmd/stream.go`: the help of `flai stream sync`. Regenerate `docs/users/flai-reference.md` with `make flai-reference` (`scripts/flai-reference.sh`). It also rewrites the flag index in `docs/operators/settings.md`, which should not change, since no flag is added.

This task waits for T-1012 and T-1013, whose behaviour it describes. It shares no path with T-1015, so the two run together.

## Done when

- The design, the user guide, and the help describe the regeneration and the trial merge as built, and link the ADR.
- `docs/users/flai-reference.md` is regenerated from the help.
- The markdown lint and `flai check --strict` pass.

## Notes
