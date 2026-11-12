---
id: T-3059
type: task
nature: improvement
title: The claude-code adapter passes the start request's outside folders as --add-dir, and Monitor joins the default --allowedTools
status: done
parent: S-0701
owner: alex
created: 2026-11-10T10:35:24Z
updated: 2026-11-12T00:05:29Z
transitions:
  - to: ready
    at: 2026-11-11T23:36:01Z
    by: agent-S-0701
  - to: in-progress
    at: 2026-11-11T23:36:01Z
    by: agent-S-0701
  - to: done
    at: 2026-11-12T00:05:29Z
    by: agent-S-0701
stream: S-0701
tags: [harness, serve]
touches: [flai/internal/harness/harness.go, flai/internal/harness/adapters.go, flai/internal/harness/harness_test.go, flai/internal/serve/agents_start.go, flai/internal/serve/agents_test.go, flai/cmd/serve_harness.go, flai/cmd/serve_agent.go, docs/users/flai-reference.md, flai/cmd/serve_actions_test.go]
after: [T-3054]
usage:
  source: log
  seconds: 1768
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 108
      output: 32081
      cache_read: 6144183
      cache_write: 276074
      cost: 3.4695
---
# T-3059 The claude-code adapter passes the start request's outside folders as --add-dir, and Monitor joins the default --allowedTools

## Work

Add the outside folders to the start request (`Request` in `flai/internal/harness/harness.go`), not only to `Host`, because pilot starts the inner claude with an empty host (`flai/internal/pilot/run.go`). In `flai/internal/serve/agents_start.go`, fill it from the story's harness's host settings for a story's agent. In `flai/internal/harness/adapters.go`, have `claudeCode.Start` add one `--add-dir <folder>` per folder after `--permission-prompt-tool`. Change `DefaultHost` to `--allowedTools Bash,Monitor,mcp__flai`. Update the default arguments the help text of `flai/cmd/serve_harness.go` and `flai/cmd/serve_agent.go` names, and regenerate `docs/users/flai-reference.md`.

Reproduce I-0265 as a test: a story's agent started on a host with outside folders gets `--add-dir` for each, and the default arguments name `Monitor`. Before you rely on the inheritance, check once with a scratch `claude -p` on the cheapest model that an explorer sub-agent's Read of an `--add-dir` folder passes. Write the outcome in the narrative's Decisions.

This task waits for T-3054, whose settings it reads.

## Done when

- Tests in `flai/internal/harness/harness_test.go` and `flai/internal/serve/agents_test.go` fail without the change and pass with it.
- The scratch check's outcome is in the narrative's Decisions.
- `flai test` passes on the changed paths.

## Notes
