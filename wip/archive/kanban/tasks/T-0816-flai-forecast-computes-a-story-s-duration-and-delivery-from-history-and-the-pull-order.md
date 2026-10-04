---
id: T-0816
type: task
nature: feature
title: flai forecast computes a story's duration and delivery from history and the pull order
status: done
parent: S-0210
owner: alex
created: 2026-10-04T19:45:44Z
updated: 2026-10-04T19:59:35Z
transitions:
  - to: ready
    at: 2026-10-04T19:46:38Z
    by: agent-S-0210
  - to: in-progress
    at: 2026-10-04T19:46:39Z
    by: agent-S-0210
  - to: done
    at: 2026-10-04T19:59:35Z
    by: agent-S-0210
stream: S-0210
tags: [flai]
touches: [flai/internal/manifest, flai/internal/workitem/planning.go, flai/internal/workitem/planning_test.go, flai/internal/planning/doc.go, flai/internal/planning/forecast.go, flai/internal/planning/forecast_test.go, flai/internal/planning/testdata, flai/cmd/forecast.go, flai/cmd/forecast_test.go, flai/cmd/root.go, docs/operators/settings.md]
usage:
  source: log
  seconds: 776
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 78
      output: 29645
      cache_read: 3855004
      cache_write: 110293
      cost: 2.0153
---
# T-0816 flai forecast computes a story's duration and delivery from history and the pull order

## Work

Add `planning.default_duration` to the manifest (a Go duration, default 1h, validated as `manifest.planning`), and `workitem.CriteriaCount`. Add the `planning` package (`doc.go` holds its package comment) with `Forecast`: a story's size is its criteria count plus its touches count; its history is the done stories with `usage.seconds`, matched on nature, model, and size band, widening to nature and size, then nature, then all, until at least three match; its duration is the median seconds per unit of size times its size; the medians of in-progress time and cycle time over agent seconds turn a duration into the time a story holds a lane and the time to its delivery. Delivery simulates the pull order: lanes up to the in-progress limit, freed by the stories in progress, then the ready and backlog stories ahead of it in order, each with its own forecast duration when it has one, starting no sooner than the stories it waits for (`after`). Without enough history the duration is `planning.default_duration` and the basis says so. `flai forecast S-nnnn [--json]` prints duration, delivery, and a one-sentence basis, with the history, position, and schedule in JSON. A settings row for the new key.

Waits for nothing: it shares no path with the touches task and runs beside it in its own task worktree.

## Done when

- [ ] `flai forecast S-nnnn` prints a Go duration, a UTC delivery, and a basis; `--json` adds the history, size, position, limit, and the stories ahead with their start and delivery
- [ ] Fixture tests pin the duration and delivery for a story with history, one ahead of a full in-progress limit, one waiting on `after`, and one without history using `planning.default_duration`
- [ ] `planning.default_duration` is validated, documented in `docs/operators/settings.md`, and its manifest test passes

## Notes
