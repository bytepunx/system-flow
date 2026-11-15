---
id: T-3231
type: task
nature: feature
title: The roles package has one condition naming the state a move goes to, and treats to-review and to-backlog as the ADR decides
status: done
parent: S-0748
owner: alex
created: 2026-11-12T16:08:17Z
updated: 2026-11-15T20:40:27Z
transitions:
  - to: ready
    at: 2026-11-15T20:16:33Z
    by: agent-S-0748
  - to: in-progress
    at: 2026-11-15T20:16:33Z
    by: agent-S-0748
  - to: done
    at: 2026-11-15T20:40:26Z
    by: agent-S-0748
stream: S-0748
tags: [cli]
touches: [flai/internal/roles/conditions.go, flai/internal/roles/conditions_test.go, flai/internal/roles/roles.go, flai/internal/check/roles_test.go, flai/internal/roles/roles_test.go]
usage:
  source: log
  seconds: 1431
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 86
      output: 389
      cache_read: 5679342
      cache_write: 472409
      cost: 2.3253
---
# T-3231 The roles package has one condition naming the state a move goes to, and treats to-review and to-backlog as the ADR decides

## Work

Add the condition the ADR names to `conditions.go`, taking a state's name, qualifying a `moved` event and the `item_move` and `flai move` calls, and matching a move into that state. Refuse at load a state the workflow lacks, as `roles.file`. Keep `to-review` and `to-backlog` as aliases or report them as `roles.file` naming the condition to write, as the ADR decides, in `check_roles.go`. Test both.

## Done when

`conditions_test.go` and `roles_test.go` pass with the new condition, a state the workflow lacks, and the old names.

## Notes
