---
id: T-3193
type: task
nature: feature
title: Settle the ADR's open choices with the operator, accept it, and mark ADR-0003 and ADR-0004 superseded or refined by it
status: backlog
parent: S-0736
owner: alex
created: 2026-11-12T16:05:47Z
updated: 2026-11-12T16:05:47Z
transitions: []
stream: S-0736
tags: [cli]
touches: [design/adrs/README.md, design/adrs/0003-work-item-hierarchy.md, design/adrs/0004-workflow-states-and-transitions.md]
after: [T-3187]
---
# T-3193 Settle the ADR's open choices with the operator, accept it, and mark ADR-0003 and ADR-0004 superseded or refined by it

## Work

Waits for T-3187: it asks about the ADR that task writes. Open one thread on the story listing the ADR's open choices, each with your recommendation, and end the session until it is answered. Then change the ADR to the answers, accept it with `flai adr accept`, and set `superseded_by` on ADR-0003 and ADR-0004 where the ADR supersedes them, the one edit an accepted ADR allows. Where it refines them, leave them as they are.

## Done when

The ADR is accepted with the operator's answers, and ADR-0003 and ADR-0004 say what became of them.

## Notes
