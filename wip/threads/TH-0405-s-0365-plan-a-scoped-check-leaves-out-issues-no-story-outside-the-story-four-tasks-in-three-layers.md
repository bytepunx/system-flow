---
id: TH-0405
title: "S-0365 plan: a scoped check leaves out issues.no-story outside the story, four tasks in three layers"
anchor:
  path: wip/kanban/stories/S-0365-flai-check-finds-issues-no-story-outside-the-story-at-close-out.md
  item: S-0365
status: resolved
participants: [planner-S-0365, orchestrator]
created: 2026-10-12T14:57:34Z
updated: 2026-10-12T14:57:53Z
---

# TH-0405 S-0365 plan: a scoped check leaves out issues.no-story outside the story, four tasks in three layers

On wip/kanban/stories/S-0365-flai-check-finds-issues-no-story-outside-the-story-at-close-out.md.

## Entries

### 2026-10-12T14:57:34Z planner-S-0365
S-0365 is planned. It has 15 touches, a forecast of 23m, and a value of 12.50 USD a week. Its `### Planning` note says where each came from.

## Proposed fix

A check scoped to a story leaves out an `issues.no-story` finding outside the story. A close-out then neither notes it nor records it in I-0134. This is the same fix ADR-0122 made for `item.archive` and ADR-0133 for `board.wip-limit`.

Why:

- All four instances (S-0343, S-0289, S-0342, S-0361) found it on I-0052, which none of those stories touched.
- The finding names an issue, not a story. Only a story made for that issue clears it, and that is your choice.
- It is still reported in two places: the review page lists every open issue no open story links, and `flai check` in the main checkout still warns.

## Tasks

| Layer | Task | Touches |
|-------|------|---------|
| 1 | T-1447: ADR, refining ADR-0085 | `design/adrs`, `design/adrs/README.md` |
| 2 | T-1448: `ScopeToStory` drops it, with tests that reproduce I-0134; `flai check` help and reference | `flai/internal/check/scope.go`, `scope_test.go`, `flai/cmd/check.go`, `check_test.go`, `docs/users/flai-reference.md` |
| 2 | T-1449: design, user guide, work-management convention and its template copy, CHANGELOG | `design/system/flai-cli.md`, `design/system/continuous-improvement.md`, `docs/users/flai.md`, `design/conventions/work-management.md`, `template/root/design/conventions/work-management.md`, `template/CHANGELOG.md` |
| 3 | T-1450: close I-0134 | `design/issues/I-0134-…md`, `design/issues/summary.md` |

T-1448 and T-1449 wait for T-1447 so they can cite the ADR's number. They share no path, so they run together. T-1450 waits for both.

## Assumptions

1. **An `issues.no-story` on an issue file the story changes stays in the check.** An example is an issue the story bumped. This matches `item.archive` and `board.wip-limit`.
2. **The ADR may choose a note instead.** T-1447 may decide on a note that is not recorded in an issue, as ADR-0115 does for `wip.overlap`. That keeps the nudge to make the story in `design/conventions/continuous-improvement.md` in front of each story's agent. I recommend leaving it out, because four close-outs saw it and none made the story.
3. **`design/conventions/continuous-improvement.md` and `flai/internal/harness/harness.go` stay unchanged.** Their rule to make a story when `flai check` warns `issues.no-story` stays true for the check in the main checkout.
4. **I-0052 itself is left alone.** It is still open with no story. Making one is yours, from the review page or with `flai issue story I-0052`.

## Proposals

None to split, merge, or drop. The plan mirrors S-0348, which delivered the same fix with the same four tasks in 24m.

## What you need to do

Review and finalize the draft.

### 2026-10-12T14:57:52Z orchestrator
Approved.

- The four tasks in three layers fit both criteria: the ADR first, then the scoped-check change with I-0134's reproduction beside the docs, then the close.
- The fix follows the rule already recorded for findings a close-out meets outside its story: ADR-0122 (`item.archive`), ADR-0125 (`narrative.state`), and ADR-0133 (`board.wip-limit`). Leaving the finding out rather than noting it fits those. The main checkout's `flai check` and the review page still show I-0052.
- The figures are flai's: 23m, and 12.50 USD a week from the 5m input I set on TH-0404.

One point for the story's agent: T-1447's touch is the folder `design/adrs`, which holds any ready story that touches an ADR while S-0365 is in progress. Once `flai adr new` names the file, narrow the touch to it with `flai touches`.

I-0052 stays the operator's to make a story for.

Source: design/adrs/0133-a-close-out-records-no-board-wip-limit-and-a-check-scoped-to-a-story-leaves-out.md

### 2026-10-12T14:57:53Z orchestrator
Resolved: Plan approved: scoped check leaves out issues.no-story outside the story, as ADR-0122/0125/0133 did for their findings.
