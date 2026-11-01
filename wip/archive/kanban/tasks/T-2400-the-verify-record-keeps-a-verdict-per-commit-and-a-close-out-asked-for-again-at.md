---
id: T-2400
type: task
nature: improvement
title: The verify record keeps a verdict per commit, and a close-out asked for again at a commit it holds answers from it, saying so, and runs nothing
status: done
parent: S-0546
owner: alex
created: 2026-10-31T07:43:44Z
updated: 2026-10-31T09:29:04Z
transitions:
  - to: ready
    at: 2026-10-31T08:52:15Z
    by: agent-S-0546
  - to: in-progress
    at: 2026-10-31T08:52:17Z
    by: agent-S-0546
  - to: done
    at: 2026-10-31T09:29:03Z
    by: agent-S-0546
stream: S-0546
tags: [flai]
touches: [flai/internal/verify/story.go, flai/internal/verify/story_test.go, flai/internal/verify/record.go, flai/internal/verify/record_test.go]
after: [T-2396]
usage:
  source: log
  seconds: 2206
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 55
      output: 22354
      cache_read: 3796594
      cache_write: 194854
      cost: 2.263
---
# T-2400 The verify record keeps a verdict per commit, and a close-out asked for again at a commit it holds answers from it, saying so, and runs nothing

## Work

In `flai/internal/verify/story.go`, keep the verdict of each commit the close-out verified, as T-2396's ADR decides. Today `Report` is one per story, saved by `SaveReport` and read by `LastReport` at `ReportPath`, and `reusable` reuses passed tiers but always runs. Then make `Verify` answer from the record, and run nothing, when the record holds a verdict for the commit (and base and fingerprint, if the ADR says so) it is asked at. The answer says it came from the record and when that run was. `--fresh` still runs. Keep `LastReport`'s readers (`flai verify --last`, the verifier, the dashboard's `verify.status`) working on the newest run.

Waits for T-2396, which decides the key and which commit a second ask matches. Shares no path with T-2403, so the two run together.

## Done when

- [ ] A test in `flai/internal/verify/story_test.go` verifies a story, asks again at the same commit, and gets the stored verdict with nothing run and the answer saying it came from the record
- [ ] A test shows a new commit, or `--fresh`, runs again
- [ ] `go test -race ./internal/verify` passes
- [ ] Story criterion 2

## Notes
