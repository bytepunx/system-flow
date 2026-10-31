---
id: T-2420
type: task
nature: remediation
title: A pack's search for the paths its texts name reads only their slash-bearing runs, so priming the repository takes less time
status: cancelled
parent: S-0548
owner: alex
created: 2026-10-31T09:41:04Z
updated: 2026-10-31T09:49:35Z
transitions:
  - to: cancelled
    at: 2026-10-31T09:49:35Z
    by: agent-S-0548
stream: S-0548
tags: []
touches: [flai/internal/context/design.go, flai/internal/context/written_test.go]
---
# T-2420 A pack's search for the paths its texts name reads only their slash-bearing runs, so priming the repository takes less time

## Work

After T-2361, `go test -race ./internal/mcpserver` took 141 s at load 14 to 31: about 70 s of serial tests and groups, then the two prime-on-this-repository tests, which cannot start before them, 55 to 70 s. T-2361's profile put 48% of their CPU in `Selection.refs` (`flai/internal/context/design.go`), which ran `strings.Contains` of the whole text once per document for each text it reads. Search instead the runs of the text that hold a slash, each a maximal run of bytes some document's path holds, which finds exactly what the whole-text search finds, and the whole text only for a path with no slash. Found in T-2364's measurement; waits for no other task.

## Done when

- [ ] A test compares the found paths with the whole-text search across links, backticks, relative prefixes, a path with a space, and a path with no slash
- [ ] The mcpserver package's time under `-race` is measured again in T-2364

## Notes
- 2026-10-31T09:49:35Z: moved to cancelled: measured no gain: the prime tests used 506 to 587 user-seconds with it and 570 to 589 without under -race; the profile that pointed at refs was taken without -race
