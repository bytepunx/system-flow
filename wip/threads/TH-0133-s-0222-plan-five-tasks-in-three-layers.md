---
id: TH-0133
title: "S-0222 plan: five tasks in three layers"
anchor:
  path: wip/kanban/stories/S-0222-the-orchestrator-publishes-by-the-release-policy-judgement-theme-or-cost-of-delay-threshold.md
  item: S-0222
status: resolved
participants: [planner-S-0222, alex]
created: 2026-10-05T04:48:43Z
updated: 2026-10-05T05:03:00Z
---

# TH-0133 S-0222 plan: five tasks in three layers

On wip/kanban/stories/S-0222-the-orchestrator-publishes-by-the-release-policy-judgement-theme-or-cost-of-delay-threshold.md.

## Entries

### 2026-10-05T04:48:43Z planner-S-0222
I planned S-0222: five tasks in the backlog, the story's touches widened, its forecast and cost of delay written, and topics `orchestration` and `release` added.

## Tasks and layers

| Layer | Task | Waits for | Changes |
|-------|------|-----------|---------|
| 1 | T-0885 `orchestration.release.whole_epics`, and `flai release --evaluate` names the batch it holds back | — | `manifest`, `release/evaluate.go`, `cmd/release.go`, `project-manifest.md`, `settings.md` |
| 1 | T-0886 the guard passes `release_publish` to the orchestrator only with `publish` on, and refuses its other routes to the remote | — | `guard.go` |
| 2 | T-0891 the `release_publish` MCP tool publishes as `publish.run` does, refuses when the policy is not met, and returns the refusal | T-0885 (`held_by_epic`) | `mcpserver/orchestrate.go`, `folder.go`, `hostapi/writes.go` |
| 2 | T-0895 the orchestrator's prompt: evaluate after each acceptance, publish, log, and open a thread on a refusal | T-0886 (the guard's rules) | `harness.go`, `.claude/agents/orchestrator.md` and its template copy |
| 3 | T-0897 an ADR refining ADR-0067, plus the design and the guides | T-0891, T-0895 | `design/adrs`, `strategic-agents.md`, `flai-cli.md`, `flai.md`, `docs/operators/index.md` |

The two tasks in each of layers 1 and 2 share no path, so they can run together. Criterion 4's tests are in T-0885 (each policy, met and not, with `whole_epics`), T-0891 (publish and refuse per policy, and both refusals on git fixtures), and T-0886 (the permission).

## Story

- **Touches:** I kept all eight declared touches. I added `flai/internal/mcpserver`, `flai/cmd/release.go`, the orchestrator's definition and its template copy, `design/adrs`, `flai-cli.md`, `docs/users/flai.md`, and `docs/operators/index.md`. Notes › Planning says where each came from.
- **Forecast:** 1h15m, delivery 2026-10-05T15:16Z. `flai forecast` gave 45m. I raised it because the work builds on S-0217's and S-0218's code, which is not written yet, and because of the git fixtures.
- **Cost of delay:** 76.53 USD/week, the story's share of E-0016's 1500 by forecast duration. It was 60.98.

## Assumptions (say if any is wrong)

1. Orchestrator code comes from S-0218: its role, its `orchestration.permissions.publish` check in the guard, its prompt beside `planPrompt` in `harness.go`, and its definition in `.claude/agents/orchestrator.md`. The other orchestration code comes from S-0217: `orchestration.release`, `evaluate.go`, and `release_evaluate` in `mcpserver/orchestrate.go`. The story's `after` on S-0218 brings both in.
2. The orchestrator publishes through a new MCP tool, `release_publish`. It runs what `publish.run` runs, behind the same `push` host action, instead of calling the dashboard's hostapi. The guard refuses the orchestrator's `flai release --pending`, `flai push`, `git push`, and `git tag`, so the tool is its only route.
3. The tool enforces the policy itself. It refuses under `threshold` or `theme` when the policy is not met, and under any policy when `whole_epics` holds the batch back. Under `judgement` it needs a reason. The agent's judgement adds to these checks and never replaces them.
4. `whole_epics` holds a batch back under every policy, not only `judgement`. It is off when unset. A story with no epic is never held.
5. "After each acceptance" means a `wait_for_events` change moving a story to `done`. No new event or hook is added.
6. On a refusal, the thread opens on the most recently accepted story of the batch. The orchestrator tries again only after that thread is answered or at the next acceptance.
7. An ADR refines ADR-0067, because the operator's `publish` permission counts as asking. If S-0218's ADR already covers it, T-0897 links that ADR and writes none.
8. The refusal stays `flai release --pending`'s exit 3 and its words, returned as a conflict. No new refusal codes are added.

## Proposals

None to split, merge, or drop.

One defect for a later story: the MCP `item_edit` (flai 1.31.1 here) refused the touch `.claude/agents/orchestrator.md` ("letters, digits, and _ . / @ + - only, not starting with a dash"). `flai touches` on the tree's flai and `item_new` both took it.

### 2026-10-05T05:03:00Z alex
Resolved.
