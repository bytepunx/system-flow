---
id: T-3199
type: task
nature: feature
title: workflow.md lists the baseline's rules by name on each transition
status: done
parent: S-0738
owner: alex
created: 2026-11-12T16:06:00Z
updated: 2026-11-15T00:10:16Z
transitions:
  - to: ready
    at: 2026-11-15T00:04:17Z
    by: agent-S-0738
  - to: in-progress
    at: 2026-11-15T00:04:17Z
    by: agent-S-0738
  - to: done
    at: 2026-11-15T00:10:16Z
    by: agent-S-0738
stream: S-0738
tags: [cli]
touches: [design/system/workflow.md]
after: [T-3189]
usage:
  source: log
  seconds: 359
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 11
      output: 3849
      cache_read: 934689
      cache_write: 29247
      cost: 0.454
---
# T-3199 workflow.md lists the baseline's rules by name on each transition

## Work

Waits for T-3189: it names the rules that task fixes. In `design/system/workflow.md` § Transitions and who makes them, name each baseline transition's rules by their vocabulary names, and say that a workflow of the operator's own lists rules from the same vocabulary on any transition.

## Done when

§ Transitions and who makes them names every baseline rule and links the ADR.

## Notes
