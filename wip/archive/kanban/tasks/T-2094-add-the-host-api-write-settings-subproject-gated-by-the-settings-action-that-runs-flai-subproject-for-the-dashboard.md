---
id: T-2094
type: task
nature: improvement
title: Add the host API write settings.subproject, gated by the settings action, that runs flai subproject for the dashboard
status: done
parent: S-0491
owner: alex
created: 2026-10-23T20:39:34Z
updated: 2026-10-23T21:47:28Z
transitions:
  - to: ready
    at: 2026-10-23T21:27:04Z
    by: agent-S-0491
  - to: in-progress
    at: 2026-10-23T21:27:05Z
    by: agent-S-0491
  - to: done
    at: 2026-10-23T21:47:28Z
    by: agent-S-0491
stream: S-0491
tags: [cli, dashboard]
touches: [flai/internal/hostapi/settings.go, flai/internal/hostapi/writes_test.go, flaiover/src/lib/server/agent.ts, design/system/dashboard-host-channel.md]
after: [T-2093]
usage:
  source: log
  seconds: 1223
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 67
      output: 27484
      cache_read: 4114677
      cache_write: 160122
      cost: 2.2521
---
# T-2094 Add the host API write settings.subproject, gated by the settings action, that runs flai subproject for the dashboard

## Work

Add the write method `settings.subproject` to `settingsSpecs` in `flai/internal/hostapi/settings.go`, beside `settings.shared` and `settings.manifest`, under the `settings` host action for the project (not `hostwide`). Its params are `op` (`add`, `set`, or `remove`), `name`, and for `add` and `set` the optional `path`, `kind`, `tags`, and `new_name`; none may be named `project`, which is the channel's key (`params_test.go`, I-0163). It builds the command line `subproject <op> ... --autocommit --trailer=<Trailer> --json`, as `manifestEdit` and `sharedEdit` do, so the folder is created and the change committed by `flai subproject` from T-2093, and journals it with `sayArgs`. A missing `op` or `name`, or an `op` it does not know, is refused with each field named in `Data`, as `manifestEdit` refuses.

The dashboard reads the sub-projects from `project.info`, which already answers `projects` with name, path, kind, and tags; extend it with `components` and `install` present-or-not only if the panel in the later task needs them.

Add the method to `REQUIRED_METHODS` in `flaiover/src/lib/server/agent.ts` with a comment naming this story, since `contract_test.go` requires the two lists to match; add a `good` and a `refused` case in `writes_test.go`, which `TestEveryWriteIsCovered` demands; and record the method in `design/system/dashboard-host-channel.md` § Decision, in the numbered list that records each method by story.

It waits for T-2093, whose command it runs.

## Done when

- `settings.subproject` with `op: add`, a name, and a kind runs `flai subproject add` with `--autocommit`, `--trailer`, and `--json`, and the method is refused without the `settings` host action.
- A call missing `name`, or with an `op` not among `add`, `set`, and `remove`, is refused with the field named, shown by `writes_test.go`.
- `go test ./internal/hostapi/...` passes, `TestEveryWriteIsCovered` and `contract_test.go` among them, and `flaiover/src/lib/server/agent.ts` lists the method.
- `design/system/dashboard-host-channel.md` names the method, its params, and its action.

## Notes

Precedents: `manifestEdit` (`settings.go` § `manifestEdit`), `sharedEdit`, and the components methods, which take `sub_project` rather than `project`.
