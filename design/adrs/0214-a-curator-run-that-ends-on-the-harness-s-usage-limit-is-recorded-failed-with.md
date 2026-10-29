---
id: ADR-0214
title: "A curator run that ends on the harness's usage limit is recorded failed with its reset time and holds the curator's starts until then, one run starts again at the reset with the acceptances held meanwhile, and a retry that ends on the limit again is raised to the operator on a thread"
status: accepted
date: 2026-10-29
supersedes: []
superseded_by: []
refines: [ADR-0196, ADR-0199]
topics: [planning, digest]
---

# ADR-0214 A curator run that ends on the harness's usage limit is recorded failed with its reset time and holds the curator's starts until then, one run starts again at the reset with the acceptances held meanwhile, and a retry that ends on the limit again is raised to the operator on a thread

## Context

[ADR-0196](0196-the-curator-runs-behind-the-curate-host-action-and-takes-in-what-the-archive.md) records a curator run in `serve/agents.json`. It ends `worked` on a clean exit and `failed` on a failed one. [ADR-0199](0199-the-curator-runs-after-acceptances-coalesced-and-on-digest-schedule-marks.md) starts the curator after acceptances and on `digest.schedule`. A start the checks refuse is not tried again until the schedule next comes round. Neither ADR says what happens when the run starts and the harness refuses it for the account's usage limit.

[I-0212](../issues/I-0212-a-curator-run-that-ends-on-the-harness-s-usage-limit-is-recorded-failed-logs-no-activity-and-is-not-started-again-until-the-next-acceptance-or-schedule.md) records that case. Accepting S-0498 and S-0501 started the curator on system-flow at 2026-10-26T16:03:52Z. Claude Code answered at once with the weekly limit, and the run exited 1 two seconds later. `serve/agents.json` recorded `failed` with `why: ended (exit 1)`, which does not say why. The run used no model, so no activity was logged. `digest.schedule` was unset, and no acceptance came after the reset. So nothing started the curator for two days, and an agent `flai serve` started cannot call `curate`.

The log says what happened as data. Claude Code's stream-json ends such a run with a `rate_limit_event` record, whose `rate_limit_info` has `status: rejected`, the limit's type, and `resetsAt` in Unix seconds. Then comes a `result` record with `is_error: true`, `api_error: usage_limit_reached`, and `api_error_status: 429`. In I-0212's run, `resetsAt` was 1793174400, 2026-10-28T08:00:00Z. The harness's usage reader ([ADR-0130](0130-every-harness-meets-flai-through-neutral-contracts-an-adapter-s-capabilities-a.md)) turns these records into a neutral mark: the run ended on a usage limit, with the limit named and the reset time when one is given. S-0514 builds that reader.

[ADR-0209](0209-the-orchestrator-runs-for-a-reason-and-ends-started-by-flai-serve-on-the-events.md) §5 decides a host-wide `budget` state for every role and every story's agent on a usage-limit refusal, with a shared hold, a probe at the reset, and a backoff. Story S-0534 builds it, on this story's reader. Until then the curator needs its own answer.

## Decision

**A curator run whose log ends on the harness's usage limit is recorded `failed` with the limit and its reset time, logs no activity, and holds the curator's starts until the reset; at the reset one run starts again, with the held acceptances; a retry that ends on the limit again, or a limit with no reset time, is raised to the operator once on a thread.** This refines ADR-0196 and ADR-0199; the rest of each stands.

