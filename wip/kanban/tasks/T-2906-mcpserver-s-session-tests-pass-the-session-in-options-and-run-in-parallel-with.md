---
id: T-2906
type: task
nature: improvement
title: mcpserver's session tests pass the session in Options and run in parallel, with no serial parent and git's identity set once for the package
status: in-progress
parent: S-0684
owner: alex
created: 2026-11-10T08:47:33Z
updated: 2026-11-10T09:07:33Z
transitions:
  - to: ready
    at: 2026-11-10T09:07:32Z
    by: system-flow
  - to: in-progress
    at: 2026-11-10T09:07:33Z
    by: system-flow
stream: S-0684
tags: [flai, go, mcp, testing]
touches: [flai/internal/mcpserver/server_test.go, flai/internal/mcpserver/task_test.go, flai/internal/mcpserver/messages_test.go, flai/internal/mcpserver/orchestrate_test.go, flai/internal/mcpserver/orchestrator_end_test.go, flai/internal/mcpserver/orchestrator_threads_test.go, flai/internal/mcpserver/shared_test.go, flai/internal/mcpserver/work_test.go, flai/internal/mcpserver/stream_test.go, flai/internal/mcpserver/issues_test.go, flai/internal/mcpserver/items_write_test.go, flai/internal/mcpserver/queue_test.go]
after: [T-2904]
---
# T-2906 mcpserver's session tests pass the session in Options and run in parallel, with no serial parent and git's identity set once for the package

## Work

This task waits for the first task, whose `Options` fields it passes. It changes test files only, none of which the second task touches, so it can run beside the second task.

Today the tests that need a session variable set it with `t.Setenv` in a serial top-level test. The ones that share a value run as `testX` subtests through `together`, as `TestInAGitSession` and `TestAsTheStorysAgent` do (D-0443). For every such test in the files this task touches:

- Make it a top-level `TestX` again that calls `t.Parallel()`. Pass `Story`, `Role`, `Session`, `Reasons`, or `StartedBy` in the `Options` its server is made with, through the package's test helpers, in place of `t.Setenv`. Replace each `inEnv` check with what the server was given.
- Drop the serial parents that only set a session variable.
- Set git's isolation and identity once in `TestMain`, in `server_test.go`, before `m.Run`, or in each test repository's own git config. `gittest.Identity` calls `t.Setenv`, which keeps `TestInAGitSession` and `issues_test.go`'s test serial.
- In `server_test.go`, the test that sets `perf.SlowEnv` passes `Options.Slow` instead.

Leave `together` in place while `permission_test.go` still calls it. The next task removes it.

## Done when

- [ ] No file this task touches calls `t.Setenv`, `gittest.Identity`, or `together`. `TestMain` may set the process environment once, before any test runs.
- [ ] Every top-level test in those files calls `t.Parallel()`.
- [ ] `go test -race -count=1 ./internal/mcpserver` passes from `flai/`, and so does `flai test` on the package.

## Notes
