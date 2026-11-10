---
id: T-2794
type: task
nature: improvement
title: Split flai/internal/hostapi/writes.go by responsibility, cutting itemSpecs in three
status: done
parent: S-0631
owner: alex
created: 2026-11-09T22:04:19Z
updated: 2026-11-10T06:38:09Z
transitions:
  - to: ready
    at: 2026-11-10T06:16:41Z
    by: agent-S-0631
  - to: in-progress
    at: 2026-11-10T06:16:42Z
    by: agent-S-0631
  - to: done
    at: 2026-11-10T06:38:09Z
    by: agent-S-0631
stream: S-0631
tags: [cli]
touches: [flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_run.go, flai/internal/hostapi/writes_actions.go, flai/internal/hostapi/writes_spec.go, flai/internal/hostapi/writes_items.go, flai/internal/hostapi/writes_streams.go, flai/internal/hostapi/writes_host.go, flai/internal/hostapi/writes_describe.go, flai/internal/hostapi/writes_methods.go, flai/internal/hostapi/writes_outcome.go, flai/internal/hostapi/settings.go]
usage:
  source: log
  seconds: 1287
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 42
      output: 8483
      cache_read: 6041413
      cache_write: 180879
      cost: 2.4075
---
# T-2794 Split flai/internal/hostapi/writes.go by responsibility, cutting itemSpecs in three

## Work

`writes.go` is 3,626 lines. Move its declarations into new files of package `hostapi` by responsibility. What any of them does stays the same. The line ranges below are approximate, from the plan's read on 2026-11-09:

| Lines | Declarations | New file |
|-------|--------------|----------|
| 54–290 | `Run`, `Runner`, `ExecRunner`, `runFlai`, `Written`, the request journal (`Request`, `Requests`, `begin`, `finish`, `forgetOld`, `repeated`, `requestID`) | `writes_run.go` |
| 291–627 | the `Action*` constants, `Actions`, `DashboardSees`, `DefinedActions`, `roleAction`, `Host`, `Entry`, `Host.enabled`, `EnableCommand` | `writes_actions.go` |
| 628–988 | `spec`, `owner`, `forgeAcceptance`, the argument checks, `decode`, `specs()`, `asksEverywhere` | `writes_spec.go` |
| 989–2464 | `itemSpecs` | see below |
| 2466–3062 | `orchestrateNow`, the `describe*` helpers, `roleRunArgs`, `testArgs`, `docDir`, `verifyArgs`, `withJSON` | `writes_describe.go` |
| 3063–3386 | `writeMethods`, `systemGit`, `queueSpecs` and the `queue*` helpers, `describeQueue`, `methodsFrom` | `writes_methods.go` |
| 3387–3626 | `writing`, `WritesRunning`, `Host.record`, `Publisher`, `describe`, `outcome`, `agentArgs`, `roleArgs` | `writes_outcome.go` |

`itemSpecs` is one function of 1,476 lines. Cut it into three functions, each returning its part of the table, and have `itemSpecs` join them in the same order:

| Specs | New file |
|-------|----------|
| item, story, board, task, issue, accept.run (about lines 989–1749) | `writes_items.go` |
| stream, thread, doc, file, adr (about 1750–2137) | `writes_streams.go` |
| publish, dashboard, host, checks, test, verify, agent, plan, analyze, curate, role, orchestrate (about 2138–2464) | `writes_host.go` |

Its closures `one`, `read`, and `moveArgs` (lines 990–1033) move to package scope, or are passed to the three functions.

Keep `writes.go` with its imports, error codes, `Trailer`, and `itemSpecs`. `flai/internal/verify/select_test.go:554` stats it. Follow the precedent of `componentSpecs` in `components.go` and `settingsSpecs` in `settings.go`, which already live in their feature's file.

Fix the comment at `flai/internal/hostapi/settings.go:615` that names `writes.go` if what it points to moved.

Change no test file. `writes_test.go` stays as it is.

This task waits for nothing and shares no path with the other splits, so all of them run together.

## Done when

- No non-test file in `flai/internal/hostapi` is over 2,000 lines, and none of the new files is over 1,200.
- The method table `MethodsFor` builds is the same as before, in the same order: the package's tests pass unchanged.
- `git diff --stat` shows no `_test.go` file changed.
- `flai test flai/internal/hostapi` passes, gofmt, goimports, and golangci-lint included.

## Notes
