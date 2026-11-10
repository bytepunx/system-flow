---
id: T-2904
type: task
nature: improvement
title: mcpserver.Options carry the session's story, role, session, reasons, and started-by, which flai mcp fills from the environment
status: done
parent: S-0684
owner: alex
created: 2026-11-10T08:46:58Z
updated: 2026-11-10T09:06:48Z
transitions:
  - to: ready
    at: 2026-11-10T08:51:12Z
    by: system-flow
  - to: in-progress
    at: 2026-11-10T08:51:12Z
    by: system-flow
  - to: done
    at: 2026-11-10T09:06:48Z
    by: agent-S-0684
stream: S-0684
tags: [flai, go, mcp]
touches: [flai/internal/mcpserver/server.go, flai/internal/mcpserver/shared.go, flai/internal/mcpserver/issues.go, flai/internal/mcpserver/session_test.go, flai/cmd/mcp.go, flai/cmd/mcp_http.go, flai/cmd/orchestrate_permissions_test.go, flai/internal/serve/restart_test.go, flai/internal/mcpserver/server_test.go]
usage:
  source: log
  seconds: 936
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 51
      output: 14621
      cache_read: 8247775
      cache_write: 221780
      cost: 3.1837
---
# T-2904 mcpserver.Options carry the session's story, role, session, reasons, and started-by, which flai mcp fills from the environment

## Work

The first layer: it waits for no task. The server reads the session from the process environment when it is made, so a test that sets one of those variables cannot run in parallel (I-0234, D-0443). Move those reads out of the package:

- Give `mcpserver.Options` the session's `Story`, `Role`, `Session`, `Reasons`, and `StartedBy`. Add one exported helper, such as `mcpserver.SessionFromEnv()`, that reads `FLAI_STORY`, `FLAI_ROLE`, `FLAI_SESSION`, `harness.ReasonsEnv`, and `guard.StartedByEnv` as `newServer` reads them today.
- `newServer` (`server.go`, about lines 203 to 214) takes them from `opt`. That covers folder mode too, since `folder.go:151` passes its `opt` to `newServer` for each project it finds. `servedSession` (`shared.go:180`) takes the role, story, and started-by as arguments. `recordingStory` (`issues.go:196`) reads `s.story`, not `FLAI_STORY`.
- `flai/cmd/mcp.go` and `mcpHTTPOptions` in `flai/cmd/mcp_http.go` fill the fields from the helper, so what `flai mcp` does over stdio and over HTTP stays the same. Over HTTP the helper reads the listening process's environment, as `New` does today. Whether an HTTP session should carry its own story and role is a separate question, outside this story.
- A child flai the server runs inherits the process environment. `task.go:62` already gives it `FLAI_SESSION` from `s.session`, and `subproject.go:152` does the same. Where a child must see the story or role for a converted test to pass, add it to that child's environment from the server's fields in `task.go`, not from the process.
- The tests outside the package that set the variables before they call `mcpserver.New` pass the fields instead: `flai/cmd/orchestrate_permissions_test.go`, `flai/cmd/role_test.go`, and `flai/internal/serve/restart_test.go`. Keep a `t.Setenv` there only where the code under test, not the server, reads the variable.
- Add `flai/internal/mcpserver/session_test.go`. It shows that a server made with `Options.Story` and `Options.Role` set acts as that story's agent and that role with the environment empty. It also shows that `SessionFromEnv` reads each variable. Only that one test sets the environment.

## Done when

- [ ] No file of `flai/internal/mcpserver` other than a test reads `FLAI_STORY`, `FLAI_ROLE`, `FLAI_SESSION`, `FLAI_REASONS`, or `FLAI_STARTED_BY`. `grep -n 'Getenv' flai/internal/mcpserver/*.go` lists only `versions.go`'s `FLAI_RELEASES_API`, which the next task moves.
- [ ] `flai mcp`, over stdio and over HTTP, fills the session from the environment through the one helper.
- [ ] The test in `session_test.go` passes.
- [ ] `flai test` passes on `flai/internal/mcpserver`, `flai/cmd`, and `flai/internal/serve`. The package's own tests may still use `together` until the later tasks.

## Notes
