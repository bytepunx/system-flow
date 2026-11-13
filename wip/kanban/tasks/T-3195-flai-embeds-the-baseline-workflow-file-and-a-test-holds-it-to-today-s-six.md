---
id: T-3195
type: task
nature: feature
title: flai embeds the baseline workflow file, and a test holds it to today's six states and transition map pair by pair
status: backlog
parent: S-0737
owner: alex
created: 2026-11-12T16:05:51Z
updated: 2026-11-12T16:05:51Z
transitions: []
stream: S-0737
tags: [cli]
touches: [flai/internal/workflow/baseline.go, flai/internal/workflow/baseline_test.go, flai/internal/workflow/baseline/board.md]
after: [T-3188]
---
# T-3195 flai embeds the baseline workflow file, and a test holds it to today's six states and transition map pair by pair

## Work

Waits for T-3188: the baseline is parsed by the parser that task writes. Write the baseline file in the ADR's format with today's six states, their kinds, WIP limits, and order, today's transitions for epics, stories, and tasks, and the baseline timers. Embed it with `//go:embed` as `flai/internal/roles/baseline.go` embeds the role files, and give `Baseline()`, parsed once. Test that every (from, to, type) pair `workitem`'s `allowed` map and `allowedFor` give today is in the baseline and no other is.

## Done when

`baseline_test.go` passes, and a change to the baseline that adds or drops a move fails it.

## Notes
