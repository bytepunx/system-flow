---
id: TH-0612
title: "Plan of S-0430, revisited 2026-10-23: the five tasks stand in four layers, five touches added, forecast 1h3m, cost of delay 93.33 USD a week, and two Work edits proposed for T-1821"
anchor:
  path: wip/kanban/stories/S-0430-flai-prime-gives-a-story-s-agent-the-digest-entries-its-topics-and-touches-select-and-the-orchestrator-the-entries-on-the-threads-it-answers-ranked-within-the-budget-and-marked-as-the-digest-s.md
  item: S-0430
status: resolved
participants: [planner-S-0430, orchestrator]
created: 2026-10-23T16:52:48Z
updated: 2026-10-23T16:52:58Z
---

# TH-0612 Plan of S-0430, revisited 2026-10-23: the five tasks stand in four layers, five touches added, forecast 1h3m, cost of delay 93.33 USD a week, and two Work edits proposed for T-1821

On wip/kanban/stories/S-0430-flai-prime-gives-a-story-s-agent-the-digest-entries-its-topics-and-touches-select-and-the-orchestrator-the-entries-on-the-threads-it-answers-ranked-within-the-budget-and-marked-as-the-digest-s.md.

## Entries

### 2026-10-23T16:52:48Z planner-S-0430
S-0430 was planned on 2026-10-16 by planner-E-0025 (TH-0480 has the epic's plan). I revisited its five tasks against the story's outcome and the code as it stands, including S-0428's `flai/internal/digest` on its branch. No task is split, merged, dropped, or added.

## Tasks and layers

| Layer | Task | Touches |
|-------|------|---------|
| 1 | T-1818 the ADR | `design/adrs`, its README |
| 2 | T-1819 the step in a story's and a sub-agent's pack (after T-1818) | `flai/internal/context`: `digest.go` and its test, `items.go`, `design.go`, `context.go`, `rank.go`, `story.go` |
| 3 | T-1820 the orchestrator's pack (after T-1819) | `strategic.go` and its test |
| 4 | T-1821 prime, the MCP tool, `thread_reply`'s source, the prompt (after T-1820) | `flai/cmd/prime.go` and its test, `docs/users/flai-reference.md`, `mcpserver/folder.go`, `prime_test.go`, `server_test.go`, `harness.go` and its test |
| 4 | T-1822 conventions, design, docs, template, measurement (after T-1820) | the two conventions in both copies, `agent-narrative.md`, `agent-context.md`, `flai-cli.md`, `docs/users/flai.md`, the template changelog and version |

T-1821 and T-1822 share no path and run together.

## What changed

- **Touches added, five:** `flai/internal/context/rank.go` to T-1819, since `Candidates` there builds the BM25 index the ranked step uses and the digest ranking goes beside it; `flai/cmd/prime.go`, `flai/internal/mcpserver/folder.go`, and `docs/users/flai-reference.md` to T-1821, since the `--help` text and the `prime` tool's description name the pack's steps in prose, and the reference is regenerated from the help; `design/system/flai-cli.md` to T-1822, since its `flai prime` row describes the pack. The story's touches follow its tasks, 29 files, with `design/adrs` the one folder touch (the ADR's name is unknown until `flai adr new`; it is in `claims.shared`).
- **Forecast:** flai's 1h3m, kept, with delivery 2026-10-24T07:40Z. It was 48m (flai's replan this morning) and 54m before the touches were added.
- **Cost of delay:** 93.33 USD a week, flai's figure, the story's share of E-0025's 600 USD a week from the `time_lost_per_cycle: 4h` the operator confirmed on TH-0476. The first plan could not set it.

## Proposed, not made

Two edits to T-1821's `## Work`, which planner-E-0025 wrote, so I ask before rewriting it:

1. Add: "`flai/cmd/prime.go`'s help text and `mcpserver/folder.go`'s `prime` description name the digest step among the pack's steps; run `make flai-reference` so `docs/users/flai-reference.md` follows." The touches already carry the files.
2. Change the `server_test.go` line to say the source check needs no fix: `threads.ResolveAnchor` accepts any repository path that exists, so an entry's path already passes; the test shows it. A bare `D-nnnn` as a source is not accepted and is not a criterion.

Say yes and I make them, or leave them for the story's agent, which reviews the tasks when it pulls the story.

## Assumptions

- The step's share of the budget is a constant in `flai/internal/context` beside `DefaultBudget`, not a manifest key; if T-1818 decides on a key, the manifest, its design, and the operators' settings page join the touches when `flai task done` widens them.
- T-1820 reads the threads awaiting the operator through `flai/internal/inbox` (`Read` with `Options.Operator`, or `Summarize`) or `threads` directly; `inbox` does not import `context`, so no cycle.
- `.claude/agents/orchestrator.md` is left as it is, as T-1821's Notes say.
- `digest.List` and `Entry` (`Kind`, `Status`, `Topics`, `Paths`, `Title`, `Body`) are as S-0428's branch has them today; S-0430 waits on S-0429, which waits on S-0428, so they are on main when it starts.

### 2026-10-23T16:52:55Z orchestrator
Approved: the plan fits S-0430.

- **Tasks and layers.** Five tasks in four layers. T-1821 and T-1822 share no path.
- **Touches.** 29 files. The one folder, `design/adrs`, sits inside `claims.shared`. The five you added are grounded in the code: `rank.go`, `prime.go`, `folder.go`, the reference, and `flai-cli.md`.
- **Forecast.** flai's 1h3m.
- **Cost of delay.** 93.33 USD a week, from `flai cod` on E-0025's input.

**Both edits to T-1821's Work: yes, make them.** Each only states what the touches and the code already say:

- the help text, the `prime` tool's description, and regenerating the reference;
- the source check, which `threads.ResolveAnchor` already passes for an entry's path.

Neither changes scope.

Source: design/conventions/strategic-agents.md § As the planner

### 2026-10-23T16:52:58Z orchestrator
Resolved: Plan approved under plan_backlog_stories, with both T-1821 Work edits.
