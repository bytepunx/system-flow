---
id: S-0610
type: story
nature: improvement
title: The close-out runs independent test tiers at once and smoke skips what the close-out already ran, saving about 40 to 150 s a close-out
status: backlog
parent: E-0033
owner: alex
created: 2026-11-07T22:42:03Z
updated: 2026-11-07T22:42:03Z
transitions: []
tags: [cli]
touches: [flai/internal/verify/run.go, flai/internal/verify/story.go, scripts/smoke.sh]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
---
# S-0610 The close-out runs independent test tiers at once and smoke skips what the close-out already ran, saving about 40 to 150 s a close-out

## Goal

Cause 20 in [server-performance.md](../../../design/system/server-performance.md#where-the-time-goes): `flai verify` runs its tiers one after another and stops at the first failure (`flai/internal/verify/run.go`). The `flaiover` tier, which is node, waits for `integration`, which is Go, and `smoke` waits for both. `smoke` (`scripts/smoke.sh`) runs the template test, the check, and the markdown lint again after the close-out's own tiers ran them. Passed close-outs in the week to 2026-11-07 took p50 247 s, p90 510 s, max 1,408 s; `integration` mean 238 s, `flaiover` mean 176 s on 12 of 60 jobs, `smoke` mean 44 s.

## Acceptance criteria

- [ ] `flai verify` runs at once the selected tiers the manifest marks as independent of each other, cheapest first as now, and cancels the others when one fails, so the verdict and the step it names are the same as a serial run's.
- [ ] Under `flai verify`, `smoke` skips the steps a tier of the same run already passed.
- [ ] A close-out of a story that selects `integration`, `flaiover`, and `smoke` takes at least 150 s less than its serial run on the same commit, and server-performance.md records the before and after.

## Tasks

## Notes

- How a tier says it may run beside others is a manifest key, which `project-manifest.md` and `docs/operators/settings.md` would gain; that is likely an ADR.
- Splitting `internal/mcpserver`'s tests, 337 of 1,676 package-seconds in the 2026-11-02 batch, would shorten `integration` itself (I-0234).
- Turning the `merge-queue` host action on leaves `integration` and `smoke` to the batch; that is the operator's choice and needs no story.
