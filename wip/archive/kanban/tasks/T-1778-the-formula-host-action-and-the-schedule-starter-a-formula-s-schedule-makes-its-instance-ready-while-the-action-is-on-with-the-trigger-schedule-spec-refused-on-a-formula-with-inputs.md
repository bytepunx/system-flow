---
id: T-1778
type: task
nature: feature
title: "The formula host action and the schedule starter: a formula's schedule makes its instance ready while the action is on, with the trigger schedule spec, refused on a formula with inputs"
status: done
parent: S-0422
owner: alex
created: 2026-10-16T05:48:01Z
updated: 2026-11-02T03:31:51Z
transitions:
  - to: ready
    at: 2026-11-02T03:22:26Z
    by: system-flow
  - to: in-progress
    at: 2026-11-02T03:22:28Z
    by: system-flow
  - to: done
    at: 2026-11-02T03:31:51Z
    by: agent-S-0422
stream: S-0422
tags: [cli]
touches: [flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go, flai/internal/serve/formula_schedule.go, flai/internal/serve/formula_schedule_test.go, flai/internal/serve/serve.go, flai/internal/formulas/formulas.go, flai/internal/formulas/formulas_test.go, flai/internal/check/check.go, flai/internal/check/formulas_test.go, flai/internal/roles/roles.go]
after: [T-2012]
usage:
  source: log
  seconds: 563
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 119
      output: 34017
      cache_read: 8641757
      cache_write: 182015
      cost: 3.4693
---
# T-1778 The formula host action and the schedule starter: a formula's schedule makes its instance ready while the action is on, with the trigger schedule spec, refused on a formula with inputs

## Work

- Add the `formula` host action to `hostapi.Actions`, off by default, with what it does: `flai serve` starts a formula on its file's schedule. `flai serve enable formula` and the Settings page take it through `knownAction`.
- A formula's front matter may carry `schedule`, a five-field cron expression in UTC or `daily`, read by `formulas.Load` as `planning.schedule` is parsed. `flai check` reports `formulas.file` for a schedule it cannot parse, and for a schedule on a formula that declares inputs, since nobody is there to give them.
- A formula scheduler per served project, `serve/formula_schedule.go`, looks whenever the project's launcher does, as `analysisScheduler` does, on the shared `cadence`. Each time a formula's schedule comes round while the action is on, it makes the instance by running `flai formula run <name> --ready --autocommit` in the project's main checkout, with no inputs, through the flai binary serve runs as (`os.Executable`, as the merge queue's `flaiAccepter` runs `flai accept`), so that `flai serve` starts its agent; the instance's creation is journalled under `formula` with the trigger `schedule <spec>`. A start that is refused, such as a formula whose check fails, is logged at info and not tried again until the schedule next comes round. Times missed between looks come round once; what the scheduler has seen is kept in memory.
- Tests: a fixture formula with a schedule makes a ready instance when its time comes with the action on, nothing with it off, and once for a time missed between looks; the check's refusals.

Second layer: it waits for T-2012, which adds the `--ready` flag the scheduler runs. The creation lives on `app` in `cmd/formula.go`, which `internal/serve` cannot import, so the scheduler runs the command rather than calling it.

## Done when

`go test ./internal/serve/... ./internal/formulas/... ./internal/check/... ./internal/hostapi/... ./cmd/...` passes.

## Notes

Rewritten by the planner on 2026-10-21, on TH-0471's answer 5 and TH-0582's approval. The role starter it held moved to T-2012; the host action moved here from the step runner.

Revised by S-0422's agent on 2026-11-02: it runs `flai formula run --ready` as a command and waits for T-2012, since the creation is in `cmd`, not `internal/formulas`.
