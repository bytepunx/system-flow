---
id: T-3200
type: task
nature: feature
title: "workflow.md, work-hierarchy.md § States, and the workflows README describe the workflow file and link the accepted ADR"
status: done
parent: S-0736
owner: alex
created: 2026-11-12T16:06:06Z
updated: 2026-11-14T17:56:23Z
transitions:
  - to: ready
    at: 2026-11-14T17:55:35Z
    by: agent-S-0736
  - to: in-progress
    at: 2026-11-14T17:55:36Z
    by: agent-S-0736
  - to: done
    at: 2026-11-14T17:56:23Z
    by: agent-S-0736
stream: S-0736
tags: [cli]
touches: [design/system/workflow.md, design/system/work-hierarchy.md, design/workflows/README.md]
after: [T-3193]
usage:
  source: log
  seconds: 47
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 12
      output: 4596
      cache_read: 1235849
      cache_write: 31664
      cost: 0.5759
---
# T-3200 workflow.md, work-hierarchy.md § States, and the workflows README describe the workflow file and link the accepted ADR

## Work

Waits for T-3193: it describes the ADR as accepted, with the operator's answers. Say in `design/system/workflow.md` § Columns and WIP limits and § Transitions and who makes them that the board is the workflow file, where it is, and that a project without one runs flai's baseline. Rewrite `design/system/work-hierarchy.md` § States from a closed list to the baseline's states and their kinds. Add the workflow file to `design/workflows/README.md`. Link the ADR from each, and say nothing is built yet where that is so.

## Done when

The three documents describe the workflow file as the accepted ADR decides it.

## Notes
