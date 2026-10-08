---
id: S-0319
type: story
nature: improvement
title: flai test never runs flaiover's prettier, eslint, or svelte-check, so formatting faults surface only at the close-out
status: done
owner: alex
created: 2026-10-07T18:59:50Z
updated: 2026-10-08T06:13:49Z
transitions:
  - to: ready
    at: 2026-10-08T00:16:14Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-08T05:51:57Z
    by: agent-S-0319
  - to: review
    at: 2026-10-08T06:12:42Z
    by: agent-S-0319
  - to: done
    at: 2026-10-08T06:13:49Z
    by: orchestrator
tags: []
topics: [testing]
touches: [scripts/flaiover-lint.sh, system-flow.yaml, flai/internal/verify/select_test.go, design/system/devex.md, design/conventions/tooling.md, design/issues/I-0105-flai-test-never-runs-flaiover-s-prettier-eslint-or-svelte-check-so-formatting-faults-surface-only-at-the-close-out.md, design/issues/summary.md, design/issues/I-0086-the-close-out-s-install-smoke-test-failed-once-and-passed-when-run-alone-with-no-cause-in-its-output.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1267
  turns:
    - day: 2026-10-08
      ceremony: 1
      test_runs: 3
      hand_edits: 1
      work: 44
  models:
    - model: claude-opus-5-5
      input: 100
      output: 19175
      cache_read: 7025696
      cache_write: 181511
      cost: 3.2411
  strategic:
    - kind: orchestrator
      seconds: 331
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 40
          output: 503
          cache_read: 9314484
          cache_write: 40442
          cost: 2.3049
        - model: claude-sonnet-5-5
          input: 8
          output: 56
          cache_read: 82288
          cache_write: 34933
          cost: 0.1011
cost_of_delay:
  inputs:
    time_lost_per_cycle: 10m
    by: flai
    at: 2026-10-07T18:59:50Z
  value: 25
  by: planner-S-0319
  at: 2026-10-08T00:15:45Z
forecast:
  duration: 30m
  delivery: 2026-10-08T05:35:00Z
  basis: "Its own forecast of 30m; 1st in the pull order with an in-progress limit of 3, behind S-0232, S-0318 and S-0320."
  by: flai
  at: 2026-10-08T05:02:40Z
finalized:
  by: orchestrator
  at: 2026-10-08T00:16:10Z
---
# S-0319 flai test never runs flaiover's prettier, eslint, or svelte-check, so formatting faults surface only at the close-out

## Goal

This story remediates [I-0105](../../../design/issues/I-0105-flai-test-never-runs-flaiover-s-prettier-eslint-or-svelte-check-so-formatting-faults-surface-only-at-the-close-out.md), "flai test never runs flaiover's prettier, eslint, or svelte-check, so formatting faults surface only at the close-out". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0105 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0105 is closed with `flai issue close I-0105 --reason` saying what fixed it

## Tasks
- T-1304 scripts/flaiover-lint.sh runs prettier --check and eslint on the flaiover files it is given
- T-1305 A flaiover-lint tier runs on changed flaiover files, and a test shows flai test selects it
- T-1306 The tier lists in devex.md and tooling.md name flaiover-lint, and I-0105 is closed

## Notes

Cost of delay inputs set by flai from I-0105. time_lost_per_cycle 10m: 10m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-07T07:23:32Z, 0.5 days before this story; under one cycle counts as one).

### Planning

Proposed fix, from the two instances (S-0212, S-0329): in both, task sub-agents left prettier faults after `flai test` ran only vitest on their flaiover files. Add a `flaiover-lint` tier to `system-flow.yaml`'s `tests`, not `all_only`, that runs a new `scripts/flaiover-lint.sh` with prettier `--check` and eslint on the changed flaiover files. svelte-check takes no file list and checks the whole project, so it stays in the `all_only` `flaiover` tier and at the close-out.

Touches, all files, no folder touch kept:

| Touch | Source | Why |
|-------|--------|-----|
| `scripts/flaiover-lint.sh` | layout | The new tier command, beside `scripts/flaiover-unit.sh` (T-1304) |
| `system-flow.yaml` | design | `tests` is where tiers are declared (`design/system/project-manifest.md`), written with `flai manifest set` (T-1305) |
| `flai/internal/verify/select_test.go` | layout | Tier selection is `flai/internal/verify/select.go`; the reproducing test sits beside it (T-1305) |
| `design/system/devex.md` | design | Its `Tier commands`, `Test tiers`, and `Close-out` rows list this repository's tiers (T-1306) |
| `design/conventions/tooling.md` | design | Its project addition lists what `flai test` runs (T-1306) |
| `design/issues/I-0105-flai-test-never-runs-flaiover-s-prettier-eslint-or-svelte-check-so-formatting-faults-surface-only-at-the-close-out.md` | goal | Criterion 2 closes it (T-1306) |
| `design/issues/summary.md` | co-change | `flai issue close` rewrites it (T-1306) |

`flai touches suggest` listed only paths that change with most stories, such as other issues, `design/system/flai-cli.md`, and `docs/users/flai.md`. None bears on this fix, so none was added.

Forecast: 30m, delivery 2026-10-08T07:15:00Z. `flai forecast` gave 13m: 83 s per unit of size times size 9. Raised because each task proves the tier by planting and reverting a prettier fault, and the first run installs flaiover's dependencies in the fresh worktree. Delivery is flai's 06:58Z, 22nd in the pull order, plus the 17m added.

Cost of delay value: 25 USD a week, as `flai cod` gives it from the inputs: 10m lost per 168h cycle at 150 USD an hour. It stands as computed. I-0105 has since recorded a second instance (S-0329), so the input may be low. Changing it is the operator's call, raised in the plan's thread.

### Accepted by the orchestrator

- Verified: 488d143780365c99dbd5d1d953792113fe05120a
- At: 2026-10-08T06:13:49Z

Verdict: accept. flai verify passed every step at the branch head 488d1437, and the verifier matched both criteria to the diff. Minor gap left for a follow-up: scripts/README.md does not list scripts/flaiover-lint.sh or name the flaiover-lint tier in its intro (main already leaves install-test.sh unlisted too).
- 1: scripts/flaiover-lint.sh, system-flow.yaml, flai/internal/verify/select_test.go, design/system/devex.md, design/conventions/tooling.md
- 2: design/issues/I-0105-flai-test-never-runs-flaiover-s-prettier-eslint-or-svelte-check-so-formatting-faults-surface-only-at-the-close-out.md, design/issues/summary.md
