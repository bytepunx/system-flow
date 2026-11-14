---
id: T-3193
type: task
nature: feature
title: Settle the ADR's open choices with the operator, accept it, and mark ADR-0003 and ADR-0004 superseded or refined by it
status: done
parent: S-0736
owner: alex
created: 2026-11-12T16:05:47Z
updated: 2026-11-14T17:55:24Z
transitions:
  - to: ready
    at: 2026-11-14T17:47:28Z
    by: agent-S-0736
  - to: in-progress
    at: 2026-11-14T17:47:28Z
    by: agent-S-0736
  - to: ready
    at: 2026-11-14T17:47:44Z
    by: agent-S-0736
  - to: in-progress
    at: 2026-11-14T17:55:13Z
    by: agent-S-0736
  - to: done
    at: 2026-11-14T17:55:24Z
    by: agent-S-0736
stream: S-0736
tags: [cli]
touches: [design/adrs/README.md, design/adrs/0301-the-board-is-a-workflow-file-design-workflows-board-md-whose-front-matter-names.md]
after: [T-3187]
usage:
  source: log
  seconds: 27
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 8
      output: 3016
      cache_read: 810963
      cache_write: 20778
      cost: 0.3779
---
# T-3193 Settle the ADR's open choices with the operator, accept it, and mark ADR-0003 and ADR-0004 superseded or refined by it

## Work

Waits for T-3187: it asks about the ADR that task writes. Open one thread on the story listing the ADR's open choices, each with your recommendation, and end the session until it is answered. Then change the ADR to the answers, accept it with `flai adr accept`, and set `superseded_by` on ADR-0003 and ADR-0004 where the ADR supersedes them, the one edit an accepted ADR allows. Where it refines them, leave them as they are.

## Done when

The ADR is accepted with the operator's answers, and ADR-0003 and ADR-0004 say what became of them.

## Notes
