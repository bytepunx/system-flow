---
id: TH-0952
title: "S-0631 plan: split six files by responsibility, then a size test; nine tasks in three layers"
anchor:
  path: wip/kanban/stories/S-0631-five-source-files-under-flai-are-over-2-000-lines.md
  item: S-0631
status: resolved
participants: [planner-S-0631, orchestrator, agent-S-0631, alex]
created: 2026-11-09T22:06:12Z
updated: 2026-11-10T07:00:36Z
---

# TH-0952 S-0631 plan: split six files by responsibility, then a size test; nine tasks in three layers

On wip/kanban/stories/S-0631-five-source-files-under-flai-are-over-2-000-lines.md.

## Entries

### 2026-11-09T22:06:12Z planner-S-0631
**Proposed fix:** split each non-test Go file under `flai/` that is over 2,000 lines into files of its own package, by responsibility. Only declarations move. Then a test fails on any such file over 2,000 lines.

## Tasks and layers

| Layer | Task | What |
|-------|------|------|
| 1 | T-2794 | Split `flai/internal/hostapi/writes.go` (3,626) into 9 files, and cut the 1,476-line `itemSpecs` into three functions |
| 1 | T-2795 | Split `flai/internal/serve/agents.go` (3,298) into 6 files |
| 1 | T-2796 | Split `flai/cmd/serve_actions.go` (2,754) into one file per `flai serve` subcommand family |
| 1 | T-2797 | Split `flai/internal/manifest/manifest.go` (2,449) into one file per manifest block |
| 1 | T-2798 | Split `flai/internal/check/check.go` (2,235) into its runner and 5 rule families |
| 1 | T-2799 | Split `flai/internal/guard/guard.go` (2,051) into 5 files |
| 1 | T-2800 | Give the `//nolint:nilerr` in `flai/internal/watch/watch.go`, now at line 94, its reason |
| 2 | T-2801 | `flai/internal/check/sourcesize_test.go`: fails on any non-test Go file under `flai/` over 2,000 lines. Not skipped under `-short` |
| 3 | T-2802 | Point the 9 design documents that name these files at where the code moved, and close I-0298 |

Layer 1 runs together: its seven tasks share no path. T-2801 waits for the six splits, because it fails until each is done. T-2802 waits for T-2801, because its closing reason names the test and needs the final file names.

## Figures

- Forecast 3h, against flai's 1h38m. flai's size counts touches; about 16,400 lines move across six packages.
- Cost of delay 75 USD a week, from the 30m `time_lost_per_cycle` set on TH-0944.

## Assumptions

1. **`guard.go` is in scope.** I-0298 does not name it, but it has grown to 2,051 lines, and the test would fail on it.
2. **The test covers non-test files only.** Test files stay as they are.
3. **Every split keeps its original file name.** `flai/internal/verify/select_test.go:554` stats `writes.go`, and the design documents' references to code that stays remain true.
4. **The new file names are proposals.** Each task may choose better ones, keeping every file under 1,200 lines. Its touches widen to what it uses when it closes.
5. **A test rather than a lint rule.** revive's `file-length-limit` (golangci-lint v2.5.0, `flai/.golangci.yaml`) would also work. A Go test meets criterion 1 directly and runs in every close-out's integration tier.
6. **Accepted ADRs, the digest, the analysis, and other issues** that name these files stay as they are.

## Proposed: one story to add, not drafted

Seven `_test.go` files are over 2,000 lines: `serve/agents_test.go` 4,919, `hostapi/writes_test.go` 3,884, `harness/harness_test.go` 3,453, `cmd/serve_actions_test.go` 3,078, `manifest/manifest_test.go` 2,503, `cmd/dashboard_test.go` 2,186, and `mcpserver/server_test.go` 2,185.

I recommend a follow-up story that splits them along the same lines as their sources and widens the test to cover them. Keeping it out of S-0631 keeps this story to moves the unchanged tests prove safe.

## Risk

These files are hubs: 95 stories changed at least one of them in the last 30 days. Other stories in progress that edit them will conflict with the moves on sync, and will need to re-apply their hunks to the new files. The story is best pulled when few `flai/` stories are in progress.

### 2026-11-09T22:18:10Z orchestrator (recommendation)
**Recommended:** keep S-0631's plan as written, and leave the seven test files over 2,000 lines to a follow-up story drafted once S-0631 is accepted, not now.

**The plan fits S-0631.**

- Layer 1's seven tasks (T-2794 to T-2800) share no path and run together. T-2801's size test reproduces I-0298 for criterion 1 and waits for the six splits it would fail on. T-2802 points the design at the moved code and closes I-0298 for criterion 2.
- Taking in `guard.go` is within criterion 1: it is over 2,000 lines now, and the test would fail on it.
- Every split keeps its file's name, so `select_test.go:554` and the design's references to code that stays hold.
- Touches name files only. 3h over flai's 1h38m gives its reason: about 16,400 lines move, and `itemSpecs` is cut in three.

**The follow-up story: later.** Splitting the test files is new scope, so it is your call. I recommend it, along the lines S-0631's splits draw, so draft it after S-0631 lands rather than now: its new file names come from S-0631's.

**The risk the planner names stands.** S-0623 changes `guard.go`, and S-0632 and S-0633 change `serve/agents.go`, all in progress now. A trial merge will hold S-0631 on any conflict (ADR-0230), and the pull order already has it last of four in ready, so those stories will likely land first.

Source: design/conventions/work-management.md

Source: design/conventions/work-management.md

### 2026-11-10T06:49:55Z agent-S-0631
S-0631 was worked to this plan as written. All nine tasks are done, and both criteria are ticked.

- **The six splits:** each moved declarations without changing them. A line-by-line comparison outside the import blocks confirms it. The one exception is `itemSpecs`, cut into `itemWriteSpecs`, `streamSpecs`, and `hostSpecs`. The largest new file is 901 lines.
- **The size test:** `TestNoSourceFileIsOver2000Lines` fails on main, naming the six files. It passes on the branch.
- **I-0298:** closed on the branch.

The seven test files over 2,000 lines are untouched, as recommended. Whether to draft that follow-up story is still yours.

### 2026-11-10T06:51:03Z orchestrator (recommendation)
**Recommended:** draft one follow-up story, once S-0631 is accepted, that splits the seven `_test.go` files over 2,000 lines along the lines S-0631's splits drew, and widens `TestNoSourceFileIsOver2000Lines` to cover test files.

**The plan is settled.** S-0631 was worked to it as written: nine tasks done, both criteria ticked, the largest new file 901 lines. Nothing in it needs an answer any more.

**The follow-up is yours.** It adds scope, and neither the orchestrator nor S-0631's agent authors stories. Its file names now exist on `story/S-0631`, so it can be drafted as soon as S-0631 lands:

| File | Lines |
|------|-------|
| `flai/internal/serve/agents_test.go` | 4,919 |
| `flai/internal/hostapi/writes_test.go` | 3,884 |
| `flai/internal/harness/harness_test.go` | 3,453 |
| `flai/cmd/serve_actions_test.go` | 3,078 |
| `flai/internal/manifest/manifest_test.go` | 2,503 |
| `flai/cmd/dashboard_test.go` | 2,186 |
| `flai/internal/mcpserver/server_test.go` | 2,185 |

The thread stays open for that decision only. S-0631's acceptance does not wait on it.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-10T07:00:36Z alex
Resolved: S-0631 was accepted
