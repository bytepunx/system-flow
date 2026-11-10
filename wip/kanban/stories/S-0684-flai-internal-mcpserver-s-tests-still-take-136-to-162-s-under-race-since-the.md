---
id: S-0684
type: story
nature: improvement
title: flai/internal/mcpserver's tests still take 136 to 162 s under -race, since the server reads the session's story, role, session, reasons, and started-by from the environment and 95 tests must run serially or in env groups before the rest
status: in-progress
owner: alex
created: 2026-11-10T08:37:09Z
updated: 2026-11-10T08:51:59Z
transitions:
  - to: ready
    at: 2026-11-10T08:50:05Z
    by: alex
  - to: in-progress
    at: 2026-11-10T08:50:12Z
    by: agent-S-0684
tags: [flai, go, mcp, testing]
topics: [testing, mcp]
touches: [flai/internal/mcpserver/server.go, flai/internal/mcpserver/shared.go, flai/internal/mcpserver/issues.go, flai/internal/mcpserver/permission.go, flai/internal/mcpserver/subproject.go, flai/internal/mcpserver/versions.go, flai/internal/mcpserver/session_test.go, flai/internal/mcpserver/server_test.go, flai/internal/mcpserver/task_test.go, flai/internal/mcpserver/messages_test.go, flai/internal/mcpserver/orchestrate_test.go, flai/internal/mcpserver/orchestrator_end_test.go, flai/internal/mcpserver/orchestrator_threads_test.go, flai/internal/mcpserver/shared_test.go, flai/internal/mcpserver/work_test.go, flai/internal/mcpserver/stream_test.go, flai/internal/mcpserver/issues_test.go, flai/internal/mcpserver/items_write_test.go, flai/internal/mcpserver/queue_test.go, flai/internal/mcpserver/permission_test.go, flai/internal/mcpserver/subproject_test.go, flai/internal/mcpserver/versions_test.go, flai/internal/mcpserver/outdated_test.go, flai/internal/mcpserver/process_test.go, flai/cmd/mcp.go, flai/cmd/mcp_http.go, flai/cmd/orchestrate_permissions_test.go, flai/cmd/role_test.go, flai/internal/serve/restart_test.go, design/issues/I-0234-flai-internal-mcpserver-s-tests-still-take-136-to-162-s-under-race-since-the.md, design/issues/summary.md, flai/internal/mcpserver/task.go]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: planner
      seconds: 83
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 86
          output: 7327
          cache_read: 5762019
          cache_write: 142702
          cost: 1.9363
cost_of_delay:
  inputs:
    time_lost_per_cycle: 43s
    by: flai
    at: 2026-11-10T08:37:09Z
  value: 1.79
  by: planner-S-0684
  at: 2026-11-10T08:48:35Z
forecast:
  duration: 2h
  delivery: 2026-11-10T17:42:00Z
  basis: "Its own forecast of 2h; 26th in the pull order with an in-progress limit of 5, behind S-0685, S-0241, S-0659, S-0660, S-0661, S-0662, S-0663, S-0665, S-0666, S-0667, S-0668, S-0669, S-0670, S-0671, S-0672, S-0673, S-0674, S-0675, S-0676, S-0677, S-0678, S-0679, S-0680, S-0681, S-0682 and S-0683."
  by: flai
  at: 2026-11-10T08:49:28Z
finalized:
  by: alex
  at: 2026-11-10T08:50:00Z
---
# S-0684 flai/internal/mcpserver's tests still take 136 to 162 s under -race, since the server reads the session's story, role, session, reasons, and started-by from the environment and 95 tests must run serially or in env groups before the rest

## Goal

This story remediates [I-0234](../../../design/issues/I-0234-flai-internal-mcpserver-s-tests-still-take-136-to-162-s-under-race-since-the.md), "flai/internal/mcpserver's tests still take 136 to 162 s under -race, since the server reads the session's story, role, session, reasons, and started-by from the environment and 95 tests must run serially or in env groups before the rest". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0234 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0234 is closed with `flai issue close I-0234 --reason` saying what fixed it

