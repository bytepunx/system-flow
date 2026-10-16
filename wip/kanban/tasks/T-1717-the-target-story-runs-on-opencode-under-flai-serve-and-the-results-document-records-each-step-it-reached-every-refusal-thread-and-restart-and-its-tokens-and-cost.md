---
id: T-1717
type: task
nature: experiment
title: The target story runs on OpenCode under flai serve, and the results document records each step it reached, every refusal, thread, and restart, and its tokens and cost
status: backlog
parent: S-0412
owner: alex
created: 2026-10-16T05:20:35Z
updated: 2026-10-16T05:20:35Z
transitions: []
stream: S-0412
tags: [cli]
touches: [design/experiments/S-0412-a-non-claude-model-works-one-story-of-this-repository-end-to-end-on-opencode-through-openrouter-with-the-guard-and-the-permission-prompt-on-and-the-results-record-what-broke-and-what-it-cost.md]
after: [T-1716]
---
# T-1717 The target story runs on OpenCode under flai serve, and the results document records each step it reached, every refusal, thread, and restart, and its tokens and cost

## Work

- Once the operator has moved the target story to ready and `flai serve` has started its agent on OpenCode, follow it through `flai show`, the board, its narrative, and its threads, without working it.
- Record in the results document each step it reached (pull, tasks, each `task_done`, close-out, review) or where it stopped and why; every guard refusal, permission thread, restart, and message; and its `usage` from `flai show`: tokens, cost, and `priced_by`, against the median cost of done Claude Code stories of the same size from `flai stats`.
- Waits for T-1716: the target and its agent are set there.

## Done when

The target story is in review or has stopped for good, and the results document records the run step by step with its usage.

## Notes
