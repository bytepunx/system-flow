---
id: S-0683
type: story
nature: remediation
title: "Two stories accepted apart leave main red: S-0535 reworded to-backlog's condition, and S-0416's TestADefinedRoleIsHeldToItsFile still expects the old words"
status: backlog
owner: alex
created: 2026-11-10T08:37:08Z
updated: 2026-11-11T18:15:37Z
transitions: []
tags: [flai, acceptance]
topics: [review, code]
touches: [design/adrs, design/adrs/README.md, flai/internal/verify/story.go, flai/internal/verify/story_test.go, flai/internal/preview/accept.go, flai/internal/preview/accept_test.go, flai/cmd/accept.go, flai/cmd/accept_stale_verify_test.go, flai/internal/guard/roles_test.go, design/system/workflow.md, design/system/flai-cli.md, docs/users/flai.md, design/issues/I-0229-two-stories-accepted-apart-leave-main-red-s-0535-reworded-to-backlog-s-condition-and-s-0416-s-testadefinedroleisheldtoitsfile-still-expects-the-old-words.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
cost_of_delay:
  inputs:
    time_lost_per_cycle: 18m
    by: flai
    at: 2026-11-10T08:37:08Z
  value: 45
  by: planner-S-0683
  at: 2026-11-10T10:34:43Z
forecast:
  duration: 2h
  delivery: 2026-11-12T03:45:00Z
  basis: "Its own forecast of 2h; 15th in the pull order with an in-progress limit of 3, behind S-0715, S-0707, S-0714, S-0711, S-0712, S-0720, S-0681, S-0718, S-0241, S-0659, S-0660, S-0661, S-0662 and S-0663."
  by: flai
  at: 2026-11-11T18:15:37Z
---
# S-0683 Two stories accepted apart leave main red: S-0535 reworded to-backlog's condition, and S-0416's TestADefinedRoleIsHeldToItsFile still expects the old words

## Goal

This story remediates [I-0229](../../../design/issues/I-0229-two-stories-accepted-apart-leave-main-red-s-0535-reworded-to-backlog-s-condition-and-s-0416-s-testadefinedroleisheldtoitsfile-still-expects-the-old-words.md), "Two stories accepted apart leave main red: S-0535 reworded to-backlog's condition, and S-0416's TestADefinedRoleIsHeldToItsFile still expects the old words". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0229 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0229 is closed with `flai issue close I-0229 --reason` saying what fixed it

## Tasks
- T-3051 ADR: an acceptance waits for a verify over main's new paths
- T-3053 TestADefinedRoleIsHeldToItsFile reads the conditions' Allows from internal/roles instead of copying them
- T-3058 The acceptance preview lists the blocker verify_stale when main's new paths send back a tier the story's last verify ran
- T-3060 flai accept refuses a story whose verify main has outrun, with a real-git test that reproduces I-0229
- T-3062 Document that an acceptance waits for a verify over main's new paths
- T-3065 Close I-0229 with what fixed it

## Notes

Cost of delay inputs set by flai from I-0229. time_lost_per_cycle 18m: 30m per occurrence × 1 occurrence ÷ 1.7 cycles of 168h (first reported 2026-10-29T03:27:16Z, 12.2 days before this story).

### Planning

Cause found by the planner's explorer, proposed in the plan thread. S-0416's close-out passed at 6c25aad7 over base 1dfa57c8, before S-0535 landed. The operator then accepted S-0416 by hand (the queue had refused it for protected paths, TH-0722). `flai accept` rebased it onto S-0535's main and fast-forwarded. No tier ran at that rebased head. Tier selection was not the gap: `flai/internal/verify/importers.go` makes `internal/guard` an importer of `internal/roles`, so a verify at the new head would have run the guard test (ADR-0156's sent-back tiers, `flai/internal/verify/story.go` `sentBack`). Acceptance never runs tiers. `queueTiersGate` in `flai/internal/preview/accept.go` compares the report only to the branch's own head, never to how far main moved. Proposed fix: every acceptance but the queue's landing is refused while main has brought in paths, since the base the story's last verify ran over, that send back a tier that verify ran. The blocker names the close-out to ask for again, which re-runs only those tiers. Also, the guard test stops copying the condition text and reads it from `internal/roles`.

Touches, and where each came from (`flai touches suggest` had nothing to start from, since the story declared none; once they were set, its co-change list named no path above 10% that the plan needs):

- `design/adrs` (folder, kept) and `design/adrs/README.md`: design. The decision is a new ADR. Its number is assigned when it is written, and other open stories may take 0279 first, so no task can name the file yet. The folder is in `claims.shared`, so it holds no ready story.
- `flai/internal/verify/story.go`, `flai/internal/verify/story_test.go`: layout. The sent-back selection and the report's base live there.
- `flai/internal/preview/accept.go`, `flai/internal/preview/accept_test.go`: layout. Every acceptance's blockers are computed there, beside `queueTiersGate`.
- `flai/cmd/accept.go`, `flai/cmd/accept_stale_verify_test.go` (new): layout. These hold the refusal and the end-to-end test that reproduces I-0229: two stories, one rewording a string and one whose test copies it, accepted one after the other.
- `flai/internal/guard/roles_test.go`: the issue. This is the test that went red; it still holds a literal copy of `Allows`.
- `design/system/workflow.md` § What waits on an acceptance, `design/system/flai-cli.md` § The acceptance gate, `docs/users/flai.md`: design and docs. They describe what an acceptance waits for (`documentation.md`).
- The I-0229 issue file and `design/issues/summary.md`: criterion 2.

No `after`: the story reads no other open story's output.

Figures:

- Forecast: 2h, raised from `flai forecast`'s 29m. That figure is 106 s per unit of size over 64 large remediation stories, times size 16. The story writes an ADR, a new acceptance blocker across `preview` and `cmd` with a real-git end-to-end test, and three documents. Every `flai/` close-out also runs the integration and smoke tiers. The delivery, 18:43Z, is flai's 17:12Z moved by the same 91 minutes.
- Cost of delay: 45 USD a week, as `flai cod` worked it out from the operator's input of 18m lost per 168h cycle at 150 USD an hour. It stands as computed. Each recurrence also fails an innocent story in the merge queue (S-0519 here), but that cost is already in the issue's 30m per occurrence, and the inputs are the operator's.
