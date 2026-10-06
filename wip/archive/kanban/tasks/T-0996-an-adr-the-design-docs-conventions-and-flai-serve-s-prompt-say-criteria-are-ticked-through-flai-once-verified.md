---
id: T-0996
type: task
nature: remediation
title: An ADR, the design, docs, conventions, and flai serve's prompt say criteria are ticked through flai once verified
status: done
parent: S-0282
owner: alex
created: 2026-10-06T03:47:59Z
updated: 2026-10-06T04:56:02Z
transitions:
  - to: ready
    at: 2026-10-06T04:39:28Z
    by: agent-S-0282
  - to: in-progress
    at: 2026-10-06T04:39:28Z
    by: agent-S-0282
  - to: done
    at: 2026-10-06T04:56:02Z
    by: agent-S-0282
stream: S-0282
tags: [cli]
touches: [design/adrs, design/system/flai-cli.md, design/system/workflow.md, design/system/flaiover-dashboard.md, docs/users/flai.md, docs/users/flai-reference.md, design/conventions/work-management.md, design/conventions/delegation.md, template/root/design/conventions/work-management.md, template/root/design/conventions/delegation.md, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, template/CHANGELOG.md]
after: [T-0992, T-0993, T-0994, T-0995]
usage:
  source: log
  seconds: 994
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 148
      output: 43556
      cache_read: 6976465
      cache_write: 205337
      cost: 3.5033
---
# T-0996 An ADR, the design, docs, conventions, and flai serve's prompt say criteria are ticked through flai once verified

## Work

Record the decision in an ADR with `flai adr new`, and say how criteria are ticked in the living design (`design/system/flai-cli.md`, `design/system/workflow.md`, `design/system/flaiover-dashboard.md`), the user docs (`docs/users/flai.md`, `docs/users/flai-reference.md` by `make flai-reference`), the conventions (`work-management.md` and `delegation.md`, template baseline first, then copied), flai serve's start prompt in `flai/internal/harness/harness.go`, and `template/CHANGELOG.md`. Who ticks follows TH-0162's answer.

Waits for every other task: it documents what they built, and its wording follows TH-0162.

## Done when

- [ ] The ADR, design, docs, conventions, and prompt say that a criterion is ticked through flai once verified, and by whom.
- [ ] `go test -race ./internal/harness/` passes and `make flai-reference` leaves the reference current.

## Notes
