---
id: T-3247
type: task
nature: feature
title: flai stats prints the workflow's own timers, and stats.get on the host carries them
status: backlog
parent: S-0752
owner: alex
created: 2026-11-12T16:09:24Z
updated: 2026-11-12T16:09:24Z
transitions: []
stream: S-0752
tags: [cli]
touches: [flai/cmd/stats.go, flai/cmd/stats_test.go, flai/internal/hostapi/reads.go, flai/internal/hostapi/reads_test.go]
after: [T-3243]
---
# T-3247 flai stats prints the workflow's own timers, and stats.get on the host carries them

## Work

Waits for T-3243: it prints the section that task adds. Make `flai stats` print the workflow's own timers after the baseline's figures, and `flai stats --json` and the host's `stats.get` carry the section. Test both with a workflow that defines a timer of its own.

## Done when

`stats_test.go` and `reads_test.go` show the timer in the text, the JSON, and `stats.get`.

## Notes
