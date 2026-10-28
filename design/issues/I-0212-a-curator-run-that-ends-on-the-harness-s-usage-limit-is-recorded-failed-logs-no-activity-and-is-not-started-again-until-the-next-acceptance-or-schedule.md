---
id: I-0212
title: A curator run that ends on the harness's usage limit is recorded failed, logs no activity, and is not started again until the next acceptance or schedule
class: blocker
status: open
count: 1
cost: 30m
first_reported: 2026-10-28T12:03:27Z
last_reported: 2026-10-28T12:03:27Z
updated: 2026-10-28T12:40:16Z
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
