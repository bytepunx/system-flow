---
id: S-0798
type: story
nature: remediation
title: A close-out asked for a story whose branch is behind main stops at sync without syncing, and the failed report makes a story the merge queue could batch not eligible
status: ready
owner: alex
created: 2026-11-16T00:41:49Z
updated: 2026-11-16T01:51:07Z
transitions:
  - to: ready
    at: 2026-11-16T01:47:04Z
    by: orchestrator
tags: [flai, close-out, merge-queue, template]
topics: [close-out, merge-queue]
touches: [scripts/close-out.sh, template/root/scripts/close-out.sh, flai/internal/closeout/script_test.go, flai/internal/serve/verifyjob_container.go, flai/internal/serve/verifyjob_test.go, design/adrs/0309-a-close-out-whose-first-flai-verify-stops-at-sync-syncs-and-verifies-again.md, design/adrs/README.md, design/conventions/work-management.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md, docs/users/flai.md, scripts/README.md, template/root/scripts/README.md, design/issues/I-0383-a-close-out-asked-for-a-story-whose-branch-is-behind-main-stops-at-sync-without.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: planner
      seconds: 416
      estimated: true
      models:
        - model: claude-haiku-5-5
          input: 146
          output: 602
          cache_read: 6707996
          cache_write: 303924
          cost: 0.4775
        - model: claude-opus-5-5
          input: 240
          output: 28913
          cache_read: 17039571
          cache_write: 593740
          cost: 6.2355
    - kind: orchestrator
      seconds: 138
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 35
          output: 290
          cache_read: 1344019
          cache_write: 187046
          cost: 0.4657
cost_of_delay:
  inputs:
    time_lost_per_cycle: 2h
    by: flai
    at: 2026-11-16T00:41:49Z
  value: 300
  by: planner-S-0798
  at: 2026-11-16T01:44:42Z
forecast:
  duration: 1h
  delivery: 2026-11-16T03:08:00Z
  basis: "Its own forecast of 1h; 1st in the pull order with an in-progress limit of 3, behind S-0750, S-0796 and S-0800."
  by: flai
  at: 2026-11-16T01:51:07Z
finalized:
  by: orchestrator
  at: 2026-11-16T01:46:42Z
---
# S-0798 A close-out asked for a story whose branch is behind main stops at sync without syncing, and the failed report makes a story the merge queue could batch not eligible

## Goal

This story remediates [I-0383](../../../design/issues/I-0383-a-close-out-asked-for-a-story-whose-branch-is-behind-main-stops-at-sync-without.md), "A close-out asked for a story whose branch is behind main stops at sync without syncing, and the failed report makes a story the merge queue could batch not eligible". The issue recommends this solution:

The close-out a story is asked for when its branch is behind main syncs first: `scripts/close-out.sh` runs `flai stream sync` when its first `flai verify` stops at sync, as it already does when main moves during the run (ADR-0156), within the same three rounds, and stops with the conflicts when the sync does. Until then a close-out that stops at sync before any tier ran should not replace a passing report as the story's last, or the queue's `verify_failed` should tell that case apart. Bring ADR-0291 §4's words and the `verify_behind_main` blocker's in step with what the close-out does.

## Acceptance criteria
- [ ] The cause I-0383 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0383 is closed with `flai issue close I-0383 --reason` saying what fixed it

## Tasks
- T-3399 scripts/close-out.sh and the template's sync and verify again when the first flai verify stops at sync, within the three rounds
- T-3400 A container story's close-out syncs and verifies again when its first flai verify stops at sync, as the script does
- T-3401 Write an ADR refining ADR-0156 and ADR-0291: a close-out whose first flai verify stops at sync syncs and verifies again within its three rounds
- T-3402 The convention, the guides, and the template's changelog say a close-out asked for a branch behind main syncs first, and I-0383 is closed

## Notes

Cost of delay inputs set by flai from I-0383. time_lost_per_cycle 2h carried over from I-0383's Impact section. Its Impact time_lost_per_cycle was taken rather than the 2h derived from its cost and count.

### Planning

Where each touch came from. The story declared none, so `flai touches suggest` ran from `scripts/close-out.sh` and `flai/internal/serve/verifyjob_container.go`.

- Design, from the goal and I-0383: `scripts/close-out.sh`, whose loop exits when its first `flai verify` fails (line 131) and syncs only after the `--sync-only` check (line 151). Also `design/issues/I-0383-…md` and `design/issues/summary.md`, which criterion 2 closes and regenerates.
- Co-change: `template/root/scripts/close-out.sh` (61%; ADR-0156 §1 changes both scripts). `scripts/README.md` (39%) and `template/root/scripts/README.md` (28%) describe the script's steps. `design/conventions/work-management.md`, `template/root/design/conventions/work-management.md`, and `template/CHANGELOG.md` (17% each) hold the close-out rule and its changelog. `flai/internal/serve/verifyjob_test.go` (11%). The rest were left out: `design/system/devex.md`, `flai/internal/harness/*`, `flai/internal/verify/story.go`, and `system-flow.yaml` say nothing this change alters.
- Layout: `flai/internal/serve/verifyjob_container.go` runs a container story's close-out in Go (ADR-0274), with the same round loop and the same gap (line 144). `flai/internal/closeout/script_test.go` is new: no test runs the script today.
- Design, the ADR the goal asks for: `design/adrs/0309-…md`, the next number after S-0796's ADR-0308, its slug from T-3401's title; and `design/adrs/README.md`, the index `flai adr new` writes. The number may change if another story takes 0309 first. `flai task done` widens the touches to the file committed.
- Design: `docs/users/flai.md`, whose close-out paragraph (line 963) says the close-out refuses a branch behind main and says to sync.
- Not touched: `flai/internal/preview/accept.go` and `flai/cmd/accept.go`. The `verify_behind_main` blocker and the help text already say the close-out syncs the branch onto main, and they hold once it does. ADR-0291 is accepted and immutable: the new ADR refines it.

No folder touch, and no `after`: the story reads nothing another open story writes. S-0796's ADR-0308 only sets the next ADR number, which is no output this story reads.

Figures:

- Forecast 1h, raised from `flai forecast`'s 33m. The fix spans two shell scripts and a Go loop. The script needs a stub harness for its first test. A story that changes `flai/` and `template/` closes out with the integration and smoke tiers, about twelve minutes (ADR-0291, Consequences). Delivery moved by the same 27 minutes, from 04:29Z to 04:56Z.
- Cost of delay 300 USD a week, as `flai cod` gives it from the input flai carried from I-0383: 2h a cycle at 150 USD an hour, one cycle a week. It stands. The input is the operator's, and I-0383 has one recorded occurrence, which held three stories on one day. That is no basis for more cycles a week.
