---
id: T-3037
type: task
nature: improvement
title: Set this project's analysis.schedule and analysis.formula analysis-cadence through flai manifest set
status: done
parent: S-0692
owner: alex
created: 2026-11-10T10:33:56Z
updated: 2026-11-11T02:26:05Z
transitions:
  - to: ready
    at: 2026-11-11T02:25:46Z
    by: system-flow
  - to: in-progress
    at: 2026-11-11T02:25:47Z
    by: system-flow
  - to: done
    at: 2026-11-11T02:26:05Z
    by: agent-S-0692
stream: S-0692
tags: [manifest, analyzer]
touches: [system-flow.yaml]
after: [T-3031]
usage:
  source: log
  seconds: 18
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 7
      output: 942
      cache_read: 907704
      cache_write: 23156
      cost: 0.339
---
# T-3037 Set this project's analysis.schedule and analysis.formula analysis-cadence through flai manifest set

## Work

Waits for T-3031: only the tree's flai, built from the branch with `scripts/flai.sh`, can write `analysis.formula`; the installed flai refuses the key.

- In the story's worktree, run `scripts/flai.sh manifest set analysis.schedule=<spec> analysis.formula=analysis-cadence`, with the cadence the operator chose on the story's plan thread (recommended there: `0 6 * * 1`, Mondays at 06:00 UTC). Do not edit `system-flow.yaml` by hand.
- Run `scripts/flai.sh check --strict`: `manifest.analysis` must report nothing, since `analysis-cadence`'s one input, `focus`, defaults to empty.
- Say in the story's Notes that the schedule makes the formula's instance only while the `analyze` host action is on for this project (`flai serve actions`), as `design/system/strategic-agents.md` § The analysis schedule says, and that turning it on is the operator's, on the host.

## Done when

- `system-flow.yaml`'s `analysis` block names the chosen `schedule` and `formula: analysis-cadence`, written by `flai manifest set`, and `flai check --strict` passes.

## Notes
