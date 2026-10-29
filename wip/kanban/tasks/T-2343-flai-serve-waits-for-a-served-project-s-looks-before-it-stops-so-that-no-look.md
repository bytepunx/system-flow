---
id: T-2343
type: task
nature: remediation
title: flai serve waits for a served project's looks before it stops, so that no look writes after it, which failed the new loop test's cleanup
status: done
parent: S-0418
owner: alex
created: 2026-10-29T17:47:04Z
updated: 2026-10-29T17:49:04Z
transitions:
  - to: ready
    at: 2026-10-29T17:48:58Z
    by: agent-S-0418
  - to: in-progress
    at: 2026-10-29T17:48:59Z
    by: agent-S-0418
  - to: done
    at: 2026-10-29T17:49:04Z
    by: agent-S-0418
stream: S-0418
tags: [cli]
touches: [flai/internal/serve/serve.go]
usage:
  source: log
  seconds: 5
  estimated: true
  priced_by: harness
  models: []
---
# T-2343 flai serve waits for a served project's looks before it stops, so that no look writes after it, which failed the new loop test's cleanup

## Work

- `Run` waits for a served project's loop and its watcher's looks once its client has stopped, before the project is done, so that no look writes to the project or the serve folder after it (I-0233).
- S-0418's loop test, `TestServeStartsADefinedRoleOnAnEventAndOnASchedule`, failed the close-out's go-test tier on its temp folder's cleanup ("directory not empty"): a look still in flight wrote after `Run` returned. The loop looks at two more starters since T-2340.

Found at the close-out, after T-2340.

## Done when

`flai test flai/internal/serve` passes, and the close-out's go-test tier passes.

## Notes
