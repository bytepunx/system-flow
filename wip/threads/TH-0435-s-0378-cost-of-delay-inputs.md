---
id: TH-0435
title: S-0378 cost of delay inputs
anchor:
  path: wip/kanban/stories/S-0378-the-dashboard-opens-a-markdown-editor-from-a-file-link-in-a-report-a-refusal-or-a-page-so-the-operator-can-fix-markdown-that-blocks-an-acceptance.md
  item: S-0378
status: resolved
participants: [planner-S-0378, orchestrator]
created: 2026-10-13T16:49:31Z
updated: 2026-10-13T16:49:40Z
---

# TH-0435 S-0378 cost of delay inputs

On wip/kanban/stories/S-0378-the-dashboard-opens-a-markdown-editor-from-a-file-link-in-a-report-a-refusal-or-a-page-so-the-operator-can-fix-markdown-that-blocks-an-acceptance.md.

## Entries

### 2026-10-13T16:49:31Z planner-S-0378
**Recommendation:** set `penalty_per_week` to 25 USD on S-0378.

S-0378 has no cost of delay inputs and no epic to take a share of, so `flai cod S-0378` cannot work out a value. The inputs are yours.

## Why 25

- Since S-0367 (ADR-0140), `flai accept` refuses an acceptance when markdown the fix cannot fix still has findings. Until S-0378 ships, you fix those files in an editor on the host and accept again.
- The fix (`scripts/lint-md.sh --fix`) clears most findings, so a refusal should be rare. Each costs a few minutes of your time and a review that waits.
- S-0367, the refusal itself, carries `penalty_per_week` 50. The editor halves the cost of a case that rarely happens, so half of that.

## Alternatives

- `time_lost_per_cycle` 10m: the time a refused acceptance costs you each time, if you would rather count it that way.
- `penalty_per_week` 50: if you expect refusals often, for example because agents write a lot of markdown the fix leaves.

I am planning the tasks, touches, and forecast meanwhile. I will write the value with `flai cod` once the inputs are set.

### 2026-10-13T16:49:38Z orchestrator
Set `penalty_per_week: 25` on S-0378, your recommended figure. Neither the thread nor the story gives a reason for 10m of time lost or a 50 penalty: the story records no frequency of refusals beyond TH-0413's one. Work out the value with flai cod on your run. Resolving.

### 2026-10-13T16:49:40Z orchestrator
Resolved: Orchestrator set penalty_per_week 25 on S-0378 under plan_backlog_stories.
