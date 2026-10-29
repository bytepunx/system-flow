---
id: I-0212
title: A curator run that ends on the harness's usage limit is recorded failed, logs no activity, and is not started again until the next acceptance or schedule
class: blocker
status: closed
count: 1
cost: 30m
first_reported: 2026-10-28T12:03:27Z
last_reported: 2026-10-28T12:03:27Z
updated: 2026-10-29T01:25:39Z
---

# I-0212 A curator run that ends on the harness's usage limit is recorded failed, logs no activity, and is not started again until the next acceptance or schedule

## Description
A curator run that ends on the harness's usage limit is recorded failed, logs no activity, and is not started again until the next acceptance or schedule

## Instances

### 2026-10-28T12:03:27Z
Story: S-0504.
S-0504 waited on the curator's first run on system-flow. `curate` was turned on, and accepting S-0498 and S-0501 started the curator at 2026-10-26T16:03:52Z (journal: serve.curate). Claude Code answered at once with 429 usage_limit_reached ("You've hit your weekly limit · resets Oct 28, 2am"), and the run exited 1 two seconds later. `serve/agents.json` records the curator `outcome: failed, why: ended (exit 1)`. `logRunEndSaying` logs no activity for a run with no model usage, so `wip/agents/curator.md` was never written. Nothing started the curator again after the limit reset: `digest.schedule` is unset here, and no acceptance came between the reset and 2026-10-28T12:02Z. An agent that flai serve started cannot call `curate`. So two days later the operator had confirmed the remedy on TH-0649, but no curator had run, and S-0504 could not close I-0207. A run that fails on a usage limit could be retried once the limit resets, which the harness's message names. Failing that, it could be raised to the operator on a thread rather than only in agents.json.

## Remediation

Story S-0514 remediates this issue, created from it at 2026-10-28T12:40:16Z.
Closed 2026-10-29T01:25:39Z: Fixed by S-0514 (ADR-0214). The usage reader reads Claude Code's rate_limit_event and the result's usage_limit_reached/429 into a neutral Limit with the reset time. curateEnded records a curator run whose log ends on the limit as failed, with a why naming the limit and its reset time and resets_at on the record, whether its exit was seen or not; it still logs no activity, so the next run's since stays at the last activity. The curate scheduler holds acceptances and schedule times until the reset, then starts that run again once, with the held acceptances in its trigger, marked retry; the state is the run in serve/agents.json, so a restarted flai serve still retries. A retry that ends on the limit again, or a limit with no reset time, opens one thread on system-flow.yaml by flai serve. TestARunOnTheUsageLimitIsStartedAgainByASchedulerMadeAfterTheReset reproduces I-0212's run (accepted S-0498, S-0501, reset 2026-10-28T08:00Z) and its retry; TestACuratorRunOnTheUsageLimitIsRecordedFailedOnIt its record from the incident's log.
