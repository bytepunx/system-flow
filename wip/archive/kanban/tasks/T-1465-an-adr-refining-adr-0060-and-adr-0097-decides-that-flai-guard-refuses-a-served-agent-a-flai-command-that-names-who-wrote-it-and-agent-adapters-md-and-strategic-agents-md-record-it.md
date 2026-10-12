---
id: T-1465
type: task
nature: improvement
title: An ADR refining ADR-0060 and ADR-0097 decides that flai guard refuses a served agent a flai command that names who wrote it, and agent-adapters.md and strategic-agents.md record it
status: done
parent: S-0369
owner: alex
created: 2026-10-12T17:05:34Z
updated: 2026-10-12T18:26:43Z
transitions:
  - to: ready
    at: 2026-10-12T18:18:04Z
    by: system-flow
  - to: in-progress
    at: 2026-10-12T18:18:05Z
    by: system-flow
  - to: done
    at: 2026-10-12T18:26:43Z
    by: agent-S-0369
stream: S-0369
tags: [flai, guard, adr]
touches: [design/adrs/README.md, design/system/agent-adapters.md, design/system/strategic-agents.md, design/adrs/0143-flai-guard-refuses-a-story-s-agent-the-planner-and-the-analyzer-a-call-that.md]
usage:
  source: log
  seconds: 518
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 12
      output: 5412
      cache_read: 835301
      cache_write: 26478
      cost: 0.4534
---
# T-1465 An ADR refining ADR-0060 and ADR-0097 decides that flai guard refuses a served agent a flai command that names who wrote it, and agent-adapters.md and strategic-agents.md record it

## Work

Write an ADR with `flai adr new` that refines ADR-0060 (the guard) and ADR-0097 (whose answer `permission_prompt` takes), as the operator confirmed on TH-0418:

- The guard refuses a served agent, a story's agent, the planner, or the analyzer, a call whose `--by` or `--placed-by` on a flai command, or whose `FLAI_AGENT` or `FLAI_SESSION` set, emptied, or unset on the command line, names anyone but that agent's own name. It does not refuse the flags or the variables as such, so an agent that sets `FLAI_AGENT` to its own name, as `tooling.md` asks, passes.
- Its own name is the `FLAI_AGENT` the guard's hook sees, as flai serve set it or the project's Claude Code settings override it, or the name flai serve gives the run: `agent-S-nnnn`, `planner-<item>` (from `FLAI_ITEM`), or `analyzer`. Its own session is the hook's `FLAI_SESSION`.
- `--by` on `flai order` (a policy) and `flai stats` (a grouping) names no author and is left alone.
- An emptied or unset `FLAI_AGENT` is refused because flai then falls back to the config's author, the operator (T-1466).
- The orchestrator keeps its own rule (`writesAsSelf`); the operator's own session, which no `FLAI_STORY`, `FLAI_ROLE`, or `FLAI_STARTED_BY` marks, and the dashboard, which runs flai outside any guard, are not held.
- Residuals: a shell write into a thread file by a story's own agent, a call to the HTTP MCP server with another `X-Flai-Agent`, and a name hidden from the guard's reading of the command line (a variable, a script file, a make target).

Record the decision in `design/system/agent-adapters.md` under `### Permissions` and `### Asking the owner before a protected write`, and in `design/system/strategic-agents.md` under the planner's `### The guard` and `### The analyzer's guard`, each linking the ADR. `flai adr new` adds the row to `design/adrs/README.md`; `flai task done` widens this task's touches with the ADR's file.

## Done when

- An accepted ADR says which agents the guard holds to their own name, what own name means, which flags and environment it refuses, what it leaves to the operator's session, the orchestrator, and the dashboard, and the residuals.
- `design/system/agent-adapters.md` and `design/system/strategic-agents.md` say the same and link the ADR.
- `design/adrs/README.md` lists it, and `flai check --strict` is clean.

## Notes

First layer, with the `threadAuthor` task: the guard task waits for this decision. Rewritten by the story's agent to the narrowing the operator confirmed on TH-0418.
