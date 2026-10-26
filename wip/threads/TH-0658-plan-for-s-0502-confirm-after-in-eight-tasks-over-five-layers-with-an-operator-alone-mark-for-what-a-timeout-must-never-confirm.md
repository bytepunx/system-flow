---
id: TH-0658
title: "Plan for S-0502: confirm_after in eight tasks over five layers, with an operator-alone mark for what a timeout must never confirm"
anchor:
  path: wip/kanban/stories/S-0502-a-thread-with-a-pending-orchestrator-recommendation-waits-on-the-operator-however-long-it-takes-with-no-setting-to-confirm-it-after-a-threshold.md
  item: S-0502
status: open
participants: [planner-S-0502]
created: 2026-10-26T04:47:53Z
updated: 2026-10-26T04:47:53Z
---

# TH-0658 Plan for S-0502: confirm_after in eight tasks over five layers, with an operator-alone mark for what a timeout must never confirm

On wip/kanban/stories/S-0502-a-thread-with-a-pending-orchestrator-recommendation-waits-on-the-operator-however-long-it-takes-with-no-setting-to-confirm-it-after-a-threshold.md.

## Entries

### 2026-10-26T04:47:53Z planner-S-0502
Plan by planner-S-0502. The story is a draft from I-0195, so its words stand. I added tags `flai` and `template`, topics `orchestration`, `cli`, and `analysis`, 38 file touches, `after: S-0498`, a forecast of 2h (flai's 1h doubled, reason under `### Planning`), and the cost of delay value 9000 USD a week from the 60h the orchestrator set on TH-0656.

## The proposed solution

The one the operator confirmed on TH-0598, which I-0195 records, with its exclusions made mechanical:

1. `orchestration.permissions.confirm_after`, a duration such as `48h`, off when unset. The one permission that is a duration, declared as `answer_threads` is the one that is a string.
2. `flai serve` confirms, on its per-minute look while the `orchestrate` host action is on, a pending recommendation by `orchestrator` that has waited `confirm_after` with no entry by the owner since, on a thread an agent opened (`agent-`, `planner-`), not one `flai serve` or the permission prompt opened, and not one marked for the operator alone. The entry is by `flai serve`: `Confirmed the recommendation of <time> orchestrator after 48h with no answer from the operator (orchestration.permissions.confirm_after).`, with the source, status `answered`, so the asking agent is started again as on the operator's confirmation. It is logged in `wip/agents/orchestrator.md` with no cost.
3. The orchestrator marks a recommendation on money or a host change for the operator alone: `thread_reply` `operator_only: true`, `flai thread reply --recommend --operator-only`, written as a one-line paragraph above the `Source:` line. flai cannot tell those classes from the text; the orchestrator already classes them under `answer_threads`.
4. `flai stats` gains `waiting.weeks[].threads.timeout` beside `orchestrator` and `confirmed`, by the ADR, so the operator's confirmations and the timeouts are counted apart (ADR-0091).

## Tasks and layers

| Layer | Tasks | After | Touches |
|-------|-------|-------|---------|
| 1 | T-2141 ADR, refining ADR-0090 and ADR-0091 | none | `design/adrs/` (new ADR, README) |
| 2 | T-2142 the permission; T-2143 the mark, `WroteSince`, the confirmation by timeout, the CLI flag, the MCP parameter | T-2141 | manifest, settings, check, settings docs; threads, cmd/thread, mcpserver |
| 3 | T-2144 metrics and metrics.md; T-2145 the orchestrator's prompt and agent definition; T-2146 the serve look and its log | T-2143; T-2143; T-2142 and T-2143 | metrics/waiting.go, charts.ts; harness.go, `.claude/agents`; serve/confirm.go, serve.go |
| 4 | T-2147 design, convention in both copies, user docs, reference, changelog | T-2144, T-2145, T-2146 | strategic-agents.md, flai-cli.md, conventions, docs |
| 5 | T-2148 this project's value, the criteria ticked, I-0195 closed | T-2147 | system-flow.yaml, I-0195, summary.md |

No two tasks of one layer share a path.

## Assumptions

- `after: S-0498`: its T-2128 declares the `threads` block in the same manifest, settings, and check files, its T-2131 adds `remind.go`'s look to the same block of `serve.go`, and its T-2132 rewrites the same design sections; the confirmer scans the threads awaiting the operator past a threshold as the reminder does. S-0502 is held in ready until S-0498 is accepted either way, since their touches overlap.
- The look runs only while the project's `orchestrate` host action is on, since the permission is the orchestrator's; with it off, no new recommendations are posted and old ones keep waiting for you.
- The exclusions by opener: a thread `flai serve` opened (upgrades, restarts, forge failures, Claude Code checks) and one the permission prompt opened (`Allow … ?`, which covers `.claude/` writes) are never confirmed, whatever the recommendation says. The owner's entry after the recommendation excludes it whoever opened the thread, which `PendingRecommendation` alone would not tell when you opened it.
- Money and host changes are the orchestrator's to mark, since flai cannot tell them from the text. Alternative: flai excludes every planner's thread, which would have excluded the two instances (TH-0583, TH-0588); rejected. Alternative: text heuristics; rejected as brittle.
- The ADR is one, refining ADR-0090 (who confirms) and ADR-0091 (the metrics), predicted as 0206; `flai task done` widens the touch to the name `flai adr new` gives. ADR-0090 and ADR-0091 are `proposed` on main although built; T-2141 refines them as they are.
- `.claude/agents/orchestrator.md` is a protected path: the story's agent writes it itself, `permission_prompt` asks you, and the story is accepted by you alone (ADR-0106).
- The dashboard changes only its `WaitWeek` type; no new chart. The orchestrator page's settings panel renders the duration from flai's catalog with no change of its own.
- `flai guard` needs no change: it reads `--recommend` and lets the other flags of `flai thread reply` through, and it already refuses every agent `flai thread confirm`; `flai serve` is not an agent.

## One question, recommendation first

**What value should this project set for `confirm_after` (T-2148)?** Recommended: `48h`, the figure you used on TH-0598. Alternatives: `24h`, since the orchestrator's recommendations here are mostly on planners' design choices and a day is what I-0151 calls a long wait; or leave it unset here and ship the setting for other projects only. T-2148 takes the figure from this thread and asks again if none is given.

Nothing to split, merge, or drop.
