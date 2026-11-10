---
id: I-0278
title: After a stop from the dashboard, the story's agent page shows the agent stopped by the operator with no Retry, while the card's menu offers Retry
class: impression
status: closed
count: 1
cost: 2m
first_reported: 2026-11-07T06:28:23Z
last_reported: 2026-11-07T06:28:23Z
updated: 2026-11-10T09:54:42Z
---

# I-0278 After a stop from the dashboard, the story's agent page shows the agent stopped by the operator with no Retry, while the card's menu offers Retry

## Description
After a stop from the dashboard, the story's agent page shows the agent stopped by the operator with no Retry, while the card's menu offers Retry

## Instances

### 2026-11-07T06:28:23Z
Story: S-0565.
Met on S-0565's criterion 2, on S-0556's container run of 2026-11-07. The operator pressed Stop at 06:24:57Z. The run ended `stopped`, why "stopped by the operator at 2026-11-07T06:24:57Z", exit 143. The operator reported on TH-0836 that the story's agent page then showed it "stopped by the operator", with no Retry button. Right-clicking the card offered Retry, and that Retry started a new agent at 06:25:30Z. Both places use the same rule, `retryable()` in `flaiover/src/lib/activity.ts`, through `StoryAgent.svelte`'s `canRetry` and the card's menu. Nothing in it depends on the harness, so this is probably not specific to `harness: container`. Not reproduced: the page's `activity` and `status` at that moment were not captured. ADR-0058 §4 says the page shows the stopped agent with Retry.

## Remediation

Story S-0669 remediates this issue, created from it at 2026-11-10T08:34:44Z.
Closed 2026-11-10T09:54:42Z: Fixed by S-0669. The page with no Retry was the Activity page, where TH-0836 had the operator press Stop. ActivityView.svelte drew Stop under stoppable() and never drew Retry. The story page's StoryAgent and the card's menu both offered Retry under retryable(). ActivityView now offers Retry in Stop's place under retryable(), with writable from /api/board (T-2909, 9d0e68af1). The regression test is ActivityView.svelte.test.ts, "offers Retry for an agent the operator stopped, under the rule the card menu uses (I-0278)" (T-2907, eca43a87). It failed before the fix. activity.test.ts pins Retry for a stopped run of a story in progress.
