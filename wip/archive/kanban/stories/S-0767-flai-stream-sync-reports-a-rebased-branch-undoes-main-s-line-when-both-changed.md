---
id: S-0767
type: story
nature: remediation
title: flai stream sync reports a rebased branch undoes main's line when both changed one long markdown line and the branch keeps every word of main's
status: done
owner: alex
created: 2026-11-13T11:42:36Z
updated: 2026-11-13T12:53:44Z
transitions:
  - to: ready
    at: 2026-11-13T12:03:56Z
    by: alex
  - to: in-progress
    at: 2026-11-13T12:04:05Z
    by: agent-S-0767
  - to: review
    at: 2026-11-13T12:51:48Z
    by: agent-S-0767
  - to: done
    at: 2026-11-13T12:53:44Z
    by: orchestrator
tags: []
touches: [flai/internal/storygit/undone.go, flai/internal/storygit/undone_test.go, flai/internal/storygit/sync_test.go, design/adrs/README.md, design/system/flai-cli.md, design/system/workflow.md, docs/users/flai.md, design/adrs/0297-a-markdown-line-of-main-s-that-a-rebased-branch-removes-is-not-undone-when-a.md, design/issues/I-0360-flai-stream-sync-reports-a-rebased-branch-undoes-main-s-line-when-both-changed.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 2958
  priced_by: harness
  turns:
    - day: 2026-11-13
      ceremony: 1
      work: 68
  models:
    - model: claude-opus-5-5
      input: 242
      output: 61956
      cache_read: 20181330
      cache_write: 489799
      cost: 8.647
  strategic:
    - kind: planner
      seconds: 172
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 42
          output: 14900
          cache_read: 2764283
          cache_write: 169957
          cost: 2.2107
    - kind: orchestrator
      seconds: 164
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 51
          output: 510
          cache_read: 1618601
          cache_write: 176943
          cost: 0.5199
        - model: claude-sonnet-5-5
          input: 7
          output: 49
          cache_read: 44321
          cache_write: 37907
          cost: 0.0635
cost_of_delay:
  inputs:
    time_lost_per_cycle: 10m
    by: flai
    at: 2026-11-13T11:42:36Z
  value: 25
  by: planner-S-0767
  at: 2026-11-13T11:47:57Z
forecast:
  duration: 40m
  delivery: 2026-11-13T23:33:00Z
  basis: "flai forecast's 23m (median 136 s per unit over 64 done medium remediation stories on claude-opus-5-5, size 10) raised to 40m for an ADR, three long design and docs passages, and a real-git rebase-stop test beyond the median remediation; delivery moved by the same 17m from flai forecast's 23:16Z, 30th in the pull order."
  by: planner-S-0767
  at: 2026-11-13T11:47:57Z
finalized:
  by: orchestrator
  at: 2026-11-13T11:48:33Z
---
# S-0767 flai stream sync reports a rebased branch undoes main's line when both changed one long markdown line and the branch keeps every word of main's

## Goal

This story remediates [I-0360](../../../design/issues/I-0360-flai-stream-sync-reports-a-rebased-branch-undoes-main-s-line-when-both-changed.md), "flai stream sync reports a rebased branch undoes main's line when both changed one long markdown line and the branch keeps every word of main's". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0360 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0360 is closed with `flai issue close I-0360 --reason` saying what fixed it

## Tasks
- T-3291 Count a markdown line of main's as kept when a line the branch adds keeps every word of it in order
- T-3292 Reproduce I-0360 with a real-git sync that stops on one long markdown line and resumes with main's words kept
- T-3293 Record the kept-words rule in an ADR refining ADR-0155, the design, and the user guide, and close I-0360

## Notes

Cost of delay inputs set by flai from I-0360. time_lost_per_cycle 10m: 10m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-13T09:21:49Z, 0.1 days before this story; under one cycle counts as one).

### Planning

