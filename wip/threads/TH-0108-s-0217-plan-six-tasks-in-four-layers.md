---
id: TH-0108
title: "S-0217 plan: six tasks in four layers"
anchor:
  path: wip/kanban/stories/S-0217-flai-exposes-the-orchestrator-s-deterministic-operations-as-commands-ordering-by-policy-promotion-candidates-and-release-evaluation.md
  item: S-0217
status: resolved
participants: [planner-S-0217, alex]
created: 2026-10-04T19:10:05Z
updated: 2026-10-04T19:13:57Z
---

# TH-0108 S-0217 plan: six tasks in four layers

On wip/kanban/stories/S-0217-flai-exposes-the-orchestrator-s-deterministic-operations-as-commands-ordering-by-policy-promotion-candidates-and-release-evaluation.md.

## Entries

### 2026-10-04T19:10:05Z planner-S-0217
I planned S-0217: six tasks in the backlog, the story's touches and forecast updated, and its cost of delay left as it was.

## Tasks and layers

| Layer | Task | Waits for | Changes |
|-------|------|-----------|---------|
| 1 | T-0809 `orchestration.policy` and `orchestration.release` in the manifest, reported by `flai check` | — | `internal/manifest`, `internal/check`, `project-manifest.md`, `settings.md` |
| 1 | T-0810 `flai order --by` and the ordering function, `--apply` | — | `workitem/policy.go`, `cmd/order.go` |
| 2 | T-0811 `flai promote --candidates [--limit]` | T-0809 (policy), T-0810 (ordering) | `workitem/promote.go`, `hold.go`, `cmd/promote.go`, `root.go` |
| 2 | T-0812 `flai release --evaluate` | T-0809 (release block) | `release/evaluate.go`, `cmd/release.go` |
| 3 | T-0813 hostapi reads, MCP tools, guard reads | T-0810, T-0811, T-0812 | `hostapi/reads.go`, `mcpserver/orchestrate.go`, `server.go`, `guard.go` |
| 4 | T-0814 `flai-cli.md`, `flai.md`, regenerated reference | T-0813 | docs |

The two tasks in each of layers 1 and 2 share no path, so they can run together. Each task's fixture tests sit in the task that writes the code, which covers criterion 5's tests.

## Story

- **Touches:** I added `flai/internal/check/check.go` and its new test, which report a bad `orchestration` block. The other touches stand.
- **Forecast:** 2h stands. Delivery moves from 12:00Z today to 2026-10-05T02:00Z, because S-0255 and S-0210 are still ahead of it on `flai/cmd`.
- **Cost of delay:** 97.56 USD/week stands. It is this story's share of E-0016's 1500 by forecast duration, and the duration did not change.

## Assumptions (say if any is wrong)

1. An unset `orchestration.policy` means `fifo`, which leaves your order alone. An unset `orchestration.release` means `judgement`.
2. `threshold` has two figures, `value` and `count`, and is met when either is at or over the manifest's figure. A story with no value counts toward `count` only.
3. `theme` names an epic or a tag, not both. It is met when every one of their stories is accepted and at least one is unreleased.
4. A promotion candidate also meets `flai move`'s definition of ready (goal, criteria, open epic), besides the criteria the story lists.
5. In the order, a missing figure sorts last and a tie keeps the current board order.
6. The hostapi read and the MCP tool for the order never apply it. Applying stays a CLI write.
7. The three MCP tools join the guard's reads, so sub-agents and the planner may call them. `.claude/settings.json` does not change.
8. `orchestration.permissions` is left to S-0218.
9. No ADR, since `metrics.md` does not change.

## Proposals

None to split, merge, or drop. Two notes for later stories:

- S-0213's cod-order chart should read the `order.by` read from T-0813.
- S-0222 should read `release.evaluate`.

S-0217 still waits for S-0210. Their shared `flai/cmd` and `flai/internal/manifest` would hold it anyway.

### 2026-10-04T19:13:57Z alex
Resolved.
