---
id: T-1120
type: task
nature: improvement
title: flai accept's preview lists the files a story's branch changes on a protected path, and an agent's acceptance of such a story is refused with the files named
status: done
parent: S-0286
owner: alex
created: 2026-10-06T22:54:04Z
updated: 2026-10-07T01:27:29Z
transitions:
  - to: ready
    at: 2026-10-07T01:20:41Z
    by: agent-S-0286
  - to: in-progress
    at: 2026-10-07T01:20:41Z
    by: agent-S-0286
  - to: done
    at: 2026-10-07T01:27:29Z
    by: agent-S-0286
stream: S-0286
tags: [flai]
touches: [flai/internal/preview/accept.go, flai/internal/preview/orchestrator.go, flai/cmd/accept.go, flai/cmd/accept_protected_test.go]
after: [T-1110]
usage:
  source: log
  seconds: 408
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 38
      output: 15333
      cache_read: 2625646
      cache_write: 66436
      cost: 1.2349
---
# T-1120 flai accept's preview lists the files a story's branch changes on a protected path, and an agent's acceptance of such a story is refused with the files named

## Work

In `flai/internal/preview/accept.go`, add a field to `Acceptance` listing the files the story's branch changes on a protected path, read with the shared list from T-1110. Add a line saying that only the operator accepts such a story. The JSON reaches `flai accept --dry-run` and the dashboard, so name the field for both.

Refuse an agent's acceptance of such a story, naming the files:

- **The orchestrator,** under `accept_reviews` (`flai accept --by orchestrator`): add a blocker code beside ADR-0093's in `orchestrator.go`, so the refusal also reaches the dashboard's host channel `accept.run`, which builds the same preview.
- **An agent on its own name:** refuse an acceptance whose `--by`, or default author, is an agent's name rather than the operator's. Use the way flai already tells an agent's name from a person's, as `workitem.IsOrchestrator` does for the orchestrator. If flai has no such way for a story's agent, ask on a thread on S-0286 before choosing one.
- The operator's own acceptance, `--by alex` or the dashboard's, goes through, with the files shown.

`flai accept`'s text output prints the files in its preview.

Waits for T-1110, which provides the shared list of protected paths.

## Done when

- `accept_protected_test.go` covers:
  - a story whose branch changes `.claude/settings.json` and `.mcp.json`, refused for the orchestrator and for an agent on its own name, each refusal naming both files;
  - the same story accepted by the operator;
  - a story with no protected change accepted by the orchestrator as before.
- `flai accept --dry-run` prints the files.
- `scripts/flai-test.sh` passes.

## Notes

Drafted by planner-S-0286. `flai/internal/hostapi/writes.go` is not touched here: `accept.run` builds its refusal from the same preview. If the story's agent finds that it does not, it widens the touches with `flai touches`.
