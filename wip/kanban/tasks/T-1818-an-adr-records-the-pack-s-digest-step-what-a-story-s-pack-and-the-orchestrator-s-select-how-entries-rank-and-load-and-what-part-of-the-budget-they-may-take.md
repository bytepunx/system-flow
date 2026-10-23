---
id: T-1818
type: task
nature: feature
title: "An ADR records the pack's digest step: what a story's pack and the orchestrator's select, how entries rank and load, and what part of the budget they may take"
status: done
parent: S-0430
owner: alex
created: 2026-10-16T06:04:32Z
updated: 2026-10-23T19:14:42Z
transitions:
  - to: ready
    at: 2026-10-23T19:13:45Z
    by: agent-S-0430
  - to: in-progress
    at: 2026-10-23T19:13:46Z
    by: agent-S-0430
  - to: done
    at: 2026-10-23T19:14:42Z
    by: agent-S-0430
stream: S-0430
tags: [cli]
touches: [design/adrs, design/adrs/README.md]
usage:
  source: log
  seconds: 56
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 13
      output: 5150
      cache_read: 805922
      cache_write: 27345
      cost: 0.4276
---
# T-1818 An ADR records the pack's digest step: what a story's pack and the orchestrator's select, how entries rank and load, and what part of the budget they may take

## Work

Write the ADR with `flai adr new`, refining ADR-0049 and ADR-0075. It records:

- Which entries a story's pack selects: active and compacted entries whose `topics` meet the story's topics, or whose `paths` meet its claim (ADR-0046), excluding absorbed ones.
- Where they go: a step of their own after the briefs, ranked by BM25 against the story's title, goal, and criteria, loading whole while the budget has room, a compacted entry as its title alone. Whether the step has a share of the budget of its own, so ranked design is not crowded out, and what it is.
- The sub-agents' packs: entries as one-line briefs only, within their half budget.
- The orchestrator's pack: `answer` and `recommendation` entries ranked against the titles and first entries of the threads awaiting the operator.
- Each item's `step: digest` and reason, in text and JSON, and that parts (ADR-0104) are unchanged.

First layer: the code tasks build what this decides.

## Done when

The ADR is accepted and listed in `design/adrs/README.md`, and `flai check --strict` reports nothing on it.

## Notes

The ADR's file name is not known until `flai adr new` makes it, so the task touches the shared `design/adrs` folder.