1. **Recorded as a limit.** When the usage reader marks a curator run's log as ended on the usage limit, the run is recorded in `serve/agents.json` as `failed`. Its `why` names the limit and when it resets, such as `ended on the harness's usage limit (usage_limit_reached, 429); it resets at 2026-10-28T08:00:00Z`. The reset time is also a field of its own, `resets_at`, an RFC 3339 time in UTC. This holds whether or not `flai serve` saw the run's exit: a run settled from its log alone is recorded the same way. The run logs no activity in `wip/agents/curator.md`. It spent nothing, and the next run's since keeps pointing at the last activity.
2. **Held until the reset.** While the curator's last run ended on the limit and its `resets_at` has not passed, the curate scheduler starts nothing. Acceptances seen meanwhile are held, as ADR-0199 §1 holds those seen while a run is going. Schedule times that come round meanwhile are held too.
3. **One retry at the reset.** At the first look once `resets_at` has passed, the scheduler starts one run. Its trigger is the trigger of the run it retries, with the held acceptances joined: retrying `accepted S-0498, S-0501` with S-0510 accepted meanwhile gives `accepted S-0498, S-0501, S-0510`. A retried `schedule` trigger with held acceptances gives way to the acceptances' trigger, as ADR-0199 §2 has it. The run is marked `retry: true` on its record, and its start is journalled as a retry. It passes the same checks as any start (ADR-0196 §2). A `flai curate` the operator asked for is retried the same way, with its trigger `asked` and its since.
4. **Kept on the run.** The hold and the retry are read from the curator's last run in `serve/agents.json`. A `flai serve` restarted meanwhile still holds, and still retries at its first look after the reset.
5. **Behind the host action.** While `curate` is off, nothing is retried and no acceptance is held, as ADR-0199 has it. A retry that came due while `curate` was off is started at the first look after it is turned on again, since the curator's last run still ended on the limit.
6. **The operator's word is not held.** A `flai curate` the operator runs during the hold starts as asked. If it ends on the limit, it is the last run, and the hold and the retry follow from it.
7. **Raised once.** Two cases are raised to the operator: a retry that ends on the limit again, and a run whose log names the limit but no reset time. `flai serve` opens one thread on `system-flow.yaml`, by `flai serve`. It names the run, the reset time when known, and what starts the curator again: `flai curate`, the next acceptance, or `digest.schedule`. The thread's ID is recorded on the run as `limit_thread`, so it is opened once. No further retry is scheduled for that limit. The next start comes from one of the three causes the thread names.
8. **The other roles keep their behaviour.** The analyzer, the planner, and the orchestrator are not changed here.
9. **The curator's part of ADR-0209 §5.** This is the curator's share of §5, built ahead of the host-wide `budget` state. When that state lands (S-0534, which reads limits with this story's reader), it holds the curator's starts as it holds every role's. The curator's retry at the reset then becomes that state's probe, or is folded into it.

## Consequences

- In I-0212's case, the curator would have started again at the first look after 2026-10-28T08:00:00Z, with the trigger `accepted S-0498, S-0501`.
- A limit costs at most one more refused start before the operator hears of it.
- The record says why the run failed and when the limit resets, without reading the log.
- A held acceptance waits until the reset. A curator that started at once would only be refused again.
- `serve/agents.json`'s curator run gains `resets_at` and `retry`, and uses `limit_thread` as a story's run does under [ADR-0108](0108-flai-serve-restarts-a-story-s-agent-that-ended-with-its-story-in-progress-up-to.md).
- The analyzer, the planner, and the orchestrator still fail on the limit as before until S-0534.
- `design/system/strategic-agents.md` and the serve docs change in a later task of S-0514.

## Alternatives considered

- **A thread alone, with no retry.** The operator would have to act on every limit. The reset time is known, and a retry at it is safe.
- **A retry for every strategic agent now.** ADR-0209 §5 decides that host-wide, with a shared state and a backoff, which S-0534 builds. A retry per role here would duplicate it.
- **The reset time parsed from the result's prose.** The text reads "resets Oct 28, 2am (America/Denver)". It depends on the locale and the time zone. `resetsAt` needs neither.
- **An activity entry for a run on the limit.** It would move the next run's since past acceptances the curator never took in.
- **Retrying until the limit clears.** A limit hit twice in a row is the operator's to look at. Repeated retries would also anticipate ADR-0209 §5's backoff.
