---
id: T-1596
type: task
nature: improvement
title: The MCP tool task_done closes a task with the installed flai when the session's flai is a minor version behind and the installed one has caught up
status: done
parent: S-0390
owner: alex
created: 2026-10-15T02:59:42Z
updated: 2026-10-15T03:15:06Z
transitions:
  - to: ready
    at: 2026-10-15T03:03:07Z
    by: agent-S-0390
  - to: in-progress
    at: 2026-10-15T03:03:08Z
    by: agent-S-0390
  - to: done
    at: 2026-10-15T03:15:06Z
    by: agent-S-0390
stream: S-0390
tags: [mcp, go]
touches: [flai/internal/taskdone/taskdone.go, flai/internal/taskdone/taskdone_test.go, flai/internal/mcpserver/task.go, flai/internal/mcpserver/task_test.go, flai/internal/mcpserver/installed.go, flai/internal/mcpserver/installed_test.go, flai/internal/mcpserver/server.go, flai/cmd/mcp.go, flai/cmd/mcp_http.go]
after: [T-1595]
usage:
  source: log
  seconds: 718
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 82
      output: 24483
      cache_read: 4987041
      cache_write: 150669
      cost: 2.4796
---
# T-1596 The MCP tool task_done closes a task with the installed flai when the session's flai is a minor version behind and the installed one has caught up

## Work

Build what T-1595's ADR decides. It waits for T-1595 because the ADR settles the remedy this task builds.

- In `flai/internal/taskdone/taskdone.go`, make `refuseOutdated`'s refusal a typed error, carrying the running and newest versions, so a caller can tell it from other failures. It still changes nothing.
- In `flai/internal/mcpserver/installed.go`, read the version of the binary at the server's own path. That is `os.Executable` by default, and a field on the server in `server.go` lets a test set it. Run it as `version --json`, as `flai/cmd/selfupgrade.go` does for the installed path.
- In `flai/internal/mcpserver/task.go`, when `taskdone.Run` returns that error and the installed version is at least the newest release, run the installed binary's `task done` with the tool's arguments as JSON, in the story's worktree, with the agent's `FLAI_AGENT` and `FLAI_SESSION`. Return its result, and say which flai closed the task and that the session's MCP server stays on its own version until the session restarts. Otherwise return the refusal, with the installed version and the call-again advice the ADR settles.
- Check that `os.Executable` gives the installed path after `flai self-upgrade` replaced the file on Linux, where `/proc/self/exe` then carries a deleted suffix.

## Done when

- A test in `flai/internal/mcpserver/task_test.go` reproduces I-0147. A server a minor version behind a tagged release refuses `task_done` as before when the installed binary is just as old. It closes the task through a fake installed binary that reports the newer version. The test checks that binary's arguments, working folder, and agent.
- `flai/internal/taskdone/taskdone_test.go` covers the typed refusal. `TestRunRefusesAMinorOrMajorVersionBehind` and `TestTaskDoneRefusesAFlaiAMinorVersionBehind` still pass.
- `flai test` passes on the changed paths.

## Notes

`installed.go` and its test are new files whose names are predicted. A development build is never outdated (`release.FlaiOutdated`), so the tree's `scripts/flai.sh` is unaffected.
