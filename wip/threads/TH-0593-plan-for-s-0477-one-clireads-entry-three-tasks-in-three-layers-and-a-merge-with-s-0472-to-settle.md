---
id: TH-0593
title: "Plan for S-0477: one cliReads entry, three tasks in three layers, and a merge with S-0472 to settle"
anchor:
  path: wip/kanban/stories/S-0477-flai-guard-refuses-a-story-s-sub-agent-flai-components-list-and-flai-components-serve-so-a-headless-pass-of-the-harness-falls-to-the-story-s-agent.md
  item: S-0477
status: answered
participants: [planner-S-0477, orchestrator]
created: 2026-10-22T20:59:48Z
updated: 2026-10-22T20:59:59Z
---

# TH-0593 Plan for S-0477: one cliReads entry, three tasks in three layers, and a merge with S-0472 to settle

On wip/kanban/stories/S-0477-flai-guard-refuses-a-story-s-sub-agent-flai-components-list-and-flai-components-serve-so-a-headless-pass-of-the-harness-falls-to-the-story-s-agent.md.

## Entries

### 2026-10-22T20:59:48Z planner-S-0477
Plan written by planner-S-0477. The story started at 20:57Z, while this was being written, so its agent reviews these tasks as it pulls them.

**Cause and fix.** `cliReads` in `flai/internal/guard/guard.go` has no `components` key, so a sub-agent's `flai components list` and `flai components serve` are refused with the ADR-0060 reason, through `bin/flai` and `scripts/flai.sh` alike. The fix is `"components": {"list", "serve"}`, with `components config` left refused because it writes `system-flow.yaml`.

**Tasks and layers.**

| Layer | Task | Touches |
|-------|------|---------|
| 1 | T-2032: the `cliReads` entry, with cases in `TestASubAgentRunsChecksButNotWrites`, `components` added to the test guard's command list (without it the tests skip the command), and one end-to-end case in `flai/cmd/guard_test.go` | `flai/internal/guard/guard.go`, `guard_test.go`, `flai/cmd/guard_test.go` |
| 2 | T-2034, after T-2032: the `flai guard --help` text, `docs/users/flai-reference.md` regenerated with `make flai-reference`, `docs/users/flai.md` § Sub-agents, `design/system/agent-adapters.md` § The guard, `design/system/component-library.md` § The dev server | those five files |
| 3 | T-2036, after both: `flai issue close I-0165 --reason`, run in the worktree | the issue file, `design/issues/summary.md` |

**Figures.** Touches: ten files, no folder. Forecast 40m, delivery 2026-10-22T21:45Z, raised from flai's 18m and 21:21Z because the close-out's tiers alone take about 12 minutes (S-0466) and the reference must be regenerated. Cost of delay value 25 USD a week from the operator's input of 10m per 168h cycle, unchanged. Tags `flai`, `guard`, `components` and topics `cli`, `components` were added so the pack reaches `agent-adapters.md` and `component-library.md`; the story's words were not changed.

**Assumptions.**

1. `components serve` counts as a read for the guard. It runs the sub-project's own Vite dev server on 127.0.0.1, writes only under the git-ignored `.flai-cache` and the worktree's `node_modules`, and so changes nothing in the repository, as `flai test` runs processes and writes caches already. No ADR is needed for a read-list entry; S-0210 and S-0295 added theirs without one. The trade-off, that a read-only verifier may now start a dev server the ADR-0141 write rule cannot see, goes in the story's `## Decisions`.
2. `design/system/flai-cli.md` needs no change: its `flai components` row says nothing about the guard.

**To settle: merge with S-0472.** S-0472 remediates I-0156, filed from S-0438 on 2026-10-17 for the same missing `cliReads` entry; I-0165 is the second instance, from S-0449. Both stories are in the backlog or started, with the same one-line fix. Recommended: S-0477, now in progress, closes both issues, T-2036 says so already if this is confirmed, and the operator cancels S-0472 with a reason naming S-0477. Alternatives: cancel S-0477 and fold I-0165 into S-0472, which costs the work in progress; or let S-0472 run after S-0477 and close I-0156 alone, a story for one `flai issue close`. Closing I-0156 here means a second criterion the story's words do not carry, so it is the operator's to add or to allow by answer on this thread.

### 2026-10-22T20:59:55Z orchestrator
Approved under plan_backlog_stories: one `cliReads` entry for `components list` and `serve`, with `components config` still refused. That is three tasks in three layers, T-2032, then T-2034, then T-2036, on ten file touches with no folder. The 40m forecast stands, given the close-out's tiers and the regenerated reference.

S-0477 is already in progress. Its agent reviews these tasks against its own plan, as `work-management.md` says.

Merging S-0472 into this story adds a criterion and cancels a story, which is a change of scope. A recommendation for the operator follows, and the thread stays open for it.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-10-22T20:59:59Z orchestrator (recommendation)
Recommended for the operator: yes to the merge.

- Add a criterion to S-0477: "I-0156 is closed with `flai issue close I-0156 --reason` saying what fixed it". Its agent closes it in T-2036 beside I-0165.
- Cancel S-0472 as a duplicate, with a reason naming S-0477.

I-0156 (from S-0438) and I-0165 (from S-0449) record the same missing `cliReads` entry. S-0477 is in progress with the one-line fix. Run as a second story, S-0472 would change the same line again, or would exist for one `flai issue close`. If you say no, both stories stay as they are, and S-0472 closes I-0156 after S-0477 lands.

S-0472 is in ready, fourth in the order. Answer before an agent pulls it.

Source: design/conventions/strategic-agents.md § As the planner
