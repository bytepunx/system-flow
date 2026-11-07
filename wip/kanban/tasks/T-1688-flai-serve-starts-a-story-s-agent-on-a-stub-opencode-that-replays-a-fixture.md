---
id: T-1688
type: task
nature: feature
title: flai serve starts a story's agent on a stub opencode that replays a fixture stream, and its refusals for opencode are tested
status: done
parent: S-0406
owner: alex
created: 2026-10-16T05:14:38Z
updated: 2026-11-07T11:06:02Z
transitions:
  - to: ready
    at: 2026-11-07T10:55:42Z
    by: system-flow
  - to: in-progress
    at: 2026-11-07T10:55:45Z
    by: system-flow
  - to: done
    at: 2026-11-07T11:06:02Z
    by: agent-S-0406
stream: S-0406
tags: [cli]
touches: [flai/internal/serve/agents_test.go, flai/internal/serve/testdata/opencode-run.jsonl, flai/internal/serve/agents.go]
after: [T-1682, T-1687]
usage:
  source: log
  seconds: 617
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 70
      output: 409
      cache_read: 3440655
      cache_write: 120884
      cost: 1.2776
---
# T-1688 flai serve starts a story's agent on a stub opencode that replays a fixture stream, and its refusals for opencode are tested

## Work

- Write `flai/internal/serve/testdata/opencode-run.jsonl`: a short `--format json` stream, `step_start`, `tool_use` of one flai MCP tool, `text`, and `step_finish`, in the shape T-1682 read from the pinned release's event types. Say in its first test's comment which release and source file it follows.
- In `flai/internal/serve/agents_test.go`, write a stub `opencode` into a temporary folder, as `claudecheck_test.go` writes its stub `claude`, set it as the host's program for `opencode`, and start a story's agent on it. Assert the arguments (`run`, the prompt, `--model`, `--format json`, and `--session` on a resume), that the configuration the stub was given carries the MCP entry `flai mcp --agent <name>`, and that the stub's replayed stream is read to the end.
- Test that a story on `opencode` with `roles` is refused at start naming `agent.harnesses.opencode.guard`, and that one without `auto-approve` is refused unless the host sets `--deny-protected`.
- No test here needs an `opencode` binary, a network, or a key, and none is skipped for want of one.
- Waits for T-1682 (the event shape the fixture follows) and T-1687 (a start that writes the configuration).

## Done when

The stub start, the resume, and the refusal tests pass under `flai test` on a host with no `opencode` on `PATH`.

## Notes
