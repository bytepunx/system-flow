---
id: T-2790
type: task
nature: improvement
title: flai release --evaluate --judged meets the judgement policy with the orchestrator's reason
status: done
parent: S-0625
owner: alex
created: 2026-11-09T22:01:59Z
updated: 2026-11-10T05:44:00Z
transitions:
  - to: ready
    at: 2026-11-10T05:19:13Z
    by: system-flow
  - to: in-progress
    at: 2026-11-10T05:19:14Z
    by: system-flow
  - to: done
    at: 2026-11-10T05:44:00Z
    by: agent-S-0625
stream: S-0625
tags: [flai, cli]
touches: [flai/cmd/release.go, flai/internal/release/evaluate.go, flai/internal/release/evaluate_test.go, flai/cmd/release_evaluate_test.go, docs/users/flai-reference.md, docs/operators/settings.md]
after: [T-2788]
usage:
  source: log
  seconds: 1486
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 35
      output: 5933
      cache_read: 4904340
      cache_write: 123047
      cost: 1.8106
---
# T-2790 flai release --evaluate --judged meets the judgement policy with the orchestrator's reason

## Work

Add `--judged "<reason>"` to `flai release --evaluate`, as T-2788's ADR decides. It waits for T-2788 because it builds the flag's name and rules that the ADR fixes.

- `flai/cmd/release.go`: declare the flag beside `--require-met`. Refuse it without `--evaluate`, as `--require-met` is refused. Say in the help text that under `judgement` a reason makes the evaluation met.
- `flai/internal/release/evaluate.go`: under `judgement`, a non-empty reason sets `Met` when something is accepted and not yet released. The evaluation's reason names the judgement. `holdByEpic` still holds the batch back under `whole_epics`. Under `threshold` and `theme` the reason changes nothing.
- Tests that reproduce I-0291 at the command:
  - Under `judgement`, `--require-met` with no reason exits 2. This is the cause.
  - With `--judged` and a reason, it exits 0 and prints the judgement.
  - It still exits 2 with nothing to release, with `whole_epics` holding the batch, or with an empty reason.
  - Under `threshold`, `--judged` does not make an unmet evaluation met.
  - Put these in `TestReleaseEvaluateRequireMet` in `flai/cmd/release_evaluate_test.go` and `TestEvaluateJudgement` in `flai/internal/release/evaluate_test.go`.
- Regenerate the reference with `make flai-reference`. It writes `docs/users/flai-reference.md` and the flag index in `docs/operators/settings.md`.

## Done when

- `flai release --evaluate --require-met --judged "<reason>"` exits 0 under `judgement` with a batch to release, and exits 2 in each case listed above.
- The tests pass with `flai test` on the changed paths, and `TestReferenceIsCurrent` passes.

## Notes

Layer 2. It runs beside T-2791; the two share no path.
