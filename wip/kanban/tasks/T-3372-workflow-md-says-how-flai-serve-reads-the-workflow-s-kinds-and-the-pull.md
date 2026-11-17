---
id: T-3372
type: task
nature: feature
title: workflow.md says how flai serve reads the workflow's kinds and the pull
status: done
parent: S-0742
owner: alex
created: 2026-11-15T01:55:42Z
updated: 2026-11-15T01:55:53Z
transitions:
  - to: ready
    at: 2026-11-15T01:55:45Z
    by: agent-S-0742
  - to: in-progress
    at: 2026-11-15T01:55:45Z
    by: agent-S-0742
  - to: done
    at: 2026-11-15T01:55:53Z
    by: agent-S-0742
stream: S-0742
tags: []
touches: [design/system/workflow.md]
after: [T-3213, T-3215, T-3216, T-3217]
usage:
  source: log
  seconds: 8
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 3
      output: 883
      cache_read: 357553
      cache_write: 8227
      cost: 0.1409
---
# T-3372 workflow.md says how flai serve reads the workflow's kinds and the pull

## Work

- Add to `design/system/workflow.md` § The workflow file how `flai serve` reads the workflow: the pull and the active-work room for a story's agent, the restarts by the active kind, ready-empty and the starved board, and the reasons by the kind of the state a move enters.
- Waits for T-3213, T-3215, T-3216, and T-3217: it describes what they built.

## Done when

- The paragraph says what the four tasks built, and the markdown lint passes.

## Notes
