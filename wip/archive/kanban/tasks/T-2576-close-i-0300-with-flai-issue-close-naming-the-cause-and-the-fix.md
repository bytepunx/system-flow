---
id: T-2576
type: task
nature: remediation
title: Close I-0300 with flai issue close, naming the cause and the fix
status: done
parent: S-0589
owner: alex
created: 2026-11-07T10:06:49Z
updated: 2026-11-07T10:23:51Z
transitions:
  - to: ready
    at: 2026-11-07T10:22:24Z
    by: system-flow
  - to: in-progress
    at: 2026-11-07T10:22:29Z
    by: system-flow
  - to: done
    at: 2026-11-07T10:23:51Z
    by: agent-S-0589
stream: S-0589
tags: [issues]
touches: [design/issues/I-0300-a-story-already-in-ready-when-the-agent-host-action-is-turned-on-gets-no-agent.md, design/issues/summary.md]
after: [T-2573]
usage:
  source: log
  seconds: 82
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 8
      output: 1425
      cache_read: 1136573
      cache_write: 16141
      cost: 0.3831
  strategic:
    - kind: orchestrator
      seconds: 3
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 0
          output: 8
          cache_read: 28935
          cache_write: 760
          cost: 0.0074
---
# T-2576 Close I-0300 with flai issue close, naming the cause and the fix

## Work

Close I-0300 with `flai issue close I-0300 --reason` (the MCP tool `issue_close` with `commit`), the reason saying that the cause it guessed does not occur: flai serve's launcher starts a story already in ready at the next look once the `agent` action is on (ADR-0041), as T-2573's test shows; S-0001 in flai-scratch-import got no agent because its project's default agent named no harness and no command was set, which `serve.log` recorded at 08:26:49Z, and S-0587 makes a default harness work out of the box. Then add a line under the issue's `## Remediation` naming T-2573's test.

## Done when

- I-0300 is closed, `design/issues/summary.md` no longer lists it, and both are committed on `story/S-0589`.

## Notes
