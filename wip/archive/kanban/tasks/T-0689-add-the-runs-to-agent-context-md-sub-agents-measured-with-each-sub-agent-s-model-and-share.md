---
id: T-0689
type: task
nature: research
title: "Add the runs to agent-context.md § Sub-agents › Measured with each sub-agent's model and share"
status: done
parent: S-0190
owner: arobson
created: 2026-10-02T12:11:00Z
updated: 2026-10-02T16:57:11Z
transitions:
  - to: ready
    at: 2026-10-02T12:11:06Z
    by: agent-S-0190
  - to: in-progress
    at: 2026-10-02T16:55:39Z
    by: agent-S-0190
  - to: done
    at: 2026-10-02T16:57:11Z
    by: agent-S-0190
stream: S-0190
tags: [cli]
touches: [design/system/agent-context.md]
after: [T-0688]
usage:
  source: log
  seconds: 92
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 21
      output: 8559
      cache_read: 756644
      cache_write: 26264
      cost: 0.5327
---
# T-0689 Add the runs to agent-context.md § Sub-agents › Measured with each sub-agent's model and share

## Work

Add S-0194's and S-0193's rows to the table in `design/system/agent-context.md § Sub-agents › Measured`, with each sub-agent's model and its share of the cost, and their comparisons, and say whether delegation now costs less than not delegating. Waits for T-0688, whose numbers it writes.

## Done when

- The table has both rows, each sub-agent's model and share, and the comparison says whether the cost of delegation fell below the cost of not delegating.

## Notes
