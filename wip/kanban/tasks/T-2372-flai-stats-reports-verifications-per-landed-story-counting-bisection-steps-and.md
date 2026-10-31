---
id: T-2372
type: task
nature: improvement
title: flai stats reports verifications per landed story, counting bisection steps, and a clean half lands without a verification of its own
status: in-progress
parent: S-0552
owner: alex
created: 2026-10-31T07:33:12Z
updated: 2026-10-31T07:35:00Z
transitions:
  - to: ready
    at: 2026-10-31T07:34:59Z
    by: system-flow
  - to: in-progress
    at: 2026-10-31T07:35:00Z
    by: system-flow
stream: S-0552
tags: [flai]
touches: [flai/internal/mergequeue/bisect.go, flai/internal/mergequeue/bisect_test.go, flai/internal/mergequeue/log.go, flai/internal/mergequeue/log_test.go, flai/internal/mergequeue/land_test.go, flai/internal/metrics/queue.go, flai/internal/metrics/queue_test.go, flai/internal/statsread/statsread.go, flai/cmd/stats.go, flai/cmd/stats_test.go, design/system/metrics.md]
---
# T-2372 flai stats reports verifications per landed story, counting bisection steps, and a clean half lands without a verification of its own

## Work

Log each bisection step's verification to the queue's event log (`.flai-cache/merge-queue/events.jsonl`) as the batch's own is, so a verification run is countable. `metrics/queue.go` adds per week the verifications run (batches and bisection steps) and verifications per landed story, beside the batch size it reports; `flai stats` prints them and `--json` carries them; `design/system/metrics.md`'s merge queue section defines them, with the ADR the docs task writes. Add a test, if none exists, that a story in a failed batch's clean half lands through `flai accept --batch` with no verification of its own. Waits for nothing: it shares no file with the other layer-1 tasks.

## Done when

- [ ] Tests cover the step events, the two figures, and the clean half's landing without a verification.

## Notes

Layer 1.
