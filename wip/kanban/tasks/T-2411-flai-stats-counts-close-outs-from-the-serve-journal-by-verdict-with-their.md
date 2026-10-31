---
id: T-2411
type: task
nature: feature
title: flai stats counts close-outs from the serve journal, by verdict, with their durations, as metrics.md defines
status: backlog
parent: S-0546
owner: alex
created: 2026-10-31T07:44:33Z
updated: 2026-10-31T07:44:33Z
transitions: []
stream: S-0546
tags: [flai]
touches: [flai/internal/metrics/closeouts.go, flai/internal/metrics/closeouts_test.go, flai/internal/metrics/metrics.go, flai/internal/statsread/statsread.go, flai/internal/statsread/statsread_test.go, flai/cmd/stats.go]
after: [T-2403]
---
# T-2411 flai stats counts close-outs from the serve journal, by verdict, with their durations, as metrics.md defines

## Work

`flai stats` reads no serve journal today. `statsread.go` loads the merge queue's events alone (`queueEvents`). Read the journal's close-out entries in `flai/internal/statsread/statsread.go`, and pass them to the metrics as a new input in `flai/internal/metrics/metrics.go`. Count them in a new `flai/internal/metrics/closeouts.go`, as T-2396's metrics.md section defines: the count, and how many passed, failed, or were answered from the record, with their durations. Print the counts in `flai stats` and its `--json` (`flai/cmd/stats.go`). A project with no journal, or no close-out entries, prints none and does not fail.

`flai/cmd/stats.go` is in S-0552's half of the split on MS-0261, which changes the queue stats there. Keep the change to the lines that print the new counts, and tell S-0552's agent on that conversation. Whichever lands second rebases.

Waits for T-2403, whose journal entries it reads. Shares no path with T-2407, so the two run together.

## Done when

- [ ] A test in `flai/internal/metrics/closeouts_test.go` counts a journal fixture's close-outs by verdict, with their durations
- [ ] A test in `flai/internal/statsread/statsread_test.go` reads the entries from a journal file and ignores other actions
- [ ] `flai stats --json` on a fixture project carries the counts
- [ ] `go test -race ./internal/metrics ./internal/statsread ./cmd -run Stats` passes
- [ ] Story criterion 4 (the stats' half)

## Notes
