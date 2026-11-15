---
id: T-3211
type: task
nature: feature
title: flai-cli.md and the user guide list the workflow findings and what clears each
status: done
parent: S-0741
owner: alex
created: 2026-11-12T16:06:48Z
updated: 2026-11-15T20:47:08Z
transitions:
  - to: ready
    at: 2026-11-15T20:46:10Z
    by: agent-S-0741
  - to: in-progress
    at: 2026-11-15T20:46:10Z
    by: agent-S-0741
  - to: done
    at: 2026-11-15T20:47:08Z
    by: agent-S-0741
stream: S-0741
tags: [cli]
touches: [design/system/flai-cli.md, docs/users/flai.md]
after: [T-3204]
usage:
  source: log
  seconds: 58
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 26
      output: 11784
      cache_read: 2131203
      cache_write: 132358
      cost: 1.4232
---
# T-3211 flai-cli.md and the user guide list the workflow findings and what clears each

## Work

Waits for T-3204: it lists the codes that task adds. Add each `workflow.` finding to `design/system/flai-cli.md`'s `flai check` section and to `docs/users/flai.md`, with what it means and what the operator changes in the workflow file to clear it.

## Done when

Both documents list every `workflow.` code with its fix.

## Notes
