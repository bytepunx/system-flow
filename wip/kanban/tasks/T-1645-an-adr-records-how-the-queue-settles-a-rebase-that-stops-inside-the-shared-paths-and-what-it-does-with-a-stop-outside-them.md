---
id: T-1645
type: task
nature: feature
title: An ADR records how the queue settles a rebase that stops inside the shared paths, and what it does with a stop outside them
status: done
parent: S-0401
owner: alex
created: 2026-10-16T05:01:06Z
updated: 2026-10-19T21:54:38Z
transitions:
  - to: ready
    at: 2026-10-19T21:54:28Z
    by: agent-S-0401
  - to: in-progress
    at: 2026-10-19T21:54:29Z
    by: agent-S-0401
  - to: done
    at: 2026-10-19T21:54:38Z
    by: agent-S-0401
stream: S-0401
tags: [cli]
touches: [design/adrs]
usage:
  source: log
  seconds: 9
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 2
      output: 935
      cache_read: 116876
      cache_write: 3585
      cost: 0.0636
---
# T-1645 An ADR records how the queue settles a rebase that stops inside the shared paths, and what it does with a stop outside them

## Work

- Write the ADR with `flai adr new`, refining ADR-0098, ADR-0126, and ADR-0157.
- Weigh how to settle a stop on a shared markdown file that no existing step settles, such as `design/system/flai-cli.md`. Options include merging section by section by heading, as ADR-0157 merges rows, keeping both sides' new sections, and leaving the story queued until the story ahead lands. Name the cases it cannot settle.
- Record that such a story is never sent back, and that a stop outside the shared paths is told to the agent by message, as ADR-0121 does.
- First in its layer: it waits for no task.

## Done when

The ADR is accepted, indexed, and `flai check --strict` passes.

## Notes
