---
id: T-3219
type: task
nature: feature
title: The forge pass pushes, reads, and sends back a story by the state the acceptance leaves from and the send-back effect
status: done
parent: S-0743
owner: alex
created: 2026-11-12T16:07:16Z
updated: 2026-11-15T01:44:14Z
transitions:
  - to: ready
    at: 2026-11-15T01:07:33Z
    by: agent-S-0743
  - to: in-progress
    at: 2026-11-15T01:07:33Z
    by: agent-S-0743
  - to: done
    at: 2026-11-15T01:44:14Z
    by: agent-S-0743
stream: S-0743
tags: [cli]
touches: [flai/internal/serve/forgestatus.go, flai/internal/serve/forgestatus_test.go, flai/internal/serve/forge.go, flai/internal/serve/forge_test.go, flai/internal/serve/forgemerge.go, flai/internal/forge/state.go]
after: [T-3369]
usage:
  source: log
  seconds: 2201
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 50
      output: 14350
      cache_read: 5610434
      cache_write: 150891
      cost: 2.2672
---
# T-3219 The forge pass pushes, reads, and sends back a story by the state the acceptance leaves from and the send-back effect

## Work

Replace the references to a state's name in the forge files: flai serve pushes a story and opens its pull request in the state the acceptance leaves from (ADR-0166), reads its checks and reviews there (ADR-0167), and sends it back along the send-back transition on a failing check or a requested change. Test the push and the send back on a workflow with a state added.

## Done when

The forge files name no state, and their tests pass on both workflows.

## Notes
