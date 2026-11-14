---
id: T-3354
type: task
nature: feature
title: The host API method log.read answers a part of a log flai wrote, confined to the project's and flai serve's log folders
status: done
parent: S-0785
owner: alex
created: 2026-11-14T10:28:45Z
updated: 2026-11-14T10:50:07Z
transitions:
  - to: ready
    at: 2026-11-14T10:29:10Z
    by: agent-S-0785
  - to: in-progress
    at: 2026-11-14T10:29:10Z
    by: agent-S-0785
  - to: done
    at: 2026-11-14T10:50:07Z
    by: agent-S-0785
stream: S-0785
tags: []
touches: [flai/internal/hostapi/logs.go, flai/internal/hostapi/logs_test.go, flai/internal/hostapi/writes_methods.go, flai/internal/hostapi/guardreads_test.go, flaiover/src/lib/server/agent.ts, design/adrs/README.md, flai/internal/hostapi/writes_actions.go, flai/internal/serve/serve.go, flai/internal/hostapi/images.go, design/adrs/0300-the-dashboard-reads-a-log-flai-wrote-through-log-read-confined-to-the-project-s.md]
usage:
  source: log
  seconds: 1257
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 109
      output: 36778
      cache_read: 7204853
      cache_write: 177425
      cost: 3.2263
  strategic:
    - kind: orchestrator
      seconds: 3
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 1
          output: 7
          cache_read: 20792
          cache_write: 579
          cost: 0.0063
---
# T-3354 The host API method log.read answers a part of a log flai wrote, confined to the project's and flai serve's log folders

## Work

Add a read-only host API method, `log.read`, that the dashboard calls to show the contents of a log file flai wrote, so that an operator without a shell on the host can read a log a status or error message names.

- Params: `path`, the log's path as a message names it: absolute, `~/`-prefixed (the home of the user flai runs as), or relative to the project's root (as `kluster.LogPart.Log` gives it); and `offset`, a whole number of zero or more.
- Answer, in the shape `images.log` and `kluster.log` answer: `path` (as given), `exists`, `offset`, `next`, `size`, `text` (at most 64 KB from offset, ending on a whole character), `reset`.
- Confinement: only a regular file whose name ends in `.log`, whose real path (symlinks resolved) lies under one of the folders flai writes logs to for this project: the project's `.flai-cache/` (but not `.flai-cache/worktrees/`), and flai serve's state folder for the logs it writes (agents, image jobs, the serve and host logs), and only the project's own files there where the file name says whose they are. Every other path is refused with a `channel.Error` that says what is readable. List every log path producer that reaches the dashboard and check each lies inside.
- Register it as a read in `writes_methods.go` (or the spec list that fits), list it in `guardreads_test.go`, and add it to flaiover's `REQUIRED_METHODS` in `flaiover/src/lib/server/agent.ts`, which `contract_test.go` checks (D-0612).
- Record the decision as an ADR with `flai adr new` (no `--commit`).
- Waits for nothing: the first layer.

## Done when

- `log.read` answers a part of each kind of log the dashboard names (agent run, image job, kluster step, verify or close-out, checks), and refuses a path outside the confinement, a non-`.log` file, a symlink out, and `..` escapes, each with a test.
- `flai test` passes on the paths changed.

## Notes
