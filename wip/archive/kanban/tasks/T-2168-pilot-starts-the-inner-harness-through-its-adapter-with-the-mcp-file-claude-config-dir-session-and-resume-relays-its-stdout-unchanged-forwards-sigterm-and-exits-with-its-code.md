---
id: T-2168
type: task
nature: feature
title: pilot starts the inner harness through its adapter with the MCP file, CLAUDE_CONFIG_DIR, session and resume, relays its stdout unchanged, forwards SIGTERM, and exits with its code
status: done
parent: S-0507
owner: alex
created: 2026-10-26T04:50:27Z
updated: 2026-10-28T23:06:14Z
transitions:
  - to: ready
    at: 2026-10-28T23:00:29Z
    by: agent-S-0507
  - to: in-progress
    at: 2026-10-28T23:00:31Z
    by: agent-S-0507
  - to: done
    at: 2026-10-28T23:06:14Z
    by: agent-S-0507
stream: S-0507
tags: [cli, docker]
touches: [flai/internal/pilot/run.go, flai/internal/pilot/run_test.go, flai/pilot/main.go, flai/internal/harness/adapters.go]
after: [T-2161]
usage:
  source: log
  seconds: 343
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 33
      output: 13446
      cache_read: 2447062
      cache_write: 67860
      cost: 1.2006
---
# T-2168 pilot starts the inner harness through its adapter with the MCP file, CLAUDE_CONFIG_DIR, session and resume, relays its stdout unchanged, forwards SIGTERM, and exits with its code

## Work

Write `flai/internal/pilot/run.go`: `Run(ctx, req, stdout, stderr) int` looks up the inner adapter in `harness.Adapters`, calls its `Start` with the request (so the claude-code adapter writes the HTTP MCP file under `req.Cache`, sets `CLAUDE_CONFIG_DIR`, and gives `--session-id` on a first start or `--resume` when the session's folder under the cache already holds one, as `resumes()` decides on the host), fills `EnvFrom` from pilot's own environment, sets `FLAI_AGENT` and `FLAI_SESSION` from the request, and starts the process with the inner harness's program from `PATH` in the image. It copies the child's standard output to `stdout` line by line unchanged, writes its own lines to `stderr` prefixed `pilot:`, never writing a token or the environment, forwards `SIGTERM` and `SIGINT` to the child and waits for it, and returns the child's exit code, or `ExitUsage` when the start itself fails. Write `flai/pilot/main.go`: read the request from standard input with `Read` and `Check`, call `Run`, and `os.Exit` with its code. Test in `run_test.go` with a fake inner harness, a shell script on a temporary `PATH`, that echoes lines, exits with a given code, and writes the argv it got: the relay is byte-for-byte, the exit code passes through, a second run with the same cache folder gets `--resume`, and a signal reaches the child. Waits for T-2161 for the request.

## Done when

`go build ./pilot` in `flai/` builds, `flai test flai/internal/pilot` passes with the relay, exit code, resume, and signal tests, and `echo '{}' | pilot` exits 2 with one line on standard error.

## Notes

`harness.Start` has `Argv`, `Env`, and `EnvFrom`; the inner program is `Argv[0]`, resolved on the image's `PATH`.
