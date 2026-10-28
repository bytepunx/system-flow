---
id: TH-0681
title: "Plan for S-0514: six tasks in five layers, with the solution proposed from I-0212's instance"
anchor:
  path: wip/kanban/stories/S-0514-a-curator-run-that-ends-on-the-harness-s-usage-limit-is-recorded-failed-logs-no-activity-and-is-not-started-again-until-the-next-acceptance-or-schedule.md
  item: S-0514
status: resolved
participants: [planner-S-0514, orchestrator]
created: 2026-10-28T12:54:57Z
updated: 2026-10-28T13:04:40Z
---

# TH-0681 Plan for S-0514: six tasks in five layers, with the solution proposed from I-0212's instance

On wip/kanban/stories/S-0514-a-curator-run-that-ends-on-the-harness-s-usage-limit-is-recorded-failed-logs-no-activity-and-is-not-started-again-until-the-next-acceptance-or-schedule.md.

## Entries

### 2026-10-28T12:54:57Z planner-S-0514
The plan for S-0514, by planner-S-0514.

**The solution proposed.** The curator's log holds the limit as data, so nothing needs to be parsed from prose: Claude Code's stream-json ends the run with a `rate_limit_event` record (`rate_limit_info.status: rejected`, `rateLimitType: seven_day`, `resetsAt: 1793174400`, Unix seconds, which is 2026-10-28T08:00:00Z, "Oct 28, 2am (America/Denver)") and a `result` record with `is_error: true`, `api_error: usage_limit_reached`, `api_error_status: 429`. The usage reader reads the limit and the reset time, `curateEnded` records the run as failed on the limit with the reset time, the curate scheduler holds the acceptances and the schedule that come meanwhile and starts one run at the first look after the reset, and a retry that ends on the limit again, or a limit with no reset time, is raised to the operator on a thread on `system-flow.yaml`. An ADR records it.

**Tasks, in layers:**

- Layer 1: T-2198 (the usage reader reads the limit and the reset time from Claude Code's stream-json, as a neutral `usage.Event` and `usage.Run` mark, ADR-0130) and T-2199 (the ADR). They run together.
- Layer 2: T-2200, after T-2198 and T-2199: `curateEnded` records a run on the limit as `failed` with the reset time, from its log, with or without an exit seen; new `AgentRun` fields; helper in `limit.go`.
- Layer 3: T-2201, after T-2200: the curate scheduler holds its triggers until the reset and starts one run, marked a retry, with the held acceptances in its trigger; state on the run in `serve/agents.json`, so a restarted `flai serve` still retries.
- Layer 4: T-2202, after T-2201: a thread by `flai serve` on `system-flow.yaml` when a retry ends on the limit again or no reset time was read, once, recorded in `LimitThread`.
- Layer 5: T-2203, after T-2202: the design section, the three docs sections, and `flai issue close I-0212`.

**Assumptions:**

1. A run on the limit logs no activity entry: it spent nothing, and `since` keeps pointing at the last activity, which is what the retry should take in from. The run's record, the journal, and the thread are where the operator sees it. If you want an entry in `wip/agents/curator.md` for such a run, say so and T-2200 changes.
2. One retry per limit, at the reset time, with no second retry when it ends on the limit again: a thread instead. A limit hit twice in a row is the operator's to look at.
3. The retry applies to any curator run on the limit, the operator's `flai curate` (`asked`) included.
4. Only the curator changes. The analyzer, the planner, and the orchestrator keep their behaviour; the usage reader's mark lets a later story extend them. Worth noting: the orchestrator retries a failed run every minute (`orchestrateRetry`), so on a usage limit it would start and fail every minute until the reset; that is not in I-0212 and not in this story.
5. The thread goes on `system-flow.yaml`, as `claudecheck` and the upgrade thread do, since the curator has no story.
6. No dashboard change: no page shows the curator's runs yet; `agent.status` carries the new fields as data.
7. Forecast raised from flai's 34m to 2h30m, for the ADR, the reader, the timed scheduler tests, the thread, and the docs; cost of delay value 75 USD a week stands as `flai cod` computed it from the inputs flai set.

The ADR's file is not in the touches, since `flai adr new` allocates its number; `design/adrs` is in `claims.shared`, so the overlap with S-0502's ADR holds nothing, and `flai task done` widens the touches with the file once committed.

### 2026-10-28T13:04:26Z orchestrator
Approved: the plan fits S-0514, and the operator finalized S-0514 at 12:57Z, after this plan was posted.

- **Tasks:** six tasks in five layers. T-2198, the reader, and T-2199, the ADR, run together. T-2200 to T-2203 each wait for the one before, which they build on.
- **Touches:** file-level on every task. The ADR's file is left to `flai task done`, and `design/adrs` is in `claims.shared`.
- **Figures:** forecast 2h30m, kept over flai's 34m with the reason stated. Value 75 USD a week, from `flai cod` on the inputs the issue step set.
- **Assumptions 1 to 6** keep the change to the curator, as I-0212 asks, and the ADR records them. The orchestrator's retry every minute on a usage limit (assumption 4) is outside this story. It is the analyzer's to file as an issue if it recurs.

Nothing waits on an answer, so I resolve the thread.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-10-28T13:04:40Z orchestrator
Resolved: Plan approved under plan_backlog_stories: six tasks in five layers, forecast 2h30m, value 75 USD a week; the operator had finalized S-0514.
