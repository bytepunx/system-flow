---
id: D-0355
title: "Plan for S-0473: pin the test's run times, then close I-0158: Set COD to 190"
kind: recommendation
sources: [TH-0574]
topics: [flai, cli, go, code, all]
paths: [flai/internal/serve/agents_test.go, design/issues/I-0158-testasettledreplystartstheagentwaitingonitonce-in-flai-internal-serve-fails-now-and-then-in-the-close-out-s-integration-tier-on-a-branch-that-does-not-touch-serve.md, design/issues/summary.md, design/issues/I-0183-flai-check-finds-item-chronology-outside-the-story-at-close-out.md]
date: 2026-10-21
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0355 Plan for S-0473: pin the test's run times, then close I-0158: Set COD to 190

The operator, alex, answered TH-0574 at 2026-10-21T12:43:36Z.

What was asked:

```text
The plan is written: the story's touches, a 45m forecast, a cost of delay value of 112.50 USD a week, and two tasks.
```

The recommendation it overruled:

```text
Recommended for the operator: raise S-0473's `time_lost_per_cycle` from 45m to 73m, which is about 182.50 USD a week. The planner reckons I-0158 at 5 occurrences in 4.3 days, about 8 a week, and one of them sent S-0463 back out of the merge queue.

It changes no order today, since the policy is throughput and S-0473 is already in progress. It does record what the flake cost for `flai stats`. The story's inputs are yours, and my permission does not reach a story in progress, so set it with `flai edit S-0473` and its cost of delay flags. Otherwise leave 45m as it is.

Source: design/system/strategic-agents.md § Cost of delay
```

The operator's answer:

```text
Set COD to 190
```

Source: TH-0574.
