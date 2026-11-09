---
id: S-0598
type: story
nature: improvement
title: A checked write runs only the rules an item write can affect, so creating or editing an item takes about 2 s, not 7 to 18 s
status: review
parent: E-0033
owner: alex
created: 2026-11-07T22:39:02Z
updated: 2026-11-09T07:41:24Z
transitions:
  - to: ready
    at: 2026-11-07T23:49:08Z
    by: alex
  - to: in-progress
    at: 2026-11-07T23:49:17Z
    by: agent-S-0598
  - to: review
    at: 2026-11-08T01:11:52Z
    by: system-flow
  - to: in-progress
    at: 2026-11-09T07:24:16Z
    by: flai
  - to: review
    at: 2026-11-09T07:41:24Z
    by: system-flow
tags: [cli]
touches: [flai/internal/itemnew/itemnew.go, flai/internal/itemedit/itemedit.go, flai/internal/check/check.go, flai/internal/check/itemwrite.go, flai/internal/check/itemwrite_test.go, flai/internal/itemnew/itemnew_test.go, flai/internal/itemedit/rules_test.go, design/system/server-performance.md, design/system/flai-cli.md, design/system/strategic-agents.md, docs/users/flai.md, design/adrs/0264-a-checked-write-of-a-work-item-runs-only-the-rules-the-write-can-affect-before.md, design/adrs/README.md, flai/internal/itemedit/criteria_test.go]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 6005
  priced_by: harness
  turns:
    - day: 2026-11-07
      hand_edits: 1
      work: 21
    - day: 2026-11-08
      ceremony: 2
      hand_edits: 2
      work: 75
    - day: 2026-11-09
      ceremony: 3
      test_runs: 1
      hand_edits: 1
      work: 66
  models:
    - model: claude-opus-5-5
      input: 518
      output: 143375
      cache_read: 98733594
      cache_write: 1789930
      cost: 34.4828
  strategic:
    - kind: orchestrator
      seconds: 107
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 22
          output: 146
          cache_read: 699967
          cache_write: 123177
          cost: 0.2113
cost_of_delay:
  inputs:
    penalty_per_week: 150
    by: alex
    at: 2026-11-07T23:52:46Z
forecast:
  duration: 45m
  delivery: 2026-11-08T00:35:00Z
  basis: "flai forecast's 26m (median 94 s per unit over 124 done improvement stories in the large band, times size 16), raised to 45m for a refusal test per kept rule on both writes and the scratch-copy timing, which size does not count; started 23:49Z."
  by: planner-S-0598
  at: 2026-11-07T23:50:00Z
finalized:
  by: alex
  at: 2026-11-07T23:48:54Z
---
# S-0598 A checked write runs only the rules an item write can affect, so creating or editing an item takes about 2 s, not 7 to 18 s

## Goal

