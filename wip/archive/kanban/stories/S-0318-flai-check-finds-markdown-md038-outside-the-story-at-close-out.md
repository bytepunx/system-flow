---
id: S-0318
type: story
nature: improvement
title: "flai check finds `markdown.MD038` outside the story at close-out"
status: done
owner: alex
created: 2026-10-07T18:59:49Z
updated: 2026-10-08T05:53:29Z
transitions:
  - to: ready
    at: 2026-10-08T00:11:45Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-08T04:34:02Z
    by: agent-S-0318
  - to: review
    at: 2026-10-08T05:52:14Z
    by: agent-S-0318
  - to: done
    at: 2026-10-08T05:53:29Z
    by: orchestrator
tags: [flai, template]
topics: [cli, conventions, template]
touches: [design/adrs, flai/internal/check/scope.go, flai/internal/check/scope_test.go, flai/cmd/check.go, flai/cmd/check_test.go, docs/users/flai-reference.md, design/conventions/work-management.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md, design/system/continuous-improvement.md, design/system/flai-cli.md, docs/users/flai.md, design/issues/I-0096-flai-check-finds-markdown-md038-outside-the-story-at-close-out.md, design/issues/summary.md, design/issues/I-0118-flai-check-finds-markdown-md034-outside-the-story-at-close-out.md, design/issues/I-0112-flai-verify-record-issues-opens-an-issue-on-the-story-branch-that-another-branch-opened-under-the-same-title-meanwhile.md, design/issues/I-0117-a-story-s-integration-tier-lints-main-s-committed-wip-so-markdown-any-agent-commits-there-fails-every-story-s-close-out-until-main-commits-a-fix.md, design/issues/I-0086-the-close-out-s-install-smoke-test-failed-once-and-passed-when-run-alone-with-no-cause-in-its-output.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 2932
  turns:
    - day: 2026-10-08
      ceremony: 3
      hand_edits: 2
      work: 69
  models:
    - model: claude-opus-5-5
      input: 224
      output: 58494
      cache_read: 16979764
      cache_write: 370306
      cost: 7.1677
  strategic:
    - kind: orchestrator
      seconds: 1746
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 101
          output: 1435
          cache_read: 25116246
          cache_write: 69147
          cost: 6.2055
        - model: claude-sonnet-5-5
          input: 16
          output: 85
          cache_read: 230072
          cache_write: 45173
          cost: 0.2373
cost_of_delay:
  inputs:
    time_lost_per_cycle: 15m
    by: orchestrator
    at: 2026-10-08T00:06:21Z
  value: 37.5
  by: planner-S-0318
  at: 2026-10-08T00:10:05Z
forecast:
  duration: 25m
  delivery: 2026-10-08T05:01:00Z
  basis: "Its own forecast of 25m; 1st in the pull order with an in-progress limit of 3, behind S-0232, S-0316 and S-0324."
  by: flai
  at: 2026-10-08T04:33:38Z
finalized:
  by: orchestrator
  at: 2026-10-08T00:11:41Z
---
# S-0318 flai check finds `markdown.MD038` outside the story at close-out

## Goal

This story remediates [I-0096](../../../design/issues/I-0096-flai-check-finds-markdown-md038-outside-the-story-at-close-out.md), "flai check finds `markdown.MD038` outside the story at close-out". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0096 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0096 is closed with `flai issue close I-0096 --reason` saying what fixed it

## Tasks
- T-1300 An ADR refining ADR-0085 and ADR-0122 records the remedy for I-0096, proposed from its instance
- T-1301 A check scoped to a story leaves out a markdown finding on another open story's narrative, so a close-out records none
- T-1302 The convention, the design, and the users' guide say that a close-out leaves out a markdown finding on another open story's narrative
- T-1303 I-0096 is closed with flai issue close, saying that a close-out leaves out a markdown finding on another open story's narrative

## Notes

### Planning

The proposed remedy, which T-1300's ADR decides: a check scoped to a story leaves out a `markdown.*` finding on the narrative of another open story, as ADR-0122 does for `item.archive`. I-0096's one instance was S-0229's hand-edited `## Decisions`, which flai's wip lint guard never sees. S-0227's close-out recorded it, though only S-0229 could fix it, and S-0229's own close-out checks its own narrative.

Where each touch came from. `flai touches suggest` listed nothing, because the story declared no touches.

| Touch | Source |
|-------|--------|
| `design/adrs` | layout: the ADR T-1300 writes |
| `flai/internal/check/scope.go`, `scope_test.go` | design: `ScopeToStory` marks outside findings and drops `item.archive` (ADR-0122) |
| `flai/cmd/check.go`, `check_test.go` | design: `recordOutside` and the command's help |
| `docs/users/flai-reference.md` | layout: generated from the help by `make flai-reference` |
| `design/conventions/work-management.md`, `template/root/design/conventions/work-management.md`, `template/CHANGELOG.md` | co-change: S-0280 changed all three for `item.archive` |
| `design/system/continuous-improvement.md`, `design/system/flai-cli.md`, `docs/users/flai.md` | co-change: S-0280's T-1167 changed these for the same rule |
| `design/issues/I-0096-…md`, `design/issues/summary.md` | criterion 2: `flai issue close` writes both |

One folder touch is kept: `design/adrs`. `flai adr new` allocates the ADR's number when it runs, so no file can be named now.

Figures:

- Forecast: 25m, delivery 2026-10-08T06:54:00Z. `flai forecast` gave 21m (78 s per unit of size times 16). Raised to 25m because S-0280 and S-0279, the same shape of remedy, took 22m and 26m, and this one must also tell an open story's narrative from a closed one's.
- Cost of delay: 37.5 USD a week, from `flai cod` on the input `time_lost_per_cycle: 15m`. The orchestrator set that input on TH-0343, as recommended there: one instance, against 45m for S-0279's 18 and S-0280's 38. Left as `flai cod` gives it.

### Accepted by the orchestrator

- Verified: 670bd2a2eb76f84aa126d9228c7f5f8dfabf84da
- At: 2026-10-08T05:53:29Z

Verdict: accept. flai verify passed every step at the branch head 670bd2a2, and the verifier matched both criteria to the diff.
- 1: design/adrs/0123-a-close-out-records-no-markdown-finding-on-another-open-story-s-narrative-and-a.md, design/adrs/README.md, flai/internal/check/scope.go, flai/internal/check/scope_test.go, flai/cmd/check.go, flai/cmd/check_test.go, docs/users/flai-reference.md, design/conventions/work-management.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md, design/system/continuous-improvement.md, design/system/flai-cli.md, docs/users/flai.md
- 2: design/issues/I-0096-flai-check-finds-markdown-md038-outside-the-story-at-close-out.md, design/issues/summary.md
