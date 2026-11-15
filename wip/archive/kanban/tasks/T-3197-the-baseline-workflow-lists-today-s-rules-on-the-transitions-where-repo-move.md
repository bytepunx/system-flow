---
id: T-3197
type: task
nature: feature
title: The baseline workflow lists today's rules on the transitions where Repo.Move applies them
status: done
parent: S-0738
owner: alex
created: 2026-11-12T16:05:57Z
updated: 2026-11-15T00:10:10Z
transitions:
  - to: ready
    at: 2026-11-15T00:04:17Z
    by: agent-S-0738
  - to: in-progress
    at: 2026-11-15T00:04:17Z
    by: agent-S-0738
  - to: done
    at: 2026-11-15T00:10:10Z
    by: agent-S-0738
stream: S-0738
tags: [cli]
touches: [flai/internal/workflow/baseline_test.go]
after: [T-3189]
usage:
  source: log
  seconds: 353
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 9
      output: 2503
      cache_read: 933209
      cache_write: 28518
      cost: 0.4388
---
# T-3197 The baseline workflow lists today's rules on the transitions where Repo.Move applies them

## Work

Waits for T-3189: the baseline names rules from the vocabulary that task adds. On the baseline's transitions, list the rules `Repo.Move` applies today, in its order: to `ready`, criteria present and not a draft; to `review`, at least one task, no open question, nothing uncommitted, and every criterion ticked as a warning; to `done`, every child closed and every criterion ticked; a reason to `cancelled` and from `review` back to `in-progress`. Extend `baseline_test.go` to hold each transition's rules.

## Done when

`baseline_test.go` passes with each baseline transition's rules as `Repo.Move` applies them today.

## Notes
