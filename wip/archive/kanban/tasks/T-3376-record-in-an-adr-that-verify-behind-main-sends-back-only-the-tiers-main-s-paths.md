---
id: T-3376
type: task
nature: remediation
title: Record in an ADR that verify_behind_main sends back only the tiers main's paths overlapping the story's change set select
status: done
parent: S-0792
owner: alex
created: 2026-11-15T07:43:13Z
updated: 2026-11-15T16:45:27Z
transitions:
  - to: ready
    at: 2026-11-15T16:44:39Z
    by: agent-S-0792
  - to: in-progress
    at: 2026-11-15T16:44:39Z
    by: agent-S-0792
  - to: done
    at: 2026-11-15T16:45:27Z
    by: agent-S-0792
stream: S-0792
tags: [adr, acceptance]
touches: [design/adrs/README.md, design/adrs/0303-an-acceptance-s-verify-behind-main-sends-back-only-the-tiers-selected-by-the.md]
usage:
  source: log
  seconds: 48
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 15
      output: 4709
      cache_read: 1137473
      cache_write: 51373
      cost: 0.6894
  strategic:
    - kind: orchestrator
      seconds: 5
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 1
          output: 7
          cache_read: 21939
          cache_write: 603
          cost: 0.0068
---
# T-3376 Record in an ADR that verify_behind_main sends back only the tiers main's paths overlapping the story's change set select

## Work

Write an ADR with `flai adr new` that refines ADR-0291, and accept it. It decides that the `verify_behind_main` blocker counts only main's new paths that overlap the story's change set, the paths its last passing verify recorded (`Report.Paths`):

- a path the story changes too; or
- a path in a Go package that a package the story changes imports, directly or through others, as the `importers` closure in `flai/internal/verify/importers.go` finds it.

Only the tiers those overlapping paths select are sent back. Main's paths that overlap nothing send back nothing, whatever tiers they select, so an acceptance that changes another part of flai no longer holds every story in review. The I-0229 case stays blocked. Name I-0378 and TH-1213, TH-1215, and TH-1216 in its Context, and the operator's criteria on S-0792. List file overlap alone as an alternative considered, rejected because it reopens I-0229, unless the operator answers otherwise on the plan's thread.

This task waits for none: the code and docs tasks build on the decision it records.

## Done when

- The ADR is accepted, refines ADR-0291, and states the overlap rule, what is sent back, and that the I-0229 case stays blocked.
- `design/adrs/README.md` lists it, and `flai check --strict` is clean.

## Notes