Proposed solution, from I-0360's one instance (S-0761 on `design/system/flai-cli.md`): `storygit.FindUndone` and `RecordedUndone` (`flai/internal/storygit/undone.go`) count a line main added as undone when the branch's net change removes a line of the same text (`undoneLines`, a multiset of whole lines, ADR-0155). In a markdown file, a removed line of main's is not undone when a line the branch adds in the same path keeps every word of it, in order: main's line's words, split on white space, are a subsequence of the added line's words. Each added line answers for one of main's lines at most. That tells an extension of main's line, as S-0761's resolution of the stop was, from an undo, and leaves ADR-0155's revert-replayed case as it is, since there the branch adds nothing that keeps main's words. It is limited to `.md` paths, as ADR-0287's sentence merge is to markdown, since in code a kept token sequence can still undo main's change (a line commented out keeps every word). It is a decision that refines ADR-0155, so the story writes an ADR for it.

Touches, and where each came from:

| Touch | Source | Why |
|-------|--------|-----|
| `flai/internal/storygit/undone.go` | design (I-0360, ADR-0155) | `FindUndone`, `RecordedUndone`, and `undoneLines` hold the check |
| `flai/internal/storygit/undone_test.go` | co-change (2 of 4 commits) | unit tests of the new rule |
| `flai/internal/storygit/sync_test.go` | co-change (2 of 4 commits) | the real-git reproduction of I-0360 through a stopped and resumed sync |
| `design/adrs/0297-a-markdown-line-a-rebased-branch-changes-keeping-every-word-of-main-s-line-in.md` | design | the ADR refining ADR-0155; its number and slug are predicted and `flai task done` corrects them |
| `design/adrs/README.md` | layout | the generated ADR index `flai adr new` writes |
| `design/system/flai-cli.md` | design | the `flai stream sync` row says the compare is a multiset of line text |
| `design/system/workflow.md` | design | § Branches and collisions describes the S-0371 compare |
| `docs/users/flai.md` | docs | the "Main's change undone" bullet says what sync warns on |

`flai/internal/storygit/sync.go` co-changed in 3 of 4 commits but is left out: `checkUndone` and `withRecordedUndone` call `FindUndone` and `RecordedUndone` and need no change. `docs/users/flai-reference.md` is left out: it is generated from the commands' help, which stays true. The merge queue (`internal/mergequeue/eligible.go`) and the acceptance preview (`internal/preview/accept.go`) read `RecordedUndone` and take the fix with it. The closing of I-0360 writes under `design/issues/`, which flai writes and the touches drift leaves out (ADR-0284). No folder touch is kept. No `after`: the story reads no other open story's output.

Figures:

- Forecast 40m, delivery 2026-11-13T23:33:00Z: `flai forecast` gave 23m (136 s per unit of size over 64 done medium remediation stories on claude-opus-5-5, size 10) and delivery 23:16Z, 30th in the pull order. Raised to 40m because the story writes an ADR and edits three long design and docs passages beside the code, and its reproduction drives a real rebase stop and a resumed sync; delivery moved by the same 17m.
- Cost of delay 25 USD a week, as `flai cod` works it from the operator's input: 10m lost per 168h cycle at 150 USD an hour. It stands: one occurrence so far, and each costs a hand `--keep-removal` or a send-back at acceptance.

### Accepted by the orchestrator

- Verified: f2ce9bc01be52056d53801f32a79c11adff2e330
- At: 2026-11-13T12:53:44Z

Verdict: pass — both criteria are met at f2ce9bc, and the diff has no breach that should stop acceptance.
- 1: flai/internal/storygit/undone.go, flai/internal/storygit/undone_test.go, flai/internal/storygit/sync_test.go, design/adrs/0297-a-markdown-line-of-main-s-that-a-rebased-branch-removes-is-not-undone-when-a.md, design/system/flai-cli.md, design/system/workflow.md, docs/users/flai.md
- 2: design/issues/I-0360-flai-stream-sync-reports-a-rebased-branch-undoes-main-s-line-when-both-changed.md, design/issues/summary.md
