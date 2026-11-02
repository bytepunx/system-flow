---
id: S-0563
type: story
nature: remediation
title: The orchestrator's conventions sit within about 1.5KB of its 80KB pack budget, so a few more lines in work-management.md or strategic-agents.md drop every design brief from its pack
status: done
owner: alex
created: 2026-11-02T04:21:24Z
updated: 2026-11-02T05:02:10Z
transitions:
  - to: ready
    at: 2026-11-02T04:22:12Z
    by: alex
  - to: in-progress
    at: 2026-11-02T04:22:51Z
    by: system-flow
  - to: review
    at: 2026-11-02T04:44:21Z
    by: agent-S-0563
  - to: done
    at: 2026-11-02T05:02:10Z
    by: alex
blocked:
  - from: 2026-11-02T04:26:06Z
    until: 2026-11-02T04:27:17Z
    reason: awaiting the operator on TH-0805
tags: []
topics: [context, prime]
touches: [design/adrs, design/adrs/README.md, design/system/agent-context.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md, flai/cmd/prime.go, flai/cmd/prime_test.go, flai/internal/context/context.go, flai/internal/context/brief_test.go, design/issues/I-0250-the-orchestrator-s-conventions-sit-within-about-1-5kb-of-its-80kb-pack-budget.md, design/issues/summary.md, design/adrs/0241-a-pack-whose-conventions-alone-exceed-its-budget-still-loads-what-is-named-and.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1224
  priced_by: harness
  turns:
    - day: 2026-11-02
      ceremony: 4
      test_runs: 2
      hand_edits: 2
      work: 55
  models:
    - model: claude-opus-5-5
      input: 136
      output: 34429
      cache_read: 10775856
      cache_write: 224914
      cost: 4.6436
  strategic:
    - kind: planner
      seconds: 482
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 78
          output: 37907
          cache_read: 6267524
          cache_write: 221352
          cost: 3.7828
    - kind: orchestrator
      seconds: 200
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 72
          output: 479
          cache_read: 2269727
          cache_write: 374382
          cost: 0.6158
        - model: claude-sonnet-5-5
          input: 8
          output: 58
          cache_read: 77804
          cache_write: 36302
          cost: 0.0797
cost_of_delay:
  inputs:
    time_lost_per_cycle: 36m
    by: flai
    at: 2026-11-02T04:21:24Z
  value: 90
  by: planner-S-0563
  at: 2026-11-02T04:29:41Z
forecast:
  duration: 3h
  delivery: 2026-11-02T07:23:00Z
  basis: "Adjusted from flai forecast's 4m (size 2, no touches) to 3h: a context-pack change with its tests, an ADR, design and user docs, and the issue's close, in three tasks"
  by: planner-S-0563
  at: 2026-11-02T04:29:41Z
finalized:
  by: alex
  at: 2026-11-02T04:21:53Z
---
# S-0563 The orchestrator's conventions sit within about 1.5KB of its 80KB pack budget, so a few more lines in work-management.md or strategic-agents.md drop every design brief from its pack

## Goal

This story remediates [I-0250](../../../design/issues/I-0250-the-orchestrator-s-conventions-sit-within-about-1-5kb-of-its-80kb-pack-budget.md), "The orchestrator's conventions sit within about 1.5KB of its 80KB pack budget, so a few more lines in work-management.md or strategic-agents.md drop every design brief from its pack". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0250 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0250 is closed with `flai issue close I-0250 --reason` saying what fixed it

## Tasks
- T-2429 A pack whose conventions exceed its budget still loads what is named and every brief, ranking nothing
- T-2430 Record the decision in an ADR refining ADR-0049 and update agent-context.md, flai-cli.md, and the prime docs
- T-2431 I-0250 is closed with the reason that a pack over its budget on its conventions now keeps its briefs

## Notes

Cost of delay inputs set by flai from I-0250. time_lost_per_cycle 36m: 18m per occurrence × 2 occurrences ÷ 1 cycle of 168h (first reported 2026-11-02T03:46:07Z, 0 days before this story; under one cycle counts as one).

Made by agent-S-0518 on the operator's word on TH-0804 (option 1, confirmed 2026-11-02T04:21:05Z). S-0518 waits in progress until this fix is on main, since its close-out and the merge queue's batch both stop at `TestAThreadPackIsUnderAQuarterOfTheGeneralPack`.

The orchestrator measured the test on TH-0804. In the main checkout the thread pack is 22.6% of the general pack, and the test passes. In S-0518's worktree it is 28.3%, and the test fails. The thread pack is 43,932 bytes in both. The general pack drops its design briefs in the worktree (155,062 against 194,234 bytes), and the only difference there is S-0518's issue records under `design/issues`. So the result depends on what the repository holds, wip and issues included. The fix should make the test independent of that, or give the orchestrator's budget a real margin, as well as trimming its conventions.

### Planning

The planner measured main at 2026-11-02T04:24Z with `flai prime --role orchestrate --json`. The budget is 81,920 bytes. The conventions take 72,769 of it:

| Convention | Bytes |
|------------|-------|
| work-management.md | 27,333 |
| strategic-agents.md | 26,592 |
| tooling.md | 7,010 |
| documentation.md | 3,578 |
| decisions.md | 3,224 |
| safety.md | 2,870 |
| communication.md | 2,162 |

The open-issue table, 22 open issues, takes about 7.6 KB, and the digest step 7,948. `addDesign` in `flai/internal/context/context.go` sizes the pack before its first check. By then `Build` has already put the issue table in. So the check that sets `exceeded: conventions` counts the conventions and the issue table together, and that check loads no brief. Each open issue a branch adds uses up the room left, which is why S-0518's issue records tipped it.

The story's agent asked on TH-0805, and the operator chose option 1 at 2026-11-02T04:27:14Z: remove the cliff. When the conventions take the pack over, it still loads what is named and every brief, ranks nothing, and keeps `exceeded: conventions`. This applies TH-0032's "briefs are never cut" (D-0048) to that case too. Options 2 and 3 are left out. Option 1 needs no trim of `work-management.md` or `strategic-agents.md`. S-0553's branch also conflicts with main in `design/conventions/strategic-agents.md`, so touching it would hold this story.

The story's agent wrote T-2429 and T-2430; the planner added T-2431 for criterion 2. The plan's layers: T-2429, then T-2430, then T-2431.

Touches, and where each came from:

| Touch | From | Why |
|-------|------|-----|
| `flai/internal/context/context.go` | declared (T-2429), design | `addDesign`'s first check is the cliff; its header line for `exceeded: conventions` says what the pack holds |
| `flai/internal/context/brief_test.go` | declared (T-2429) | `TestTheHeaderSaysWhatAloneExceedsTheBudget` asserts today's no-briefs result |
| `flai/cmd/prime_test.go` | co-change | asserts the header line `over budget: the conventions alone exceed it`, which T-2429 rewrites |
| `flai/cmd/prime.go` | co-change | its help says the pack is the conventions and a catalog when they alone exceed the budget |
| `docs/users/flai-reference.md` | co-change | generated from `prime.go`'s help by `make flai-reference` |
| `docs/users/flai.md` | co-change | the `flai prime` section says the same |
| `design/system/agent-context.md`, `design/system/flai-cli.md` | declared (T-2430) | they describe the conventions case |
| `design/adrs` | design | the ADR T-2430 writes, whose file name its title sets, so no task can name it yet: the one folder touch kept |
| `design/adrs/README.md` | design | the ADR index flai regenerates |
| `design/issues/I-0250-…md`, `design/issues/summary.md` | criterion 2 | `flai issue close I-0250` writes both (T-2431) |

`flai touches suggest` also listed `flai/internal/context/story.go`, `strategic.go`, `design.go`, and `flai/internal/mcpserver/folder.go`. They are left out. Option 1 does not change how the pack is built, and the prime tool's description already says briefs are never cut.

Forecast: `flai forecast` gave 4m from size 2 (2 criteria, 0 touches) before any touches were set. Adjusted to 3h, counted from its start at 2026-11-02T04:23:00Z. That covers a context-pack change with its tests, an ADR, design and user docs, and the issue's close, in three tasks.

Cost of delay: `flai cod` gives 90.00 USD a week from the inputs, 36m lost per 168h cycle at 150 USD an hour. It stands as the value. The inputs are the operator's, and they understate the delay: until this is on main, every story whose change selects the integration tier stops at this test, S-0518 among them. The plan's thread says so.

No `after`: this story reads no other open story's output.
