---
id: S-030
type: story
nature: feature
title: Logging and telemetry conventions
status: review
parent: E-005
owner: alex
created: 2026-09-16T23:18:16Z
updated: 2026-09-16T23:25:10Z
transitions:
  - to: ready
    at: 2026-09-16T23:23:14Z
    by: agent
  - to: in-progress
    at: 2026-09-16T23:23:14Z
    by: agent
  - to: review
    at: 2026-09-16T23:25:10Z
    by: agent
tags: [conventions]
---

# S-030 Logging and telemetry conventions

## Goal
The two draft conventions become real: `logging.md` says what is logged, at which level, in what shape, with what correlation, and what never appears in a log; `telemetry.md` says which metrics, traces, and health signals every service emits, how they are named and labelled, and where they go. Both are written in the template, agreed with the operator, and mirrored here.

## Acceptance criteria
- [ ] template/root/design/conventions/logging.md has Rules and When in doubt written, `status: active`, under 120 lines
- [ ] template/root/design/conventions/telemetry.md the same
- [ ] README rows drop the "(draft)" suffix in both copies; design/conventions here matches the template above the markers
- [x] Project additions here say how flai and flaiover apply them (flai: structured stderr logging and a --verbose flag; flaiover: request logs, health endpoint, metrics endpoint) and open follow-up stories where code must change
- [x] design/tech records any logging or telemetry library the conventions imply, with version and rationale, or states that the standard library suffices
- [x] flai check --strict passes in this repository and on a rendered template

## Tasks
- T-054 Draft logging.md rules
- T-055 Draft telemetry.md rules
- T-056 Project additions, follow-up stories, design/tech, mirror here

## Notes
- Rules drafted; the first three criteria (status active, README suffix) wait for the operator to confirm the text.
- Created 2026-09-16 at the operator's request after S-025 found the two empty placeholder files. The agent proposes the rules; the operator edits or confirms before they leave draft.
- Proposed starting points: structured logs (one event per line, key-value, UTC timestamps, level, component, correlation ID), levels limited to error, warn, info, debug with meanings, no secrets or personal data, sampling for high-volume debug; telemetry as OpenTelemetry-shaped metrics and traces with a fixed naming scheme, a health endpoint per service, and the four golden signals per service as the minimum set.
