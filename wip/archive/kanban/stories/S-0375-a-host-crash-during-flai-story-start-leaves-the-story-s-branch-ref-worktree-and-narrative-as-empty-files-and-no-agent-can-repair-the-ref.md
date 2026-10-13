---
id: S-0375
type: story
nature: remediation
title: A host crash during flai story start leaves the story's branch ref, worktree, and narrative as empty files, and no agent can repair the ref
status: done
owner: alex
created: 2026-10-12T16:52:56Z
updated: 2026-10-13T17:16:35Z
transitions:
  - to: ready
    at: 2026-10-13T16:38:06Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-13T16:38:12Z
    by: agent-S-0375
  - to: review
    at: 2026-10-13T17:09:21Z
    by: agent-S-0375
  - to: done
    at: 2026-10-13T17:16:35Z
    by: orchestrator
tags: [flai, git]
topics: [git, narrative]
touches: [flai/internal/atomicfile/atomicfile.go, flai/internal/atomicfile/atomicfile_test.go, flai/internal/workitem/narrative.go, flai/internal/workitem/narrative_test.go, flai/internal/storygit/storygit.go, flai/internal/storygit/open.go, flai/internal/storygit/open_test.go, flai/cmd/stream.go, flai/internal/storystart/start_test.go, design/system/flai-cli.md, design/system/agent-narrative.md, docs/users/flai.md, docs/operators/runbooks/restore.md, design/issues/I-0137-a-host-crash-during-flai-story-start-leaves-the-story-s-branch-ref-worktree-and-narrative-as-empty-files-and-no-agent-can-repair-the-ref.md, design/issues/summary.md, flai/cmd/stream_reopen_test.go]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1885
  turns:
    - day: 2026-10-13
      ceremony: 2
      hand_edits: 2
      work: 41
  models:
    - model: claude-opus-5-5
      input: 296
      output: 103436
      cache_read: 13611026
      cache_write: 560675
      cost: 8.3147
  strategic:
    - kind: orchestrator
      seconds: 1104
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 125
          output: 1991
          cache_read: 10070597
          cache_write: 124370
          cost: 2.5705
        - model: claude-sonnet-5-5
          input: 12
          output: 71
          cache_read: 177287
          cache_write: 61307
          cost: 0.1928
cost_of_delay:
  inputs:
    time_lost_per_cycle: 30m
    by: flai
    at: 2026-10-12T16:52:56Z
  value: 75
  by: planner-S-0375
  at: 2026-10-13T16:37:32Z
forecast:
  duration: 1h
  delivery: 2026-10-13T18:27:00Z
  basis: "flai forecast's 24m (84 s per unit over 26 large remediation stories, size 17) raised to 1h for three code changes in layer 1, one of which removes files under .git and needs crash-reproducing tests with real git, plus a test task, docs in four files, and the issue close; delivery moved by the same 36m from flai's 17:51."
  by: planner-S-0375
  at: 2026-10-13T16:37:32Z
finalized:
  by: orchestrator
  at: 2026-10-13T16:38:02Z
---
# S-0375 A host crash during flai story start leaves the story's branch ref, worktree, and narrative as empty files, and no agent can repair the ref

## Goal

This story remediates [I-0137](../../../design/issues/I-0137-a-host-crash-during-flai-story-start-leaves-the-story-s-branch-ref-worktree-and-narrative-as-empty-files-and-no-agent-can-repair-the-ref.md), "A host crash during flai story start leaves the story's branch ref, worktree, and narrative as empty files, and no agent can repair the ref". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0137 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0137 is closed with `flai issue close I-0137 --reason` saying what fixed it

## Tasks
- T-1500 atomicfile.WriteFile syncs the file before the rename and the folder after it
- T-1501 Opening a stream over an empty narrative writes it afresh from the template
- T-1502 flai stream open repairs a broken story ref that holds no commits and a worktree git no longer knows
- T-1503 flai story start succeeds over the empty narrative, ref, and worktree a crashed start left
- T-1504 The design, the user guide, and the restore runbook say how a start recovers from a crash
- T-1505 Close I-0137 with what fixed it

