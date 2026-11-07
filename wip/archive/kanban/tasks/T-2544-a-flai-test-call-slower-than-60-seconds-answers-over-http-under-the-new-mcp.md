---
id: T-2544
type: task
nature: remediation
title: A flai test call slower than 60 seconds answers over HTTP under the new MCP file, as I-0228's third instance could not, and I-0228 is closed
status: done
parent: S-0583
owner: alex
created: 2026-11-07T07:06:57Z
updated: 2026-11-07T07:35:12Z
transitions:
  - to: ready
    at: 2026-11-07T07:32:12Z
    by: system-flow
  - to: in-progress
    at: 2026-11-07T07:32:14Z
    by: system-flow
  - to: done
    at: 2026-11-07T07:35:12Z
    by: agent-S-0583
stream: S-0583
tags: [cli, docker]
touches: [design/issues/I-0228-claude-code-drops-any-call-to-flai-s-http-mcp-server-that-answers-after-60-seconds-so-wait-for-events-wait-for-work-verify-and-a-slow-task-done-fail-for-an-agent-in-a-container.md, design/issues/summary.md]
after: [T-2542, T-2543]
usage:
  source: log
  seconds: 178
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 20
      output: 5080
      cache_read: 2622342
      cache_write: 38936
      cost: 0.9229
  strategic:
    - kind: orchestrator
      seconds: 2
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 0
          output: 6
          cache_read: 18267
          cache_write: 781
          cost: 0.0047
---
# T-2544 A flai test call slower than 60 seconds answers over HTTP under the new MCP file, as I-0228's third instance could not, and I-0228 is closed

## Work

This task waits for T-2542, whose MCP file it runs, and T-2543, so that the issue closes only once the docs say what fixed it.

Repeat I-0228's third instance against flai's own HTTP server. In `flai-agent-claude-code:local`, run one headless `claude -p` with `--strict-mcp-config` and the MCP file T-2542's `httpMCPConfig` writes, asking it to call the MCP tool `test` once on `flai/internal/guard` with this story, and nothing else. That call took about four minutes when I-0228 recorded it. Use the file's bytes as the adapter writes them, with the URL and the token by their variables, never written out. Name the session's agent as this story's probe, so that nothing it does reads as this story's agent. Record how long the call took and what it answered in the narrative's `## Decisions`.

If the call answers in under 60 seconds, the run proves nothing: pick paths whose tiers take longer, and say which.

`task_done` and `story_start` are not repeated. They meet the same per-server limit the `test` call does, and T-2541's held call measured that limit to 660 s. No call is measured between 660 s and the figure, so say so in the closing reason if the figure is higher.

Then close the issue with `flai issue close I-0228 --reason`, naming the ADR, what each call now does (the two waits and `permission_prompt` hold 45 seconds, `verify` answers at once, and `task_done`, `story_start`, and `test` run under the per-server `timeout`), and that a running image must be rebuilt with `make agent-images` to carry it.

## Done when

- A `test` call slower than 60 seconds over HTTP in the container returned its result, with its duration in the narrative.
- I-0228 is closed with that reason, and `design/issues/summary.md` no longer lists it.
- Both of the story's acceptance criteria can be ticked from T-2542's tests, this run, and this closing.

## Notes
