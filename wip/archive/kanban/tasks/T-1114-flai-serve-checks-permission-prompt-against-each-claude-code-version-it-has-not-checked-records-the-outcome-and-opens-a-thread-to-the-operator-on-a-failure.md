---
id: T-1114
type: task
nature: improvement
title: flai serve checks permission_prompt against each Claude Code version it has not checked, records the outcome, and opens a thread to the operator on a failure
status: done
parent: S-0286
owner: alex
created: 2026-10-06T22:53:44Z
updated: 2026-10-07T01:18:55Z
transitions:
  - to: ready
    at: 2026-10-07T01:04:24Z
    by: agent-S-0286
  - to: in-progress
    at: 2026-10-07T01:04:25Z
    by: agent-S-0286
  - to: done
    at: 2026-10-07T01:18:55Z
    by: agent-S-0286
stream: S-0286
tags: [flai]
touches: [flai/internal/serve/claudecheck.go, flai/internal/serve/claudecheck_test.go, flai/internal/serve/serve.go, flai/internal/harness/adapters.go]
after: [T-1106]
usage:
  source: log
  seconds: 870
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 131
      output: 52629
      cache_read: 9012531
      cache_write: 228043
      cost: 4.2388
---
# T-1114 flai serve checks permission_prompt against each Claude Code version it has not checked, records the outcome, and opens a thread to the operator on a failure

## Work

In `flai/internal/serve`, add the check the ADR describes:

- Find the version of the `claude` program flai serve starts for the claude-code harness. That is the host's program from `harness.DefaultHost` or the operator's `flai serve agent set`.
- If that version has a recorded outcome on this host, do nothing.
- Otherwise, make a scratch flai project in a temporary folder with one in-progress story and its worktree. Run one headless `claude -p` in it, with the same `--permission-mode` and `--permission-prompt-tool mcp__flai__permission_prompt` that `adapters.go` gives a story's agent. Ask it for one Write of a file under a `.claude/` folder in the story's worktree.
- The write passes when the file exists afterwards with the content asked for.
- Record the version, the outcome, and the time in flai serve's state on the host.
- On a failure, open one thread to the operator. Quote Claude Code's error from the run's output, and say that `.claude/` writes will not go through until it is fixed.

Run the check when flai serve starts, and before it starts a claude-code agent if the version has changed. Do not block or delay the agent's start on it. Reuse the harness's argv in `adapters.go` rather than copying the flags.

Use a model the operator's configuration names for the cheapest role, or Haiku. One check costs about 0.02 USD.

Waits for T-1106, because the ADR decides what the check runs and where its outcome lives. It shares no file with T-1110, so they run together.

## Done when

- `claudecheck_test.go` runs the check with a stand-in for `claude`, a script on `PATH` that prints a version and does or does not write the file. It shows:
  - an unchecked version is checked and recorded;
  - a recorded version is not checked again;
  - a failure opens one thread to the operator that quotes the stand-in's error.
- `scripts/flai-test.sh` passes.

## Notes

Drafted by planner-S-0286. `claudecheck.go` is the plan's name for a new file; the story's agent may name it otherwise and update the touches with `flai touches`.