## Tasks
- T-2904 mcpserver.Options carry the session's story, role, session, reasons, and started-by, which flai mcp fills from the environment
- T-2905 mcpserver takes what its tests now set in package variables and the environment, the permission poll, the flai it runs, its version, the releases API and token, and a child's config, through Options
- T-2906 mcpserver's session tests pass the session in Options and run in parallel, with no serial parent and git's identity set once for the package
- T-2910 mcpserver's permission, subproject, versions, and outdated tests pass Options in place of package variables and run in parallel, and together goes
- T-2915 A test keeps mcpserver's tests off the process, the package's -race time is measured against I-0234's, and I-0234 is closed

## Notes

Cost of delay inputs set by flai from I-0234. time_lost_per_cycle 43s: 1m per occurrence × 1 occurrence ÷ 1.4 cycles of 168h (first reported 2026-10-31T09:50:04Z, 9.9 days before this story).

### Planning

The proposed solution is the first of the two remedies the issue's instance names. It also follows D-0443. The session goes into `mcpserver.Options`, which `flai mcp` fills from the environment, so the package's tests stop setting the process and the serial phase goes. The package's time then becomes that of the parallel prime-on-this-repository tests. The second remedy, priming faster under race, is left out: T-2420 measured no gain from a search change and was cancelled.

What is serial today, as counted by the explorer on main at planning: 36 top-level tests, each serial for a variable of its own, plus 6 serial parents that run 54 subtests through `together`. `t.Setenv` is called 50 times in 12 test files. Of those calls, 40 set a session variable. The other 10 set the releases API and its tokens, `FLAI_CONFIG`, and `FLAI_SLOW_REQUEST`. Besides those, git's identity and three package variables are set on the process.

Touches, all files, none a folder:

| Touch | Source |
|-------|--------|
| `flai/internal/mcpserver/server.go`, `shared.go`, `issues.go` | Declared by the issue: its environment reads at `server.go:203-209`, `shared.go:180`, and `issues.go:196` |
| `flai/internal/mcpserver/task.go` | Layout: it builds the environment of the child flai that `task_done` runs, which may need the session from the server's fields. Added in the 08:48:56Z touches edit, logged as by system-flow. T-2904 touches it |
| `flai/cmd/mcp.go`, `flai/cmd/mcp_http.go` | Co-change (`flai touches suggest`: 16 and 18 of 99 commits), and the layout: they build `mcpserver.Options` |
| `flai/internal/mcpserver/permission.go`, `subproject.go`, `versions.go` | Layout: they hold `permissionPoll`, `flaiExecutable`, and the `FLAI_RELEASES_API` read that the serial tests swap |
| `flai/internal/mcpserver/server_test.go` | Co-change (45 of 99) and the layout: it holds `TestMain`, `together`, and `TestAsTheStorysAgent` |
| The other `flai/internal/mcpserver/*_test.go` files named | Layout: each calls `t.Setenv` or `together`, or assigns a package variable |
| `session_test.go`, `process_test.go` | New test files that T-2904 and T-2915 add |
| `flai/cmd/orchestrate_permissions_test.go`, `flai/cmd/role_test.go`, `flai/internal/serve/restart_test.go` | Layout: they set the session variables before `mcpserver.New`, so they must pass `Options` once the server stops reading the environment |
| The `design/issues/I-0234-…` file and `design/issues/summary.md` | The goal: `flai issue close` writes both |

`folder.go` is left out, though `flai touches suggest` ranks it first: 63 percent co-change on the first plan, and 76 of the 1044 commits on the second. It passes its `opt` to `newServer` already and reads no session variable. `docs/users/flai.md` and `design/system/flai-cli.md` are left out too, since what `flai mcp` does does not change.

Forecast: `flai forecast` gave 56m from size 33 on the first plan and 58m from size 34 on the second. It is raised to 2h. The work converts about 90 tests across 16 files, and each `-race` run of the package takes 2 to 3 minutes: at least two runs in each of T-2906, T-2910, and T-2915, plus the close-out's integration tier. flai's replan after S-0636's acceptance set the delivery and basis from the 2h. The story started at 08:50Z, so it is due about 10:50Z.

Cost of delay: the value of 1.79 USD a week is `flai cod`'s, from the operator's input of 43 s lost per cycle at 150 USD an hour. It stands as computed on both plans. The input is the operator's, set by flai from I-0234's one occurrence of 1m, and this plan changes no input.

`after`: none. The story reads no other story's output.
