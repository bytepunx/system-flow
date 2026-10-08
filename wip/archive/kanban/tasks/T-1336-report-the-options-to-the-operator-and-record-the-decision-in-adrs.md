---
id: T-1336
type: task
nature: research
title: Report the options to the operator and record the decision in ADRs
status: done
parent: S-0339
owner: alex
created: 2026-10-08T07:20:56Z
updated: 2026-10-08T08:32:17Z
transitions:
  - to: ready
    at: 2026-10-08T08:04:09Z
    by: agent-S-0339
  - to: in-progress
    at: 2026-10-08T08:04:09Z
    by: agent-S-0339
  - to: done
    at: 2026-10-08T08:32:17Z
    by: agent-S-0339
stream: S-0339
tags: [research, agents, adr]
touches: [design/system/agent-adapters.md, design/adrs, design/adrs/README.md]
after: [T-1335]
usage:
  source: log
  seconds: 422
  estimated: true
  models:
    - model: claude-fable-5-1
      input: 923
      output: 50751
      cache_read: 4891141
      cache_write: 275128
      cost: 8.0659
---
# T-1336 Report the options to the operator and record the decision in ADRs

## Work

Open a thread on S-0339 with `thread_open`. Lead with the recommendation from `## Abstractions`, then list the options and the questions only the operator can answer:

- the harness, provider, and model split in the agent schema
- which adapters the epic builds first
- whether a harness without a guard hook or a permission handler may run a story's agent, and with what limits
- how usage and cost are measured when a harness logs no cost
- whether a paid trial against LiteLLM or OpenRouter is wanted before the epic.

Wait for the answer with `wait_for_events`. Record the decision at once. Write one ADR with `flai adr new` per decision the answer makes, and name in it the ADRs it supersedes or refines. On each ADR it supersedes, set `superseded_by` in the front matter, the one edit an accepted ADR allows. Add a `## Decision` section to `design/system/agent-adapters.md` that links every new ADR. Waits for T-1335 because the thread reports the whole document.

## Done when

- The thread is answered, and the answer is recorded in the narrative's `## Decisions`.
- Each decision has an ADR, listed in `design/adrs/README.md`.
- `design/system/agent-adapters.md` has a `## Decision` section linking each new ADR.
- The story's second criterion is ticked.

## Notes
