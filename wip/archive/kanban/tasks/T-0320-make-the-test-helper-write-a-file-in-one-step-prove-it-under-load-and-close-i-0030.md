---
id: T-0320
type: task
nature: feature
title: Make the test helper write a file in one step, prove it under load, and close I-0030
status: done
parent: S-0093
owner: alex
created: 2026-09-21T22:49:15Z
updated: 2026-09-21T22:55:02Z
transitions:
  - to: ready
    at: 2026-09-21T22:50:18Z
    by: system-flow
  - to: in-progress
    at: 2026-09-21T22:50:18Z
    by: system-flow
  - to: review
    at: 2026-09-21T22:55:01Z
    by: system-flow
  - to: done
    at: 2026-09-21T22:55:02Z
    by: system-flow
stream: S-0093
tags: []
---
# T-0320 Make the test helper write a file in one step, prove it under load, and close I-0030

## Work
`write` in `watch_test.go` uses `os.WriteFile`, which truncates and then writes, so a tick landing between the two sees a zero-length file and a stalled test goroutine can make that a second report. Write to a scratch folder outside the watched root and rename into place. Then run `go test -race -short -count=200 ./internal/watch` on this host, once idle and once with every core busy (busy processes stopped by exact PID), and close I-0030 with this story as its remediation.

## Done when
Both 200-run passes are clean and I-0030 is closed with the remediation recorded.

## Notes
The helper and the debounce are the same defect at two levels: something the test does not control changing what the watcher sees between two ticks.