## Notes

Cost of delay inputs set by flai from I-0137. time_lost_per_cycle 30m: 30m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-12T15:29:54Z, 0.1 days before this story; under one cycle counts as one).

### Planning

Proposed remedy, from I-0137's instance and TH-0390: make flai's own writes durable, and let `flai stream open` repair what a crashed start left, but only what holds no commits.

1. `atomicfile.WriteFile` renames without `Sync`, so a crash can keep the rename and lose the data: the empty narrative. It syncs the file and its folder (T-1500).
2. An empty narrative holds nothing, so `OpenStream` writes it afresh; one with content is never replaced (T-1501).
3. `BranchExists` reads a broken ref as missing, and git refuses to recreate it. `storygit.Open` removes the story's own ref and reflog when the ref resolves to no commit, and a worktree folder git does not know when the branch holds no commit, and refuses otherwise. It creates the branch with `core.fsync=reference` (T-1502).

Touches, all files, no folder touch:

| Touch | Source |
|-------|--------|
| `flai/internal/atomicfile/atomicfile.go`, `atomicfile_test.go` | layout: the explorer found the helper without `Sync` |
| `flai/internal/workitem/narrative.go` | layout: `OpenStream` and `ReopenStream` |
| `flai/internal/workitem/narrative_test.go` | co-change with `narrative.go` (4 of 12 commits) |
| `flai/internal/storygit/storygit.go`, `open.go`, `open_test.go` | layout: `BranchExists` and `Open`, named by I-0137's `git worktree add` |
| `flai/cmd/stream.go` | co-change (4 of 12) and layout: `canReopen` |
| `flai/internal/storystart/start_test.go` | layout: the end-to-end test for criterion 1 |
| `design/system/flai-cli.md` | co-change (5 of 12) and design: the atomic-write and story branch text |
| `design/system/agent-narrative.md` | co-change (4 of 12) and design: § Recovery procedure |
| `docs/users/flai.md` | co-change (3 of 12) and the documentation convention |
| `docs/operators/runbooks/restore.md` | layout: the runbook that puts back a story worktree |
| the I-0137 file, `design/issues/summary.md` | criterion 2; `summary.md` also co-change (2 of 12) |

The story declared no touches; `flai touches suggest` was seeded with `open.go`, `atomicfile.go`, and `narrative.go`.

Figures:

- Forecast: flai's 24m raised to 1h. Three code changes, one of which removes files under `.git` and needs tests with real git, plus a test task, four documents, and the issue close. Delivery moved by the same 36m from flai's 17:51.
- Cost of delay: 75 USD a week, `flai cod`'s figure from flai's 30m input, kept: one occurrence, but it held a story in progress for four days and stopped another's acceptance.

### Accepted by the orchestrator

- Verified: 0797010753f9957423bd95c3668bd169dc94fbc1
- At: 2026-10-13T17:16:35Z

Verdict: meets all criteria (verifier at 0797010753f9957423bd95c3668bd169dc94fbc1; flai verify passed every step at that commit)

- 1: flai/internal/atomicfile/atomicfile.go, flai/internal/atomicfile/atomicfile_test.go, flai/internal/workitem/narrative.go, flai/internal/workitem/narrative_test.go, flai/internal/storygit/storygit.go, flai/internal/storygit/open.go, flai/internal/storygit/open_test.go, flai/cmd/stream.go, flai/cmd/stream_reopen_test.go, flai/internal/storystart/start_test.go, design/system/flai-cli.md, design/system/agent-narrative.md, docs/users/flai.md, docs/operators/runbooks/restore.md
- 2: design/issues/I-0137-a-host-crash-during-flai-story-start-leaves-the-story-s-branch-ref-worktree-and-narrative-as-empty-files-and-no-agent-can-repair-the-ref.md, design/issues/summary.md
