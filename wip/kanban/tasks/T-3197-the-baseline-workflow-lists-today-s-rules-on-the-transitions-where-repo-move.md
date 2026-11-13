---
id: T-3197
type: task
nature: feature
title: The baseline workflow lists today's rules on the transitions where Repo.Move applies them
status: backlog
parent: S-0738
owner: alex
created: 2026-11-12T16:05:57Z
updated: 2026-11-12T16:05:57Z
transitions: []
stream: S-0738
tags: [cli]
touches: [flai/internal/workflow/baseline/board.md, flai/internal/workflow/baseline_test.go]
after: [T-3189]
---
# T-3197 The baseline workflow lists today's rules on the transitions where Repo.Move applies them

## Work

Waits for T-3189: the baseline names rules from the vocabulary that task adds. On the baseline's transitions, list the rules `Repo.Move` applies today, in its order: to `ready`, criteria present and not a draft; to `review`, at least one task, no open question, nothing uncommitted, and every criterion ticked as a warning; to `done`, every child closed and every criterion ticked; a reason to `cancelled` and from `review` back to `in-progress`. Extend `baseline_test.go` to hold each transition's rules.

## Done when

`baseline_test.go` passes with each baseline transition's rules as `Repo.Move` applies them today.

## Notes