Cause 8 in [server-performance.md](../../../design/system/server-performance.md#where-the-time-goes): a checked creation or edit runs the whole of `flai check` twice, before and after. Each run lists all 3,233 items cold, lints all 5,068 markdown files under `wip`, starts `node` for the component examples, and scans the digest for secrets. `story new`, `epic new`, `task new` (and MCP `item_new`), `edit`, finalize, `criteria tick`, and `issue story` take 7 to 18 s wall, 8 to 12.7 s CPU, about 97% of it in the two checks. The serve log shows `item.new` at p50 7.9 s and up to 32 s.

## Acceptance criteria

- [x] A checked creation or edit runs only the rules a work item write can affect: the written items and their parent, the after, overlap, and touches graph, the board, the story's narrative, and the markdown lint and conflict markers of the written files; and it computes the "before" result for those rules alone.
- [x] It still refuses exactly the findings ADR-0197 and ADR-0213 refuse, which tests show for a creation and an edit that bring a finding of each kept rule.
- [x] `flai story new`, `flai edit`, and MCP `item_new` of a task, timed in a scratch copy of this repository, each take under 3 s, and server-performance.md records the before and after.

## Tasks
- T-2620 check runs an item write's rule set: the written items and their parent, the graph rules, the board, the story's narrative, and the written files' lint and conflict markers
- T-2622 itemnew.Create checks a creation with the item-write rule set before and after, in place of check.Run
- T-2625 itemedit.Apply checks an edit with the item-write rule set before and after, in place of check.Run
- T-2635 Time story new, edit, and MCP item_new of a task in a scratch copy, before and after, and record them in server-performance.md
- T-2640 Write the ADR that a checked write runs only the rules an item write can affect, and say so in flai-cli.md, strategic-agents.md, and the user guide

## Notes

- Proposed fix: a rule set for item writes in `check`, used by `itemnew.Create` and `itemedit.Apply`, in place of `check.Run`. The whole check stays in `flai check`, the close-out, and CI.
- If the scope of what a checked write refuses changes, that refines ADR-0197 and needs an ADR.

- Criterion 3 was measured twice in a scratch clone. `design/system/server-performance.md` § A checked write runs only the rules it can affect (S-0598) has the figures.
  - On 2026-11-08, with every item parsed cold, the median wall time was 2.4 to 2.7 s. Some runs went over 3 s at load 46 to 131.
  - On 2026-11-09, rebased onto S-0600's kept parses, every run of `story new`, `edit`, and MCP `item_new` of a task took 0.7 to 1.3 s wall. Current `main` took 6.1 to 12.1 s.

### Planning

Touches, and where each came from:

| Touch | Source |
|-------|--------|
| `flai/internal/itemnew/itemnew.go`, `flai/internal/itemedit/itemedit.go`, `flai/internal/check/check.go` | declared |
| `flai/internal/check/itemwrite.go`, `flai/internal/check/itemwrite_test.go` | layout: the new rule set and its parity test (T-2620) |
| `flai/internal/itemnew/itemnew_test.go` | co-change and layout: the creation's refusal tests (T-2622) |
| `flai/internal/itemedit/rules_test.go` | layout: the edit's refusal tests (T-2625) |
| `design/system/server-performance.md` | design: the goal links it; criterion 3 records there (T-2635) |
| `design/system/flai-cli.md`, `docs/users/flai.md` | co-change (18% and 19%) and design: both say a write runs the whole check (T-2640) |
| `design/system/strategic-agents.md` | design: it says a planner's `item_new` runs `flai check` (T-2640) |
| `design/adrs/0264-a-checked-write-of-a-work-item-runs-only-the-rules-the-write-can-affect.md`, `design/adrs/README.md` | design: the ADR that refines ADR-0197 and ADR-0213 (T-2640); the number is the next free one today, and T-2640 narrows it to the file `flai adr new` names |

No folder touches. `flai touches suggest` also listed `flai/internal/check/check_test.go` (27%), `flai/cmd/edit.go`, `flai/cmd/items.go`, and the manifest files. None is predicted: the commands call `itemnew.Create` and `itemedit.Apply` and need no change, and no manifest key is planned. `flai/internal/serve/replan.go`, `flai/cmd/criteria.go`, and `flai/internal/mcpserver/` save through `Apply` or `Create` and get the change with no edit.

No `after`: the story reads no other story's output. S-0601 (check's rules at once) and S-0599 (the replanner) touch `check.go` and the replan path too, and the trial merge judges that.

Forecast: `flai forecast` gave 26m (median 94 s per unit of size over 124 done improvement stories in the large band, times size 16). It is raised to 45m: criterion 2 asks for a refusal test of each kept rule, for both a creation and an edit, and criterion 3 asks for timings in a scratch copy, before and after. Size counts neither. The delivery is 45m after the story started at 23:49Z.

Cost of delay: `flai cod` refuses for now, since neither the story nor E-0033 has inputs. TH-0888 asks the operator for them and recommends `penalty_per_week: 150` for S-0598. The value is written once the operator answers.
- 2026-11-09T07:24:16Z: moved to in-progress: the rebase of story/S-0598 onto main stopped on conflicts in flai/internal/check/check.go, which the acceptance of S-0598 by alex cannot settle, so it aborted the rebase; its agent settles it in the worktree and syncs it with flai stream sync S-0598
